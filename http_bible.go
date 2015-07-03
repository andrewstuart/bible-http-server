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
