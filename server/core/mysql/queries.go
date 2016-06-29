package mysql

import (
	"reflect"

	"github.com/michaelg9/ISOC/server/services/models"
)

const (
	timeLayout    = "2006-01-02 15:04:05"
	insertData    = "INSERT INTO Data (device, timestamp) VALUES (?, ?);"
	insertBattery = "INSERT INTO BatteryStatus VALUES (?, ?);"
	insertUser    = "INSERT INTO User (email, passwordHash) VALUES (?, ?)"
	getUser       = "SELECT uid, email, passwordHash, COALESCE(apiKey, '') AS apiKey FROM User WHERE email = ? OR apiKey = ?;"
	getBattery    = "SELECT d.timestamp, b.batteryPercentage " +
		"FROM Data d, BatteryStatus b " +
		"WHERE d.device = ? AND d.id = b.id;"
	getDevices = "SELECT dev.id, dev.manufacturer, dev.modelName, dev.osVersion " +
		"FROM Device dev, User u " +
		"WHERE u.email = ? AND u.uid = dev.user;"
)

// QueryStruct saves insert and retrieve queries for some data that is stored.
// For example, the battery data.
type QueryStruct struct {
	StoredData       interface{}
	Retrieve, Insert string
}

// queries stores the type of a given slice of structs as a key
// The value is the QueryStruct for the type of struct that is stored
// in the slice of the key.
// TODO: Explain use case
var queries = map[reflect.Type]QueryStruct{
	reflect.TypeOf([]models.Battery{}):      QueryStruct{models.Battery{}, getBattery, insertBattery},
	reflect.TypeOf([]models.DeviceStored{}): QueryStruct{models.DeviceStored{}, getDevices, ""},
	reflect.TypeOf([]models.User{}):         QueryStruct{models.User{}, getUser, insertUser},
}
