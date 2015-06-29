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

	err = checkCycle(&b, "Gen.1.1")
	if err != nil {
		fmt.Println(len(visited))
		log.Fatal(err)
	}
}

var AlreadyVisited = fmt.Errorf("Visited")
var visited = make(map[string]bool)

func checkCycle(b *osis.Bible, r string) error {
	if strings.Index(r, "-") > -1 {
		return nil
	}

	v, err := b.GetVerse(r)

	if err != nil {
		return err
	}

	visited[v.ID] = true

	for i := range v.Refs {
		if visited[v.Refs[i].RefID] {
			fmt.Println(v.ID, v.Refs[i].RefID)
			//Anchor/escape
			return AlreadyVisited
		}

		if err := checkCycle(b, v.Refs[i].RefID); err != nil {
			return err
		}
	}
	return nil
}
