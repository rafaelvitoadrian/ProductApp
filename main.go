package main

import (
	"log"
	"net/http"
	"product_app/controller/productcontroller"
	"product_app/controller/storecontroller"
	"product_app/database"
	"product_app/middlewares"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

func init() {
	database.ConnectDB()
}

func main() {
	defer database.DBConn.Close()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Post("/product/create", productcontroller.CreateProduct)
	r.Get("/product", productcontroller.GetAll)

	r.With(middlewares.APIKeyAuth).Get("/product/{id}", productcontroller.GetProductById)
	r.Put("/product/update", productcontroller.UpdateProduct)

	r.Get("/store", storecontroller.GetAll)
	r.Get("/store/{id}", storecontroller.GetStoreById)
	r.Post("/store/create", storecontroller.CreateStore)
	r.Put("/store/update", storecontroller.UpdateStore)
	r.Delete("/store/delete", storecontroller.DeleteStore)

	log.Fatal(http.ListenAndServe(":3000", r))
}
