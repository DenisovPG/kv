package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"kv/store"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func keyValuePutHandler(storage store.Crud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
func keyValueGetHandler(storage store.Crud) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := storage.Get(key)
		if errors.Is(err, store.ErrorNoSuchKey) {
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

func keyValueDeleteHandler(storage store.Crud) http.HandlerFunc {
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

func main() {
	mem_store := store.NewMemoryStore()
	defer mem_store.Save()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/", helloMuxHandler)
		r.HandleFunc("/v1/{key}", keyValuePutHandler(mem_store)).Methods("PUT")
		r.HandleFunc("/v1/{key}", keyValueGetHandler(mem_store)).Methods("GET")
		r.HandleFunc("/v1/{key}", keyValueDeleteHandler(mem_store)).Methods("DELETE")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()
	<-stop
	log.Println("Shutting down...")
}
