package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func SetHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
}

func GetVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ref := fmt.Sprintf("%s.%s.%s", vars["book"], vars["chapter"], vars["verse"])
	v, err := b.GetVerse(ref)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(v)
}
