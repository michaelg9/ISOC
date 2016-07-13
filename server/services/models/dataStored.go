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

// GetBattery gets all the battery data from the given device and scans it into the given slice
func (db *DB) GetBattery(device DeviceStored, batteries *[]Battery) error {
	getBatteryQuery := `SELECT timestamp, batteryPercentage FROM BatteryStatus WHERE device = :id;`
	stmt, err := db.PrepareNamed(getBatteryQuery)
	if err != nil {
		return err
	}

	return stmt.Select(batteries, device)
}

// CreateBattery inserts the given data into the database
func (db *DB) CreateBattery(device DeviceStored, batteries *[]Battery) error {
	createBatteryQuery := "INSERT INTO BatteryStatus (timestamp, batteryPercentage, device) VALUES (?, ?, ?);"
	for _, battery := range *batteries {
		_, err := db.Exec(createBatteryQuery, battery.Time, battery.Value, device.ID)
		if err != nil {
			return err
		}
	}

	return nil
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
