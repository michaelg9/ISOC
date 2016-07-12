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
	UpdateUser(email, field string, value interface{}) error
	DeleteUser(user User) error
	GetDevicesFromUser(user User) ([]DeviceStored, error)
	CreateDeviceForUser(user User, device DeviceStored) error
	UpdateDevice(id int, field string, value interface{}) error
	DeleteDevice(device DeviceStored) error
	GetBattery(device DeviceStored, batteries *[]Battery) error
	CreateBattery(data DeviceStored, batteries *[]Battery) error
}

// DB is the database struct
type DB struct {
	*sqlx.DB
}

// NewDB returns a new database connection
func NewDB(dsn string) (*DB, error) {
	// NOTE: Might change to MustConnect
	// Connect to database
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
