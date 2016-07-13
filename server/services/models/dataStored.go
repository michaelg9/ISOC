package models

import (
	"errors"
	"reflect"

	"github.com/fatih/structs"
)

// Battery is the struct for the battery
// percentage element
type Battery struct {
	Time  string `xml:"time,attr" json:"time" db:"timestamp"`
	Value int    `xml:",chardata" json:"value" db:"batteryPercentage"`
}

var typeName = map[reflect.Type]string{
	reflect.TypeOf([]Battery{}): "Battery",
}

var createQueries = map[string]string{
	"Battery": `INSERT INTO BatteryStatus (timestamp, batteryPercentage, device) VALUES (:Time, :Value, :ID);`,
}

var getQueries = map[string]string{
	"Battery": `SELECT timestamp, batteryPercentage FROM BatteryStatus WHERE device = :id;`,
}

// GetData gets the data from the given data type.
func (db *DB) GetData(device DeviceStored, ptrToData interface{}) error {
	value, err := getValueOfPtr(ptrToData)
	if err != nil {
		return err
	}
	query, ok := getQueries[typeName[value.Type()]]
	if !ok {
		return errors.New("Type of data is not stored.")
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.Select(value.Addr().Interface(), device)
}

// CreateData inserts the given data into the database.
func (db *DB) CreateData(device DeviceStored, ptrToData interface{}) error {
	value, err := getValueOfPtr(ptrToData)
	if err != nil {
		return err
	}
	query, ok := createQueries[typeName[value.Type()]]
	if !ok {
		return errors.New("Type of data is not stored.")
	}

	// TODO: check if value is slice
	for i := 0; i < value.Len(); i++ {
		data := value.Index(i).Interface()
		m := structs.Map(data)
		m["ID"] = device.ID
		_, err := db.NamedExec(query, m)
		if err != nil {
			return err
		}
	}

	return nil
}

// getValueOfPtr gets the reflect value of a pointer
func getValueOfPtr(ptrToValue interface{}) (reflect.Value, error) {
	ptr := reflect.ValueOf(ptrToValue)
	// Check if ptrToValue is really a pointer
	if ptr.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("Argument is not a pointer.")
	}

	// Get the value the ptrToValue points to
	value := reflect.Indirect(ptr)
	return value, nil
}
