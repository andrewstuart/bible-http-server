package osis

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Bible struct {
	Books      []*Book        `json:"-"xml:"-"`
	BooksById  map[string]int `json:"-"xml:"-"`
	XMLName    xml.Name       `xml:"osis"json:"-"`
	Testaments []Testament    `xml:"osisText>div"json:"testaments"`
	Version    Version        `xml:"osisText>header>work"json:"version"`
}

func NewBible(r io.Reader) (*Bible, error) {
	b := &Bible{}
	dec := xml.NewDecoder(r)

	err := dec.Decode(b)
	if err != nil {
		return nil, err
	}
	b.index()

	return b, nil
}

var (
	NoSuchVerse = fmt.Errorf("No such verse")
	InvalidRef  = fmt.Errorf("Invalid verse reference")
)

func (b *Bible) GetVerse(ref string) (*Verse, error) {
	vr, err := NewVerseRef(ref)
	if err != nil {
		return nil, err
	}

	if b.BooksById == nil {
		b.index()
	}

	if bIdx, ok := b.BooksById[vr.BookID]; ok {
		bk := b.Books[bIdx]
		if vr.Chapter < len(bk.Chs) && vr.Verse < len(bk.Chs[vr.Chapter].Vrs) {
			return bk.Chs[vr.Chapter].Vrs[vr.Verse], nil
		}
	}

	return nil, NoSuchVerse
}
func (b *Bible) index() {
	b.BooksById = make(map[string]int)
	b.Books = make([]*Book, 0, 56)
	for i := range b.Testaments {
		//Append to list of all books
		for j := range b.Testaments[i].Books {
			b.BooksById[strings.ToLower(b.Testaments[i].Books[j].ID)] = len(b.Books)
			b.Books = append(b.Books, &b.Testaments[i].Books[j])
		}
	}
}

type Version struct {
	Title     string `xml:"title"json:"title"`
	ID        string `xml:"identifier"json:"id"`
	RefSystem string `xml:"refSystem"json:"refSystem"`
}

type Testament struct {
	Books []Book `xml:"div"`
}

type Book struct {
	Chs []*Chapter `xml:"chapter"`
	ID  string     `xml:"osisID,attr"`
}

type Chapter struct {
	Vrs []*Verse `xml:"verse"`
	ID  string   `xml:"osisWork,attr"`
}
