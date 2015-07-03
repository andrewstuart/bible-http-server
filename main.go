package main

import (
	"net/http"

	"github.com/andrewstuart/bible/osis"
	"github.com/gorilla/mux"
)

var b osis.Bible

func main() {
	r := mux.NewRouter()

	r.Path("/").HandlerFunc(SetHeaders(SearchVerse))
	r.Path("/{book}/{chapter}/{verse}").HandlerFunc(SetHeaders(GetVerse))

	http.ListenAndServe(":8089", r)
}
