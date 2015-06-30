package osis

import (
	"encoding/xml"
	"fmt"
)

type Bible struct {
	Books      []*Book        `json:"-"xml:"-"`
	BooksById  map[string]int `json:"-"xml:"-"`
	XMLName    xml.Name       `xml:"osis"json:"-"`
	Testaments []Testament    `xml:"osisText>div"json:"testaments"`
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
			b.BooksById[b.Testaments[i].Books[j].ID] = len(b.Books)
			b.Books = append(b.Books, &b.Testaments[i].Books[j])
		}
	}
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
	ID  string   `xml:"osisID,attr"`
}
