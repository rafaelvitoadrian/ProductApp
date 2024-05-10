package productcontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"product_app/database"
	"product_app/model"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

type ProductResponse struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Created_at  time.Time   `json:"created_at`
	Updated_at  time.Time   `json:"updated_at`
	Store       model.Store `json:store`
}

type ProductAllResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at`
	Updated_at  time.Time `json:"updated_at`
}

func GetAll(w http.ResponseWriter, r *http.Request) {

	stmt_s := "SELECT Id, Name, Description, Created_at, Updated_at FROM product"
	stmt, err := database.DBConn.Query(stmt_s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var data []ProductAllResponse

	for stmt.Next() {
		var d ProductAllResponse
		err := stmt.Scan(&d.Id, &d.Name, &d.Description, &d.Created_at, &d.Updated_at)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, d)
	}

	err = json.NewEncoder(w).Encode(data)
}

func CheckStoreProduct(idParams int, apiKey int) bool {
	stmt_s := "SELECT Id_store FROM product where id=$1"
	var idStore int
	err := database.DBConn.QueryRow(stmt_s, idParams).Scan(&idStore)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Id Tidak Ditemukan:", idParams)
			return false
		}
		panic(err.Error())
	}

	if apiKey == idStore {
		return true
	} else {
		return false
	}
}

func GetProductById(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	var rspProduct model.Product
	var rspStore model.Store

	stmt_s := "SELECT Id, Name, Description, Id_store, Created_at, Updated_at FROM product where id=$1"
	stmt := database.DBConn.QueryRow(stmt_s, idParam)

	err := stmt.Scan(&rspProduct.Id, &rspProduct.Name, &rspProduct.Description, &rspProduct.Id_store, &rspProduct.Created_at, &rspProduct.Updated_at)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No product found with ID:", idParam)
			return
		}
		panic(err.Error())
	}

	stmt_s = "SELECT Id, Name, Address FROM store where id=$1"
	stmt = database.DBConn.QueryRow(stmt_s, rspProduct.Id_store)
	err = stmt.Scan(&rspStore.Id, &rspStore.Name, &rspStore.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No store found with ID:", idParam)
			return
		}
		panic(err.Error())
	}

	rsp := ProductResponse{
		Id:          rspProduct.Id,
		Store:       rspStore,
		Name:        rspProduct.Name,
		Description: rspProduct.Description,
		Created_at:  rspProduct.Created_at,
		Updated_at:  rspProduct.Updated_at,
	}

	err = json.NewEncoder(w).Encode(rsp)

}

type ProductCreateRequest struct {
	Id_store    int    `json:"id_store" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product ProductCreateRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(product)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	stmt_s := "insert into product (name, description, id_store,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)"
	stmt, err := database.DBConn.Prepare(stmt_s)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Id_store, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Produk %s berhasil dibuat", product.Name)

}

type ProductUpdateRequest struct {
	Id          int    `json:"id" validate:"required"`
	Id_store    int    `json:"id_store" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product ProductUpdateRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(product)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation error: %s", errors), http.StatusBadRequest)
		return
	}

	stmt_s := "UPDATE product set name=$1, description=$2, id_store=$3, updated_at=$4 where id=$5"
	stmt, err := database.DBConn.Prepare(stmt_s)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Id_store, time.Now(), product.Id)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Produk %s berhasil diubah", product.Name)
}
