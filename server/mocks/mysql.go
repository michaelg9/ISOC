package mocks

import (
	"database/sql"
	"reflect"

	"github.com/michaelg9/ISOC/server/models"
)

// MockDB implements the Datastore interface so we can use it to mock the DB
type MockDB struct{}

func (mdb *MockDB) GetUser(user models.User) (models.User, error) {
	if user.Email == Users[0].Email {
		return Users[0], nil
	}
	return models.User{}, sql.ErrNoRows
}

func (mdb *MockDB) CreateUser(user models.User) error {
	return nil
}

func (mdb *MockDB) UpdateUser(user models.User) error {
	return nil
}

func (mdb *MockDB) DeleteUser(user models.User) error {
	return nil
}

func (mdb *MockDB) GetDevice(device models.Device) (models.Device, error) {
	if device.AboutDevice.ID == 1 {
		return Devices[0], nil
	}
	return models.Device{}, sql.ErrNoRows
}

func (mdb *MockDB) GetDevicesFromUser(user models.User) ([]models.Device, error) {
	return Devices[:1], nil
}

func (mdb *MockDB) GetDeviceInfos(user models.User) ([]models.AboutDevice, error) {
	return AboutDevices[:1], nil
}

func (mdb *MockDB) CreateDeviceForUser(user models.User, aboutDevice models.AboutDevice) error {
	return nil
}

func (mdb *MockDB) UpdateDevice(aboutDevice models.AboutDevice) error {
	return nil
}

func (mdb *MockDB) DeleteDevice(aboutDevice models.AboutDevice) error {
	return nil
}

func (mdb *MockDB) GetData(aboutDevice models.AboutDevice, ptrToData interface{}) error {
	if aboutDevice.ID == 1 {
		var getData = map[reflect.Type]interface{}{
			reflect.TypeOf([]models.Battery{}): BatteryData[:1],
		}
		v := reflect.ValueOf(ptrToData).Elem()
		v.Set(reflect.ValueOf(getData[v.Type()]))
		return nil
	}
	return sql.ErrNoRows
}

func (mdb *MockDB) CreateData(aboutDevice models.AboutDevice, ptrToData interface{}) error {
	return nil
}
