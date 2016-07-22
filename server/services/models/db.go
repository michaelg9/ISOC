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
	GetUser(user User) (fullUser User, err error)
	CreateUser(user User) error
	UpdateUser(user User) error
	DeleteUser(user User) error
	GetDevicesFromUser(user User) ([]Device, error)
	GetDeviceInfos(user User) ([]DeviceStored, error)
	CreateDeviceForUser(user User, device DeviceStored) error
	UpdateDevice(device DeviceStored) error
	DeleteDevice(device DeviceStored) error
	GetData(device DeviceStored, ptrToData interface{}) error
	CreateData(device DeviceStored, ptrToData interface{}) error
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
