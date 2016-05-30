package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	user = "treigerm"
	pwd  = "Hip$terSWAG"

	// Gets the hash of the password given the username
	passwordQuery = "SELECT passwordHash FROM User WHERE username = ?"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@/mobile_data?parseTime=true", user, pwd))
	if err != nil {
		// TODO: Error handling
		panic(err)
	}
}

// GetHashedPassword gets the hash of the password from a given user
func GetHashedPassword(username string) (hash string) {
	stmt, err := db.Prepare(passwordQuery)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&hash)
	if err != nil {
		// TODO: Error handling
		panic(err)
	}

	return
}
