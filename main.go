package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewstuart/bible/osis"
	"github.com/gorilla/mux"
)

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

			fmt.Printf("b = %+v\n", b)

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

	r.Path("/").HandlerFunc(SetHeaders(SearchVerse))
	r.Path("/verse/{verseId}").HandlerFunc(SetHeaders(GetVerseById))
	r.Path("/{book}").HandlerFunc(SetHeaders(GetVerse))
	r.Path("/{book}/{chapter}").HandlerFunc(SetHeaders(GetVerse))
	r.Path("/{book}/{chapter}/{verse}").HandlerFunc(SetHeaders(GetVerse))
	r.Path("/book/{book}/chapter/{chapter}/verse/{verse}").HandlerFunc(SetHeaders(GetVerse))

	http.ListenAndServe(":8089", r)
}
