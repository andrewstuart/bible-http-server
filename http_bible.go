//Package main is an http api for bible verses with a postgres backend.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//IDQuery gets a verse by id
const idQuery = `
SELECT vv.text, vv.verseid, v.book, v.chapter, v.verse 
FROM verse_version vv 
INNER JOIN verse v 
	ON v.id = vv.verseid 
WHERE vv.verseid = $1`

//GetVerseByID gets a verse by id.
func GetVerseByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["verseId"]

	if id == "" {
		w.WriteHeader(404)
		return
	}

	v := VerseResult{}

	err := db.QueryRow(idQuery, id).Scan(&v.Verse.Text, &v.Verse.ID, &v.Book, &v.Chapter, &v.VerseNum)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(v)
}

const versesQueryBase = `
SELECT vv.text, ver.name, v.book, v.chapter, v.verse, v.id
FROM verse_version vv
INNER JOIN verse v
	ON v.id = vv.verseid
INNER JOIN version ver
	ON ver.id = vv.versionid
WHERE `

//GetVerse returns json for a verse by query
func GetVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//Dynamic WHERE, still sql-safe
	conditions := make([]string, 0, 3)
	params := make([]interface{}, 0, 0)

	for _, param := range []string{"book", "chapter", "verse"} {
		if p, ok := vars[param]; ok {
			c := fmt.Sprintf(" v.%s = $%d", param, len(params)+1)
			conditions = append(conditions, c)
			params = append(params, p)
		}
	}

	if p, ok := vars["version"]; ok {
		c := fmt.Sprintf(" lower(ver.extid) = lower($%d)", len(params)+1)
		conditions = append(conditions, c)
		if strings.Index(p, "Bible.") == -1 {
			p = "Bible." + p
		}
		params = append(params, p)
	}

	query := versesQueryBase + strings.Join(conditions, " AND ")
	rows, err := db.Query(query, params...)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	verses := make([]VerseResult, 0, 0)

	for rows.Next() {
		v := VerseResult{}
		err := rows.Scan(&v.Verse.Text, &v.VersionId, &v.Book, &v.Chapter, &v.VerseNum, &v.Verse.ID)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		verses = append(verses, v)
	}

	enc := json.NewEncoder(w)
	enc.Encode(verses)
}

const bookQ = `
SELECT id, name, ord, chaps
FROM v_book_counts
ORDER BY ord asc`

func getBooks(w http.ResponseWriter, r *http.Request) {
	books := make([]book, 0, 66)

	rows, err := db.Query(bookQ)
	if err != nil {
		log.Println("Error getting books", err)
		w.WriteHeader(500)
		return
	}

	for rows.Next() {
		b := book{}
		err := rows.Scan(&b.ID, &b.Name, &b.Order, &b.ChapterCount)
		if err != nil {
			log.Println("Error scanning book from db", err)
			w.WriteHeader(500)
			return
		}

		books = append(books, b)
	}

	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Println("Error encoding books to JSON.", err)
	}
}

type book struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	Order        int    `json:"order"`
	ChapterCount int    `json:"chapterCount"`
}
