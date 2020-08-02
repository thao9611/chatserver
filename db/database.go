package db

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbname    = flag.String("db", "chatdb", "db name")
	tablename = flag.String("table", "", "table name")
)

// CreateDB ..
func CreateDB(namedb string) {
	db, err := sql.Open("mysql", "thaovu:password@tcp(localhost:3306)/"+namedb)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	} else {
		log.Println("Ping successfully")
	}
	_, err = db.Exec("USE " + namedb)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE " + namedb)
	if err != nil {
		panic(err)
	}
	log.Printf("Created db %s", namedb)
}

// CreateTable ..
func CreateTable(db *sql.DB, tablename string) {
	cmd := fmt.Sprintf("CREATE TABLE %s (user VARCHAR(1000), text VARCHAR(1000), room VARCHAR(1000), date TIMESTAMP)", tablename)
	_, err := db.Exec(cmd)
	if err != nil {
		panic(err)
	}
	log.Println("Create table successfully")
}

/*
	stmtIns, err := db.Prepare("INSERT INTO chat1 VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	stmtOut, err := db.Prepare("SELECT text FROM chat1 WHERE user = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	user := "thao"
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("hello %d", i)
		log.Printf("Insert %s - %s\n", user, text)
		_, err = stmtIns.Exec(user, text) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	var textQuery []byte
	rows, err := stmtOut.Query(user)
	for rows.Next() {
		err := rows.Scan(&textQuery) // WHERE number = 13
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		fmt.Printf("Text from thao: %s\n", textQuery)
	}
*/
