package main

import (
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("./esv/ot.bzz", os.O_RDONLY, 0600)

	zr, err := zlib.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}

	bibleBytes, err := ioutil.ReadAll(zr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(bibleBytes))

	fmt.Println(string(bibleBytes[len(bibleBytes)-100:]))
}
