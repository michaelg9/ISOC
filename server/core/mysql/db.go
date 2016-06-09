package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"

	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/michaelg9/ISOC/server/services/models"
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

// GetUser gives back the user with a given username
func GetUser(username string) (models.User, error) {
	result, err := getData(getUser, models.User{}, username)
	if err != nil {
		// TODO: Check if error "sql: no rows in result set"
		return models.User{}, err
	}

	users, ok := result.([]models.User)
	if !ok {
		return models.User{}, errors.New("Failed to convert SQL result into type []models.User.")
	}

	if len(users) != 1 {
		return models.User{}, errors.New("Multiple users with same username.")
	}

	return users[0], nil
}

// IDEA: Factor all Get functions into one scan like function

// GetDevices gets all devices which are referenced to the given API key
func GetDevices(key string) ([]models.DeviceStored, error) {
	// Query the database
	result, err := getData(getDevices, models.DeviceStored{}, key)
	if err != nil {
		return nil, err
	}

	// Convert result to array of device structs
	devices, ok := result.([]models.DeviceStored)
	if !ok {
		return nil, errors.New("Failed to convert SQL result into type []models.DeviceStored.")
	}

	return devices, nil
}

// GetBattery gets all battery data associated with a device ID
func GetBattery(deviceID int) ([]models.Battery, error) {
	result, err := getData(getBattery, models.Battery{}, deviceID)
	if err != nil {
		return nil, err
	}

	battery, ok := result.([]models.Battery)
	if !ok {
		return nil, errors.New("Failed to convert SQL result into type []models.Battery.")
	}

	return battery, nil
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

// GetData queries the database with the given query and arguments. ResultStruct
// is the struct which has the format of one row of the resulting table of the query.
func getData(query string, resultStruct interface{}, args ...interface{}) (interface{}, error) {
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
