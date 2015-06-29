package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type Section struct {
	Start, Length, Uncompressed uint32
}

func (s *Section) Slice(bs []byte) []byte {
	return bs[s.Start : s.Start+s.Length]
}

func main() {
	fmt.Println("\nOld Testament\n")

	otChaps, err := readGroup("./esv/ot")
	if err != nil {
		log.Fatal(err)
	}
	for i, chap := range otChaps {
		fmt.Println(i+1, len(chap))
	}

	fmt.Println("\nNew Testament\n")
	chapters, err := readGroup("./esv/nt")
	if err != nil {
		log.Fatal(err)
	}
	for i, chap := range chapters {
		fmt.Println(i+1, len(chap))
	}

	err = ioutil.WriteFile("./out.xml", chapters[24], 0660)
	if err != nil {
		log.Fatal(err)
	}

}

func readGroup(name string) ([][]byte, error) {
	sects, err := readHeader(name + ".bzs")
	if err != nil {
		return nil, err
	}
	ss, err := readFile(name+".bzz", sects)
	if err != nil {
		return nil, err
	}
	return ss, nil
}
