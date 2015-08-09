package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andrewstuart/bible-http-server/osis"
)

func SearchVerse(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	log.Printf("Search for verse term: %s\n", q)

	if q != "" {
		verses, err := search(q)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("Error searching %s: %v\n", q, err)
			fmt.Fprintf(w, "Error searching: %v", err)
			return
		}

		if len(verses) == 0 {
			w.WriteHeader(404)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(verses)
		return
	}

	w.WriteHeader(404)
}

//VerseResult is simply a representation of a Verse with a few additional
//desired items that couldn't be inferred from context
type VerseResult struct {
	Book      string  `json:"book"`
	Chapter   int     `json:"chapter"`
	VerseNum  int     `json:"verse"`
	VersionId string  `json:"version,omitempty"`
	Match     float64 `json:"match,omitempty"`
	osis.Verse
}

const VerseQuery = `
SELECT v.book, v.chapter, v.verse, vv.verseid, vv.text, ver.name, ts_rank(vect, q)
FROM to_tsquery($1) q, verse_version vv 
INNER JOIN verse v
	ON v.id = vv.verseid
INNER JOIN version ver
	ON ver.id = vv.versionid
WHERE vv.vect @@ q
`

func search(str string) ([]VerseResult, error) {
	verseCurs, err := db.Query(VerseQuery, str)

	if err != nil {
		return nil, err
	}

	verses := make([]VerseResult, 0, 5)
	for verseCurs.Next() {
		v := VerseResult{}
		verseCurs.Scan(&v.Book, &v.Chapter, &v.VerseNum, &v.Verse.ID, &v.Verse.Text, &v.VersionId, &v.Match)

		verses = append(verses, v)
	}

	return verses, nil
}
