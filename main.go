package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/andrewstuart/bible/osis"
)

func main() {
	b := osis.Bible{}

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

	fmt.Println(b.Testaments[0].Books[0].Chapters[0].Verses[0])
}
