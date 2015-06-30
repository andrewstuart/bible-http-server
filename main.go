package main

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewstuart/bible/osis"
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

var b osis.Bible

func main() {
	b = osis.Bible{}

	zipped, err := os.Open("./bibles/esv.xml.gz")
	if err != nil {
		log.Fatal(err)
	}

	r, err := gzip.NewReader(zipped)
	if err != nil {
		log.Fatal(err)
	}

	dec := xml.NewDecoder(r)
	err = dec.Decode(&b)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ESV Loaded")

	rtr := mux.NewRouter()

	rtr.HandleFunc("/{book}/{chapter}/{verse}", SetHeaders(GetVerse))

	http.ListenAndServe(":8080", rtr)
}
