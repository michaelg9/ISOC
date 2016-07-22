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

func (db *DB) update(queries map[string]string, arg interface{}) error {
	value := reflect.ValueOf(arg)

	// Return an error if arg is not a struct
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
