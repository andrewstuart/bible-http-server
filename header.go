package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var BadHeader = fmt.Errorf("Bad headers")

func readHeader(name string) ([]Section, error) {

	f, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	sects := make([]Section, 0, 20)

	for {
		s := Section{}
		err = binary.Read(f, binary.LittleEndian, &s.Start)
		err = binary.Read(f, binary.LittleEndian, &s.Length)
		err = binary.Read(f, binary.LittleEndian, &s.Uncompressed)

		if err == io.EOF {
			return sects, nil
		} else if err != nil {
			return nil, BadHeader
		}

		sects = append(sects, s)
	}
}
