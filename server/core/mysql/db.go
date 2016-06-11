package mysql

import (
	"database/sql"
	"errors"
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

// InsertBatteryData inserts the given data for the battery status
// TODO: Refactor into more general InsertData
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

// Get takes a value and stores the result of the required query in the value
func Get(ptrToValue interface{}, args ...interface{}) error {
	// Check if ptrToValue is really a pointer
	if reflect.ValueOf(ptrToValue).Kind() != reflect.Ptr {
		return errors.New("Argument is not a pointer.")
	}

	// Get the value the ptrToValue points to
	value := reflect.Indirect(reflect.ValueOf(ptrToValue))
	// Use the type of the value to retrieve its matching struct
	queryStruct := queries[reflect.TypeOf(value.Interface())]

	stmt, err := db.Prepare(queryStruct.Retrieve)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get the type if the result struct and create a new slice of that type
	typeOfResult := reflect.TypeOf(queryStruct.StoredData)
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
			return err
		}

		value.Set(reflect.Append(value, row))
	}

	return nil
}
