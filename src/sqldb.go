package main

import (
	"database/sql"
	"fmt"
	"log"

	// Using a CGO_ free version of the driver makes the docker image much smaller,
	// this is completly fine since we have a super simple application
	_ "modernc.org/sqlite"
)


func openDB() *sql.DB {
	// WAL for better read concurrency; wait up to 5s on locks
	// DSN format supports _pragma=... with modernc.org/sqlite
	// To be honest for this simple usecase this is a bit overkill
	dsn := "file:/app/data/app.db?cache=shared&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil { log.Fatalf("open: %v", err) }

	// SQLite is single-writer;
	// This way in the future there is no way for me to do bad things
	db.SetMaxOpenConns(1)

	return db
}

type sqlDb struct {
	db *sql.DB
}

func (db *sqlDb) closeSqlDb() {
  err := db.db.Close() 
  if err != nil { log.Fatal("Could not close db connection")}
  log.Println("db Closed")
}

func startSqlDb() (sqlDb,error) {
	db := openDB()
  _,err := db.Exec(`
  		CREATE TABLE IF NOT EXISTS urls(
      id INTEGER PRIMARY KEY,
			sUrl TEXT NOT NULL,
			url TEXT NOT NULL
		);`)
  if err != nil { 
    log.Println("ERROR: Could not exec setup call to db")
    return sqlDb{},err
  }
  return sqlDb{db: db},nil
}

func (db *sqlDb) CheckUrl(url string) (string, error) {
    var sUrl string
    err := db.db.QueryRow(
        `SELECT sUrl FROM urls WHERE url = ? LIMIT 1`,
        url,
    ).Scan(&sUrl)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        log.Println("ERROR: Could not checkUrl:", err)
        return "", err
    }
    return sUrl, nil
}


func (db *sqlDb) WriteUrl(url, shortUrl string) error {
    _, err := db.db.Exec(
        "INSERT INTO urls(sUrl, url) VALUES(?, ?)",
        shortUrl,
        url,
    )
    if err != nil {
        log.Println("ERROR: Could not write to db:", err)
        return err
    }
    return nil
}

func (db *sqlDb) CheckForCollision(shortUrl *string) error {
    var exists int
    err := db.db.QueryRow(
        `SELECT EXISTS(SELECT 1 FROM urls WHERE sUrl = ?)`,
        *shortUrl,
    ).Scan(&exists)
    if err != nil {
        return err // some database error
    }
    if exists == 1 {
        return fmt.Errorf("shortUrl '%s' already exists", *shortUrl)
    }
    return nil // no collision
}
func (db *sqlDb) CheckSurl(surl string) (string, error) {
    var url string
    err := db.db.QueryRow(
        `SELECT url FROM urls WHERE sUrl = ? LIMIT 1`,
        surl,
    ).Scan(&url)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", nil
        }
        log.Println("ERROR: Could not checkSurl:", err)
        return "", err
    }
    return url, nil
}
