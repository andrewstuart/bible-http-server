package main

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andrewstuart/bible/osis"
	_ "github.com/lib/pq"
)

const (
	versionInsert = `INSERT INTO version (extid, name) VALUES ($1, $2) RETURNING id`
	verseInsert   = `INSERT INTO book (id, name, ord) VALUES ($1, $1, $2)`
)

var db *sql.DB

const PassEnvName = "PGPASSWORD"

func init() {
	var err error

	pass := strings.TrimSpace(os.Getenv(PassEnvName))
	if pass == "" {
		log.Fatalf("Please set %s environment variable with postgres password.", PassEnvName)
	}

	dbConn := fmt.Sprintf("postgres://bible:%s@localhost/bible?sslmode=disable", pass)

	db, err = sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal(err)
	}
}

func store(b *osis.Bible) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var versionId int
	err = tx.QueryRow(versionInsert, b.Version.ID, b.Version.Title).Scan(&versionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	wg := &sync.WaitGroup{}

	for i, bk := range b.Books {
		go func(i int, bk *osis.Book) {
			wg.Add(1)
			_, err = tx.Exec(verseInsert, bk.ID, i)
			if err != nil {
				return
			}

			for j, ch := range bk.Chs {
				for k, vs := range ch.Vrs {
					var verseId int
					err := tx.QueryRow(`INSERT INTO verse (book, chapter, verse) VALUES ($1, $2, $3) RETURNING id`, bk.ID, j+1, k+1).Scan(&verseId)
					if err != nil {
						return
					}

					log.Println(verseId)

					_, err = tx.Exec(`INSERT INTO verse_version (versionId, verseId, text) values ($1, $2, $3)`, versionId, verseId, vs.Text)
					if err != nil {
						return
					}
				}
			}
			wg.Done()
		}(i, bk)
	}

	time.Sleep(time.Second)

	wg.Wait()

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func loadFromGzippedFile(path string) (*osis.Bible, error) {
	zipped, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(zipped)
	if err != nil {
		return nil, err
	}

	return osis.NewBible(r)
}
