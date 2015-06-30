package osis

import (
	"strconv"
	"strings"
)

type Verse struct {
	Text string      `xml:",chardata"`
	ID   string      `xml:"osisID,attr"`
	Refs []Reference `xml:"note>reference"`
}

type Reference struct {
	Text  string `xml:",chardata"`
	RefID string `xml:"osisRef,attr"`
}

type VerseRef struct {
	BookID         string
	Chapter, Verse int
}

func NewVerseRef(ref string) (*VerseRef, error) {
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
