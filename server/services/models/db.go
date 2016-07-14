package models

import (
	"github.com/jmoiron/sqlx"
	// mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

// Datastore is the interface for the database actions
type Datastore interface {
	GetUser(user User) (fullUser User, err error)
	CreateUser(user User) error
	UpdateUser(user User, field string) error
	DeleteUser(user User) error
	GetDevicesFromUser(user User) ([]Device, error)
	GetDeviceInfos(user User) ([]DeviceStored, error)
	CreateDeviceForUser(user User, device DeviceStored) error
	UpdateDevice(device DeviceStored, field string) error
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
