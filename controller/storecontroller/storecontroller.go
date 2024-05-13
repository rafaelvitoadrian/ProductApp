package storecontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product_app/database"
	"product_app/model"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type ProductResponseStore struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GetStoreByIdResponse struct {
	Id      int                    `json:"id"`
	Name    string                 `json:"name"`
	Address string                 `json:"address"`
	Product []ProductResponseStore `json:"product"`
}

func GetStoreById(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	var rspProduct []ProductResponseStore
	var rspStore model.Store

	stmt_s := "SELECT Id, Name, Address FROM store where id=$1"
	stmt := database.DBConn.QueryRow(stmt_s, idParam)

	err := stmt.Scan(&rspStore.Id, &rspStore.Name, &rspStore.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Store tidak ditemukan", http.StatusUnauthorized)
			return
		}
		panic(err.Error())
	}

	fmt.Println(rspStore)

	stmt_s = "SELECT Id, Name FROM product where id_store=$1"
	stmt_r, err_r := database.DBConn.Query(stmt_s, rspStore.Id)
	if err_r != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt_r.Close()

	for stmt_r.Next() {
		var d ProductResponseStore
		err := stmt_r.Scan(&d.Id, &d.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rspProduct = append(rspProduct, d)
	}

	fmt.Println(rspProduct)

	rsp := GetStoreByIdResponse{
		Id:      rspStore.Id,
		Name:    rspStore.Name,
		Address: rspStore.Address,
		Product: rspProduct,
	}

	err = json.NewEncoder(w).Encode(rsp)

}

type ProductResponseStoreAll struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ResponseStoreAll struct {
	Id      int                       `json:"id"`
	Name    string                    `json:"name"`
	Address string                    `json:"address"`
	Product []ProductResponseStoreAll `json:"product"`
}

func GetAll(w http.ResponseWriter, r *http.Request) {

	storeRows, err := database.DBConn.Query("SELECT id, name, address FROM store")
	if err != nil {
		log.Fatal(err)
	}
	defer storeRows.Close()

	storeDetails := make(map[int]struct {
		name    string
		address string
	})

	for storeRows.Next() {
		var storeID int
		var storeName, storeAddress string

		err := storeRows.Scan(&storeID, &storeName, &storeAddress)
		if err != nil {
			log.Fatal(err)
		}

		storeDetails[storeID] = struct {
			name    string
			address string
		}{
			name:    storeName,
			address: storeAddress,
		}

		// fmt.Print(storeDetails[storeID])
		// fmt.Println(" ")
	}
	if err := storeRows.Err(); err != nil {
		log.Fatal(err)
	}

	productRows, err := database.DBConn.Query("SELECT s.id, p.id, p.name FROM store s LEFT JOIN product p ON s.id = p.id_store")
	if err != nil {
		log.Fatal(err)
	}
	defer productRows.Close()

	storeProducts := make(map[int][]ProductResponseStoreAll)

	for productRows.Next() {
		var storeID, productID sql.NullInt64
		var productName sql.NullString

		var idStore int
		var idProduct int
		var nameProduct string

		err := productRows.Scan(&storeID, &productID, &productName)
		if err != nil {
			log.Fatal(err)
		}

		if storeID.Valid {
			idStore = int(storeID.Int64)
		}

		if productID.Valid {
			idProduct = int(productID.Int64)
		}

		if productName.Valid {
			nameProduct = productName.String
		}

		storeProducts[idStore] = append(storeProducts[idStore], ProductResponseStoreAll{
			Id:   idProduct,
			Name: nameProduct,
		})
	}
	if err := productRows.Err(); err != nil {
		log.Fatal(err)
	}

	var responseStores []ResponseStoreAll
	for storeID, products := range storeProducts {
		storeDetail := storeDetails[storeID]

		if products[0].Id == 0 {
			products = make([]ProductResponseStoreAll, 0)
		}

		responseStore := ResponseStoreAll{
			Id:      storeID,
			Name:    storeDetail.name,
			Address: storeDetail.address,
			Product: products,
		}
		responseStores = append(responseStores, responseStore)
	}

	err = json.NewEncoder(w).Encode(responseStores)

}

type StoreCreateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

func CreateStore(w http.ResponseWriter, r *http.Request) {
	var store StoreCreateRequest

	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(store)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	stmt_s := "insert into store (name, address) VALUES ($1, $2)"
	stmt, err := database.DBConn.Prepare(stmt_s)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(store.Name, store.Address)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Store %s berhasil dibuat", store.Name)

}

func DeleteStore(w http.ResponseWriter, r *http.Request) {
	var store model.Store

	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt_s := "DELETE from store where id=$1"
	stmt, err := database.DBConn.Prepare(stmt_s)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(store.Id)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Produk berhasil dihapus")
}

type StoreUpdateRequest struct {
	Id      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

func UpdateStore(w http.ResponseWriter, r *http.Request) {
	var store StoreUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(store)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	stmt_s := "UPDATE store set name=$1, address=$2 where id=$3"
	stmt, err := database.DBConn.Prepare(stmt_s)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(store.Name, store.Address, store.Id)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Store %s berhasil diupdate", store.Name)
}
