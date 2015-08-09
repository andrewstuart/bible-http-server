package main

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/andrewstuart/bible-http-server/osis"
	_ "github.com/lib/pq"
)

const (
	versionInsert = `INSERT INTO version (extid, name) VALUES ($1, $2) RETURNING id`
)

var db *sql.DB

const passEnvName = "PGPASSWORD"

func init() {
	var err error

	pass := strings.TrimSpace(os.Getenv(passEnvName))
	if pass == "" {
		log.Fatalf("Please set %s environment variable with postgres password.", passEnvName)
	}

	dbHost := stringDef(getLinkedPort(), "localhost")

	dbConn := fmt.Sprintf("postgres://bible:%s@%s/bible?sslmode=disable", pass, dbHost)

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

	var versionID int
	err = tx.QueryRow(versionInsert, b.Version.ID, b.Version.Title).Scan(&versionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	wg := &sync.WaitGroup{}

	for i, bk := range b.Books {
		wg.Add(1)
		go func(i int, bk *osis.Book) {
			defer wg.Done()

			for j, ch := range bk.Chs {
				for k, vs := range ch.Vrs {

					//Handle words (osis uses for greek/hebrew)
					if len(vs.Words) != 0 {
						txt := make([]string, len(vs.Words))
						for i := range vs.Words {
							txt[i] = vs.Words[i].Text
						}

						vs.Text = strings.Join(txt, " ")
					}

					var verseID int
					err = tx.QueryRow(`SELECT id FROM verse  where book = $1 and chapter = $2 and verse = $3`, bk.ID, j+1, k+1).Scan(&verseID)

					if err != nil {
						log.Printf("Could not find %s %d:%d, inserting.\n", bk.ID, j+1, k+1)
						err = tx.QueryRow(`INSERT INTO verse (book, chapter, verse) values ($1, $2, $3) RETURNING id`, bk.ID, j+1, k+1).Scan(&verseID)
						continue
					}

					_, err = tx.Exec(`INSERT INTO verse_version (versionId, verseId, text) values ($1, $2, $3)`, versionID, verseID, vs.Text)
					if err != nil {
						return
					}
				}
			}
		}(i, bk)
	}

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

func getLinkedPort() string {
	e := os.Getenv("POSTGRES_PORT")
	if e == "" {
		return ""
	}

	vals := strings.Split(e, "://")
	if len(vals) < 2 {
		return ""
	}

	return vals[1]
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
