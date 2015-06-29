package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"os"
)

func readFile(file string, sects []Section) ([][]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	sections := make([][]byte, 0, 20)

	for _, s := range sects {
		zr, err := zlib.NewReader(bytes.NewBuffer(s.Slice(bs)))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		sectB, err := ioutil.ReadAll(zr)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		sections = append(sections, sectB)
	}

	return sections, nil
}
