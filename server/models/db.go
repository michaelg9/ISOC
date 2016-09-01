package models

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

// Datastore is the interface for the database actions
type Datastore interface {
	// User methods
	GetUser(user User) (fullUser User, err error)
	CreateUser(user User) (insertedID int, err error)
	UpdateUser(user User) error
	DeleteUser(user User) error

	// Device methods
	GetAllDevices() ([]AboutDevice, error)
	GetDeviceFromUser(user User, device Device) (fullDevice Device, err error)
	GetDevicesFromUser(user User) ([]Device, error)
	CreateDeviceForUser(user User, aboutDevice AboutDevice) (insertedID int, err error)
	UpdateDevice(aboutDevice AboutDevice) error
	DeleteDevice(aboutDevice AboutDevice) error

	// Feature methods
	GetAllFeatureData(ptrToFeature interface{}) error
	GetFeatureOfDevice(aboutDevice AboutDevice, ptrToFeature interface{}) error
	CreateFeatureForDevice(aboutDevice AboutDevice, ptrToFeature interface{}) error
}

// DB is the database struct
type DB struct {
	*sqlx.DB
}

// NewDB returns a new database connection
func NewDB(dsn string) *DB {
	// Connect to database
	db := sqlx.MustConnect("mysql", dsn)
	return &DB{db}
}

// update takes a struct and the a map from the field names to database queries for updating
// the value stored in the field. If a value in the given struct is non-empty and there is an
// update query for the field stored in the map, we update the database.
func (db *DB) update(queries map[string]string, arg interface{}) error {
	value := reflect.ValueOf(arg)

	if value.Kind() != reflect.Struct {
		return errors.New("Need to pass a struct to update.")
	}

	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i)
		fieldZeroValue := reflect.Zero(fieldValue.Type())
		query, queryExists := queries[fieldName]

		fieldNotEmpty := fieldValue.Interface() != fieldZeroValue.Interface()

		if fieldNotEmpty && queryExists {
			_, err := db.NamedExec(query, arg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
