package mysql

// TODO: Consistent variable naming

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"

	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	connectionInfo := fmt.Sprintf("%v:%v@tcp(%v:3306)/mobile_data?parseTime=true", user, pwd, host)
	db, err = sql.Open("mysql", connectionInfo)
	if err != nil {
		panic(err)
	}
}

// Insert inserts a given array of data structs (for example battery data) into
// the database
func Insert(ptrToData interface{}, args ...interface{}) error {
	// Get the value the ptrToData points to
	dataValue, err := getValueOfPtr(ptrToData)
	if err != nil {
		return err
	}

	// Get the matching query struct for the given data struct
	queryStruct, ok := queries[reflect.TypeOf(dataValue.Interface())]
	if !ok {
		return errors.New("The type of data you want to insert is not stored in the database.")
	}

	if dataValue.Kind() != reflect.Slice {
		return errors.New("Input data is not a slice.")
	}

	// Loop over the given array and insert its data into the database
	for i := 0; i < dataValue.Len(); i++ {
		// Get the value at index i and throw an error if it is not a struct
		v := dataValue.Index(i)
		if v.Kind() != reflect.Struct {
			return errors.New("Input data is not a slice of structs.")
		}

		// Create slice for the arguments to the database query
		data := make([]interface{}, v.NumField())

		// Transform the data struct into an array of interfaces so
		// we can pass that array as an argument to our prepared statement
		for j := 0; j < v.NumField(); j++ {
			data[j] = v.Field(j).Interface()
		}

		// Remove all potential zero values from the data slice and append the arguments
		// that were passed as arguments of the function
		data = append(removeZeroValues(data), args...)

		// Insert data into the table specific to the given data type
		if _, err = executeInsert(queryStruct.Insert, data...); err != nil {
			return err
		}
	}

	return nil
}

// Get takes a value and stores the result of the required query in the value
func Get(ptrToValue interface{}, args ...interface{}) error {
	// Get the value the ptrToValue points to
	value, err := getValueOfPtr(ptrToValue)
	if err != nil {
		return err
	}

	// Use the type of the value to retrieve its matching struct
	queryStruct, ok := queries[reflect.TypeOf(value.Interface())]
	if !ok {
		return errors.New("The type of data you want to retrieve is not stored in the database.")
	}

	// Query the database with the retrieve query for the given datatype
	rows, err := executeQuery(queryStruct.Retrieve, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get the type of the result struct
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

		// Append the scaned row to the result
		value.Set(reflect.Append(value, row))
	}

	return nil
}

// executeInsert executes an SQL insert statement
func executeInsert(insertStmt string, args ...interface{}) (sql.Result, error) {
	// Create prepared statement
	stmt, err := db.Prepare(insertStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute statement with arguments
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// executeQuery executes an SQL query
func executeQuery(query string, args ...interface{}) (rows *sql.Rows, err error) {
	// Create prepared statement
	stmt, err := db.Prepare(query)
	if err != nil {
		return rows, err
	}
	defer stmt.Close()

	// Query the database with statement
	rows, err = stmt.Query(args...)
	return rows, err
}

// getValueOfPtr gets the reflect value of a pointer
func getValueOfPtr(ptrToValue interface{}) (reflect.Value, error) {
	// Check if ptrToValue is really a pointer
	if reflect.ValueOf(ptrToValue).Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("Argument is not a pointer.")
	}

	// Get the value the ptrToValue points to
	value := reflect.Indirect(reflect.ValueOf(ptrToValue))

	return value, nil
}

// removeZeroValues removes all the values from the slice which are the default zero value of that type
func removeZeroValues(slice []interface{}) []interface{} {
	filteredSlice := slice[:0]
	for _, elem := range slice {
		zeroValue := reflect.Zero(reflect.TypeOf(elem)).Interface()
		if elem != zeroValue {
			filteredSlice = append(filteredSlice, elem)
		}
	}

	return filteredSlice
}
