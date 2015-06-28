package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Section struct {
	Start, Length, Uncompressed uint32
}

func (s *Section) Slice(bs []byte) []byte {
	return bs[s.Start : s.Start+s.Length]
}

func main() {
	f, err := os.OpenFile("./esv/ot.bzs", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	sects := make([]Section, 0, 20)

	for {
		s := Section{}
		check(binary.Read(f, binary.LittleEndian, &s.Start))
		check(binary.Read(f, binary.LittleEndian, &s.Length))
		done := check(binary.Read(f, binary.LittleEndian, &s.Uncompressed))

		sects = append(sects, s)

		if done {
			break
		}
	}

	fmt.Printf("sects = %+v\n", sects)

	ss := read("./esv/ot.bzz", sects)

	fmt.Println(len(sects), len(ss))
}

func check(err error) bool {
	if err != nil {
		if err == io.EOF {
			return true
		} else {
			log.Fatal(err)
		}
	}
	return false
}

func read(file string, sects []Section) [][]byte {
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	sections := make([][]byte, 0, 20)

	for _, s := range sects {
		zr, err := zlib.NewReader(bytes.NewBuffer(s.Slice(bs)))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		sectB, err := ioutil.ReadAll(zr)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		sections = append(sections, sectB)
	}

	return sections
}
