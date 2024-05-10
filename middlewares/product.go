package middlewares

import (
	"net/http"
	"product_app/controller/productcontroller"
	"strconv"

	"github.com/go-chi/chi"
)

func APIKeyAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idParams, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Parameter tidak sesuai", http.StatusBadRequest)
			return
		}

		apiKey, err := strconv.Atoi(r.Header.Get("X-API-Key"))
		if err != nil {
			http.Error(w, "Api Key Tidak Sesuai", http.StatusBadRequest)
			return
		}

		checkAPI := productcontroller.CheckStoreProduct(idParams, apiKey)

		if !checkAPI {
			http.Error(w, "Anda tidak memiliki akses", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
