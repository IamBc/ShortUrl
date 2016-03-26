package main

import (
	"database/sql"
	"github.com/golang/glog"
	"os"
	"strings"
	_ "github.com/lib/pq"
	"expvar" 
)


var (
  dbCounters = expvar.NewMap("dbCounters")
  )

/*
Implementetaion of the storage_interface for relational databases
*/

var db *sql.DB //This is a connection pool. It must be global so it is visible in all files

/*
Must be called in the main function. It will create the nessecary environment for the storage.
*/
func InitStorage() {
	var err error
	//var tx    *sql.Tx
	db, err = sql.Open(os.Getenv("DB_CONNECTION_DRIVER"), os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		glog.Error(err)
	}
	if db == nil {
		glog.Fatal(db)
	}
}

func GetURLFromStorage(urlHash string) (string, error) {
	var err error

	tx, err := db.Begin()
	if err != nil {
	    glog.Error(err)
	    return ``, err
	} else {
	    glog.Info("ERROR ON BEGIN IS NiLL!!!")
	}

	dbCounters.Add(`queryCount`, 1)
	rows, err := tx.Query(`SELECT url FROM urls WHERE url_hash = $1`, urlHash)
	if err != nil {
		glog.Error(err)
		tx.Rollback()
		return ``, err
	}
	defer rows.Close()

	var url string
	for rows.Next() {
		err = rows.Scan(&url)
		if err != nil {
			tx.Rollback()
			return ``, err
		}
	}

	tx.Commit()
	return url, err
}

func AddURLToStorage(urlHash string, url string) (string, error) {
	var err error

	tx, err := db.Begin()

	dbCounters.Add(`queryCount`, 1)
	stmt, err := db.Prepare(`INSERT INTO urls(url_hash, url) VALUES($1, $2)`)
	if err != nil {
	    glog.Error("ERROR ON PREPARING FIRST STATEMENT")
	}
	_, err = tx.Stmt(stmt).Exec( urlHash, url)
        if err != nil {
                glog.Error(`ERROR IN FIRST QUERY: `, err)
                tx.Rollback()
        }
	// Hash already exists
	if err != nil && strings.ContainsAny(err.Error(), `Error 1062`) {
		dbCounters.Add(`queryCount`, 1)
		rows, err := tx.Query(`SELECT url_hash FROM urls WHERE url = $1`, url)
		if err != nil {
			glog.Error(err)
			tx.Rollback()
			return urlHash, err
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&urlHash)
			if err != nil {
				tx.Rollback()
				return ``, err
			}
		}

	} else if err != nil { // some exception has occured
		glog.Error(err)
		tx.Rollback()
	}

	tx.Commit()
	return urlHash, nil
}

func DeleteURL(urlHash string) error {
	var err error

	tx, err := db.Begin()
	dbCounters.Add(`queryCount`, 1)
	_, err = tx.Query(`DELETE FROM urls WHERE url_hash = $1`, urlHash)
	if err != nil {
		glog.Error(err)
		tx.Rollback()
	}
	tx.Commit()
	return err
}
