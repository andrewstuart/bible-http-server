package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func SetHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	}
}

const IdQuery = `
SELECT vv.text, vv.verseid, v.book, v.chapter, v.verse 
FROM verse_version vv 
INNER JOIN verse v 
	ON v.id = vv.verseid 
WHERE vv.verseid = $1`

func GetVerseById(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["verseId"]

	if id == "" {
		w.WriteHeader(404)
		return
	}

	v := VerseResult{}

	err := db.QueryRow(IdQuery, id).Scan(&v.Verse.Text, &v.Verse.ID, &v.Book, &v.Chapter, &v.VerseNum)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	enc := json.NewEncoder(w)
	enc.Encode(v)
}

const VersesQueryBase = `
SELECT vv.text, ver.name, v.book, v.chapter, v.verse, v.id
FROM verse_version vv
INNER JOIN verse v
	ON v.id = vv.verseid
INNER JOIN version ver
	ON ver.id = vv.versionid
WHERE `

func GetVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	//Dynamic WHERE, still sql-safe
	conditions := make([]string, 0, 3)
	params := make([]interface{}, 0)

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

	query := VersesQueryBase + strings.Join(conditions, " AND ")
	rows, err := db.Query(query, params...)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	verses := make([]VerseResult, 0)

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
SELECT id, name, ord
FROM book
ORDER BY ord asc`

func getBooks(w http.ResponseWriter, r *http.Request) {
	books := make([]Book, 0, 66)

	rows, err := db.Query(bookQ)
	if err != nil {
		log.Println("Error getting books", err)
		w.WriteHeader(500)
		return
	}

	for rows.Next() {
		b := Book{}
		err := rows.Scan(&b.ID, &b.Name, &b.Order)
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

type Book struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Order int    `json:"order"`
}
