package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"

	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	user = "treigerm"
	pwd  = "123"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:3306)/mobile_data?parseTime=true", user, pwd, os.Getenv("DB_HOST")))
	if err != nil {
		panic(err)
	}
}

// GetHashedPassword gets the hash of the password from a given user
func GetHashedPassword(username string) (hash string, err error) {
	stmt, err := db.Prepare(passwordQuery)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&hash)
	if err != nil {
		return "", err
	}

	return
}

// InsertBatteryData inserts the given data for the battery status
// TODO: Change timestamp type to time.Time
func InsertBatteryData(deviceID, batteryStatus int, timestamp string) (err error) {
	stmt, err := db.Prepare(insertIntoData)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(deviceID, timestamp)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(insertIntoBattery)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, batteryStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetData queries the database with the given query and arguments. ResultStruct
// is the struct which has the format of one row of the resulting table of the query.
func GetData(query string, resultStruct interface{}, args ...interface{}) (interface{}, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get the type if the result struct and create a new slice of that type
	typeOfResult := reflect.TypeOf(resultStruct)
	result := reflect.Zero(reflect.SliceOf(typeOfResult))
	for rows.Next() {
		// Create a struct of the result struct type
		row := reflect.New(typeOfResult).Elem()
		// Row is the slice which stores the pointers to the fields of r
		rowPointers := make([]interface{}, row.NumField())
		for i := range rowPointers {
			rowPointers[i] = row.Field(i).Addr().Interface()
		}

		// Save one result row to the the row struct
		err = rows.Scan(rowPointers...)
		if err != nil {
			return nil, err
		}

		result = reflect.Append(result, row)
	}

	return result.Interface(), nil
}
