package osis

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type Bible struct {
	BooksById  map[string]*Book
	XMLName    xml.Name    `xml:"osis"`
	Testaments []Testament `xml:"osisText>div"`
}

var (
	NoSuchVerse = fmt.Errorf("No such verse")
	InvalidRef  = fmt.Errorf("Invalid verse reference")
)

func (b *Bible) GetVerse(ref string) (*Verse, error) {
	strs := strings.Split(ref, ".")
	if len(strs) < 3 {
		return nil, InvalidRef
	}

	book := strs[0]

	chap, err := strconv.Atoi(strs[1])
	if err != nil {
		return nil, InvalidRef
	}
	//Human -> arr index
	chap--

	vs, err := strconv.Atoi(strs[2])
	if err != nil {
		return nil, InvalidRef
	}
	//Human -> arr index
	vs--

	if b.BooksById == nil {
		b.Index()
	}

	if bk, ok := b.BooksById[book]; ok {
		if chap < len(bk.Chs) && vs < len(bk.Chs[chap].Vrs) {
			return &b.BooksById[book].Chs[chap].Vrs[vs], nil
		}
	}

	return nil, NoSuchVerse
}

func (b *Bible) Index() {
	b.BooksById = make(map[string]*Book)
	for i := range b.Testaments {
		for j := range b.Testaments[i].Books {
			b.BooksById[b.Testaments[i].Books[j].ID] = &b.Testaments[i].Books[j]
		}
	}
}

type Testament struct {
	Books []Book `xml:"div"`
}

type Book struct {
	Chs []Chapter `xml:"chapter"`
	ID  string    `xml:"osisID,attr"`
}

type Chapter struct {
	Vrs []Verse `xml:"verse"`
	ID  string  `xml:"osisID,attr"`
}

type Verse struct {
	Text string      `xml:",chardata"`
	ID   string      `xml:"osisID,attr"`
	Refs []Reference `xml:"note>reference"`
}

type Reference struct {
	Text  string `xml:",chardata"`
	RefID string `xml:"osisRef,attr"`
}
