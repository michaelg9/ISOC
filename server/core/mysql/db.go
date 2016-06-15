package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

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

// InsertData inserts a given array of data structs (for example battery data) into
// the database
func InsertData(deviceID int, data interface{}) error {
	// Get the matching query struct for the given data struct
	queryStruct, ok := queries[reflect.TypeOf(data)]
	if !ok {
		return errors.New("This type of input data is not stored in the database.")
	}

	// Prepared statement to insert in the "parent" data table
	stmtParent, err := db.Prepare(insertData)
	if err != nil {
		return err
	}
	defer stmtParent.Close()

	// Prepared statement to insert into specific data table
	stmtChild, err := db.Prepare(queryStruct.Insert)
	if err != nil {
		return err
	}
	defer stmtChild.Close()

	// Get a the value the given data
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return errors.New("Input data is not a slice.")
	}

	// Loop over the given array and insert its data into the database
	for i := 0; i < dataValue.Len(); i++ {
		v := dataValue.Index(i)
		if v.Kind() != reflect.Struct {
			return errors.New("Input data is not a slice of structs.")
		}

		args := make([]interface{}, v.NumField())

		// Transform the data struct into an array of interfaces so
		// we can pass that array as an argument to our prepared statement
		for j := 0; j < v.NumField(); j++ {
			args[j] = v.Field(j).Interface()
		}

		// Check if the first contains a timestamp that
		// is in the format specified in the models package
		_, err = time.Parse(timeLayout, args[0].(string))
		if err != nil {
			return err
		}

		timestamp := args[0]
		result, err := stmtParent.Exec(deviceID, timestamp)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// args[0] was the timestamp before but now we don't need it
		// anymore. Instead we replace it with the inserted ID of the
		// "parent" table so we can link both entries later on.
		args[0] = id
		_, err = stmtChild.Exec(args...)
		if err != nil {
			return err
		}
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
	queryStruct, ok := queries[reflect.TypeOf(value.Interface())]
	if !ok {
		return errors.New("The type of data you want to retrieve is not stored in the database.")
	}

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

	// Get the type if the result struct
	typeOfResult := reflect.TypeOf(queryStruct.StoredData)
	for rows.Next() {
		// Create a struct of the result struct type
		row := reflect.New(typeOfResult).Elem()
		// Row is the slice which stores the pointers to the fields of row
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
