package api

import (
	"errors"
	"io"
	"net"
	"net/http"

	"kv/domain"
	"kv/infrastructure/adaptors/middleware"
	"kv/infrastructure/repository/storage"

	"github.com/gorilla/mux"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func throttledKeyValuePutHandler(storage storage.Crud) http.HandlerFunc {
	throttled := middleware.Throttle(1, 1)
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if !throttled(ip) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		vars := mux.Vars(r)
		key := vars["key"]

		value, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = storage.Put(key, string(value))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func keyValueGetHandler(storage storage.Crud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := storage.Get(key)
		if errors.Is(err, domain.ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(value + "\n"))
	}
}

func keyValueDeleteHandler(storage storage.Crud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		err := storage.Delete(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetStoreHttpHandler(storage storage.Crud) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", helloMuxHandler)
	r.HandleFunc("/v1/{key}", throttledKeyValuePutHandler(storage)).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler(storage)).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler(storage)).Methods("DELETE")
	return r
}
