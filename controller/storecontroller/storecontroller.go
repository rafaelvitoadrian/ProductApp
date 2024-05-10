package storecontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func GetAll(w http.ResponseWriter, r *http.Request) {

	stmt_s := "SELECT * FROM store"
	stmt, err := database.DBConn.Query(stmt_s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var data []model.Store

	for stmt.Next() {
		var d model.Store
		err := stmt.Scan(&d.Id, &d.Name, &d.Address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, d)
	}

	err = json.NewEncoder(w).Encode(data)
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
