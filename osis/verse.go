package osis

import (
	"strconv"
	"strings"
)

type Verse struct {
	Text  string      `xml:",chardata"json:"text"`
	ID    string      `xml:"osisID,attr"json:"id"`
	Words []Word      `xml:"w"json:"words,omitempty"`
	Refs  []Reference `xml:"note>reference"json:"references,omitempty"`
}

type Reference struct {
	Text  string `xml:",chardata"json:"text"json:"comment"`
	RefID string `xml:"osisRef,attr"json:"refid"`
}

type VerseRef struct {
	BookID         string
	Chapter, Verse int
}

func NewVerseRef(ref string) (*VerseRef, error) {
	ref = strings.ToLower(ref)
	vr := VerseRef{}

	strs := strings.Split(ref, ".")
	if len(strs) < 3 {
		return nil, InvalidRef
	}

	vr.BookID = strs[0]

	var err error
	vr.Chapter, err = strconv.Atoi(strs[1])
	if err != nil {
		return nil, InvalidRef
	}
	//Human -> arr index
	vr.Chapter--

	vr.Verse, err = strconv.Atoi(strs[2])
	if err != nil {
		return nil, InvalidRef
	}
	//Human -> arr index
	vr.Verse--

	return &vr, nil
}

type Word struct {
	N          int    `xml:"n,attr"`
	Text       string `xml:",chardata"`
	Translit   string `xml:"xlit,attr"`
	Lemma      string `xml:"lemma,attr"`
	Morphology string `xml:"morph,attr"`
}
