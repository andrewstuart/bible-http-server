package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewstuart/bible-http-server/osis"
	"github.com/gorilla/mux"
)

//CORSRouter wraps an http.Handler and adds a header.
type CORSRouter struct {
	r *mux.Router
}

func (cr CORSRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	cr.r.ServeHTTP(w, r)
}

var b osis.Bible

func main() {
	if len(os.Args) == 1 {
		serve()
		return
	}

	if len(os.Args) == 3 {
		if os.Args[1] == "import" {
			b, err := loadFromGzippedFile(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}

			v, err := b.GetVerse("Gen.1.1")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("v = %+v\n", v)

			err = store(b)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Successfully stored", os.Args[2])
		}
	}
}

func serve() {
	r := mux.NewRouter()

	r.Path("/").HandlerFunc(findVerses)
	r.Path("/verse/{verseId}").HandlerFunc(getVerseByID)

	r.Path("/text/{version:[a-zA-Z]+}/{book:[1-3]?[a-zA-Z]+}").HandlerFunc(getVerse)
	r.Path("/text/{version:[a-zA-Z]+}/{book:[1-3]?[a-zA-Z]+}/{chapter:[0-9]+}").HandlerFunc(getVerse)
	r.Path("/text/{version:[a-zA-Z]+}/{book:[1-3]?[a-zA-Z]+}/{chapter:[0-9]+}/{verse:[0-9]+}").HandlerFunc(getVerse)
	r.Path("/text/version/{version:[a-zA-Z]+}/book/{book:[a-zA-Z]+}/chapter/{chapter:[0-9]+}/verse/{verse:[0-9]+}").HandlerFunc(getVerse)

	r.Path("/text/{book:[1-3]?[a-zA-Z]+}").HandlerFunc(getVerse)
	r.Path("/text/{book:[1-3]?[a-zA-Z]+}/{chapter:[0-9]+}").HandlerFunc(getVerse)
	r.Path("/text/{book:[1-3]?[a-zA-Z]+}/{chapter:[0-9]+}/{verse:[0-9]+}").HandlerFunc(getVerse)
	r.Path("/text/book/{book:[a-zA-Z]+}/chapter/{chapter:[0-9]+}/verse/{verse:[0-9]+}").HandlerFunc(getVerse)

	r.Path("/books").HandlerFunc(getBooks)

	port := ":" + stringDef(os.Getenv("BIBLE_PORT"), "8080")
	log.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(port, CORSRouter{r})
	if err != nil {
		log.Fatal(err)
	}
}

func stringDef(s ...string) string {
	for i := range s {
		if s[i] != "" {
			return s[i]
		}
	}
	return ""
}
