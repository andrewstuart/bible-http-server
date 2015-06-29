package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"

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
	b.Index()

	v, err := b.GetVerse("Gen.1.1")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range v.Refs {
		if strings.Index(r.RefID, "-") > -1 {
			continue
		}
		fmt.Printf("r.RefID = %+v\n", r.RefID)
		ref, err := b.GetVerse(r.RefID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", ref)
	}
}
