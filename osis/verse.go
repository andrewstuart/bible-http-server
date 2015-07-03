package osis

import (
	"strconv"
	"strings"
)

type Verse struct {
	Book    string      `json:"book"`
	Chapter string      `json:"chapter"`
	Verse   string      `json:"verse"`
	Text    string      `xml:",chardata"json:"text"`
	ID      string      `xml:"osisID,attr"json:"id"`
	Refs    []Reference `xml:"note>reference"json:"references"`
}

type Reference struct {
	Text  string `xml:",chardata"json:"text"`
	RefID string `xml:"osisRef,attr"`
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
