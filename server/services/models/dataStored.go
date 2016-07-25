package models

// TODO: Check everywhere for semicolons

import (
	"errors"
	"reflect"

	"github.com/fatih/structs"
)

// Battery is the struct for the battery
// percentage element
// TODO: Add charging attribute
type Battery struct {
	Time  string `xml:"time,attr" json:"time" db:"timestamp"`
	Value int    `xml:",chardata" json:"value" db:"batteryPercentage"`
}

// Call is the struct for a tracked Call
type Call struct {
	Type    string `xml:"type,attr" db:"callType"` // Type of call, either "Outgoing" or "Ingoing"
	Start   string `xml:"start,attr" db:"start"`
	End     string `xml:"end,attr" db:"end"`
	Contact string `xml:",chardate" db:"contact"`
}

// App is the struct for an installed application
type App struct {
	Name      string `xml:"name,attr" db:"name"` // Internal Name of the application
	UID       int    `xml:"uid,attr" db:"uid"`
	Version   string `xml:"version,attr" db:"version"`
	Installed string `xml:"installed,attr" db:"installed"`
	Label     string `xml:",chardata" db:"label"` // Label is the name that is visible to the user
}

// Runservice is the struct for when an app is used
type Runservice struct {
	AppName string `xml:",chardata" db:"appName"`
	RX      int    `xml:"rx,attr" db:"rx"` // Received bytes from this service
	TX      int    `xml:"tx,attr" db:"tx"` // Transmitted bytes from this service
	Start   string `xml:"start,attr" db:"start"`
	End     string `xml:"end,attr" db:"end"`
}

var typeName = map[reflect.Type]string{
	reflect.TypeOf([]Battery{}):    "Battery",
	reflect.TypeOf([]Call{}):       "Call",
	reflect.TypeOf([]App{}):        "App",
	reflect.TypeOf([]Runservice{}): "Runservice",
}

var createQueries = map[string]string{
	"Battery":    `INSERT INTO BatteryStatus (timestamp, batteryPercentage, device) VALUES (:Time, :Value, :ID);`,
	"Call":       `INSERT INTO ` + "`" + `Call` + "`" + ` (callType, start, end, contact, device) VALUES (:Type, :Start, :End, :Contact, :ID);`,
	"App":        `INSERT INTO App (name, uid, version, installed, label, device) VALUES (:Name, :UID, :Version, :Installed, :Label, :ID);`,
	"Runservice": `INSERT INTO Runservice (appName, rx, tx, start, end, device) VALUES (:AppName, :RX, :TX, :Start, :End, :ID);`,
}

var getQueries = map[string]string{
	"Battery":    `SELECT timestamp, batteryPercentage FROM BatteryStatus WHERE device = :id;`,
	"Call":       `SELECT callType, start, end, contact FROM ` + "`" + `Call` + "`" + ` WHERE device = :id;`,
	"App":        `SELECT name, uid, version, installed, label FROM App WHERE device = :id;`,
	"Runservice": `SELECT appName, rx, tx, start, end FROM Runservice WHERE device = :id;`,
}

// GetData gets the data from the given data type.
func (db *DB) GetData(device DeviceStored, ptrToData interface{}) error {
	value, err := getValueOfPtr(ptrToData)
	if err != nil {
		return err
	}
	query, ok := getQueries[typeName[value.Type()]]
	if !ok {
		return errors.New("Type of data is not stored.") // TODO: error as string
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}

	return stmt.Select(value.Addr().Interface(), device)
}

// CreateData inserts the given data into the database.
func (db *DB) CreateData(device DeviceStored, ptrToData interface{}) error {
	// Get the reflect value of the given data
	value, err := getValueOfPtr(ptrToData)
	if err != nil {
		return err
	}

	// Get the right query for that type of data
	query, ok := createQueries[typeName[value.Type()]]
	if !ok {
		return errors.New("Type of data is not stored.")
	}

	// Insert each data point in the slice into the database
	for i := 0; i < value.Len(); i++ {
		data := value.Index(i).Interface()
		m := structs.Map(data)           // Transform the struct of the type into map[string]string
		m["ID"] = device.ID              // Append the device ID to that map
		_, err := db.NamedExec(query, m) // Insert the data
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
