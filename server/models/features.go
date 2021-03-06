package models

import (
	"errors"
	"reflect"

	"github.com/fatih/structs"
)

var (
	errFeatureNotStored = errors.New("Type of data is not stored.")
)

// Features contains all the data that we store
type Features struct {
	Battery    []Battery    `xml:"battery,omitempty" json:"battery,omitempty"`
	Call       []Call       `xml:"call,omitempty" json:"call,omitempty"`
	App        []App        `xml:"app,omitempty" json:"app,omitempty"`
	Runservice []Runservice `xml:"runservice,omitempty" json:"runservice,omitempty"`
}

// GetContents returns a slice of pointers to all the data of the device in the struct.
// We need to return pointers in order to preserve the underlying type.
func (deviceData *Features) GetContents() []interface{} {
	v := reflect.Indirect(reflect.ValueOf(deviceData))
	contents := make([]interface{}, v.NumField())

	for i := range contents {
		contents[i] = v.Field(i).Addr().Interface()
	}

	return contents
}

// Battery is the struct for the battery
// percentage element
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

// typeName maps the type of a slice of data to a key. The key can be used
// to retrieve queries from the createQueries and getQueries maps.
var typeName = map[reflect.Type]string{
	reflect.TypeOf([]Battery{}):    "Battery",
	reflect.TypeOf([]Call{}):       "Call",
	reflect.TypeOf([]App{}):        "App",
	reflect.TypeOf([]Runservice{}): "Runservice",
}

// createQueries maps the keys for the feature types to INSERT queries.
var createQueries = map[string]string{
	"Battery":    `INSERT INTO BatteryStatus (timestamp, batteryPercentage, device) VALUES (:Time, :Value, :ID);`,
	"Call":       `INSERT INTO ` + "`" + `Call` + "`" + ` (callType, start, end, contact, device) VALUES (:Type, :Start, :End, :Contact, :ID);`,
	"App":        `INSERT INTO App (name, uid, version, installed, label, device) VALUES (:Name, :UID, :Version, :Installed, :Label, :ID);`,
	"Runservice": `INSERT INTO Runservice (appName, rx, tx, start, end, device) VALUES (:AppName, :RX, :TX, :Start, :End, :ID);`,
}

// getQueries maps the key for the feature types to SELECT queries.
var getQueries = map[string]string{
	"Battery":    `SELECT timestamp, batteryPercentage FROM BatteryStatus`,
	"Call":       `SELECT callType, start, end, contact FROM ` + "`" + `Call` + "`",
	"App":        `SELECT name, uid, version, installed, label FROM App`,
	"Runservice": `SELECT appName, rx, tx, start, end FROM Runservice`,
}

// GetAllFeatureData gets all entries from a given feature (e.g. battery logs from all devices).
func (db *DB) GetAllFeatureData(ptrToFeature interface{}) error {
	// Get the reflect.Value to which the pointer points
	value, err := getValueOfPtr(ptrToFeature)
	if err != nil {
		return err
	}

	// Retrieve the SELECT query for the given feature type
	query, ok := getQueries[typeName[value.Type()]]
	if !ok {
		return errFeatureNotStored
	}

	// Query the database and return result
	return db.Select(ptrToFeature, query)
}

// GetFeatureOfDevice gets the data from the given data type.
func (db *DB) GetFeatureOfDevice(aboutDevice AboutDevice, ptrToFeature interface{}) error {
	// Get the reflect.Value to which the pointer points
	value, err := getValueOfPtr(ptrToFeature)
	if err != nil {
		return err
	}

	// Retrieve the SELECT query for the given feature type
	query, ok := getQueries[typeName[value.Type()]]
	if !ok {
		return errFeatureNotStored
	}

	stmt, err := db.PrepareNamed(query + ` WHERE device = :id`)
	if err != nil {
		return err
	}

	// Query the database and return result
	return stmt.Select(ptrToFeature, aboutDevice)
}

// CreateFeatureForDevice inserts the given data into the database.
func (db *DB) CreateFeatureForDevice(aboutDevice AboutDevice, ptrToFeature interface{}) error {
	// Get the reflect.Value to which the pointer points
	value, err := getValueOfPtr(ptrToFeature)
	if err != nil {
		return err
	}

	// Retrieve the INSERT query for the given feature type
	query, ok := createQueries[typeName[value.Type()]]
	if !ok {
		return errFeatureNotStored
	}

	// Insert each data point in the slice into the database
	for i := 0; i < value.Len(); i++ {
		data := value.Index(i).Interface()
		m := structs.Map(data)           // Transform the struct of the type into map[string]string
		m["ID"] = aboutDevice.ID         // Append the device ID to that map
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
