package osis

import "encoding/xml"

type Bible struct {
	XMLName    xml.Name    `xml:"osis"`
	Testaments []Testament `xml:"osisText>div"`
}

type Testament struct {
	Books []Book `xml:"div"`
}

type Book struct {
	Chapters []Chapter `xml:"chapter"`
}

type Chapter struct {
	Verses []Verse `xml:"verse"`
}

type Verse string
