package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
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
