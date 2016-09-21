package mocks

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/michaelg9/ISOC/server/models"
)

// MockDB implements the Datastore interface so we can use it to mock the DB
type MockDB struct{}

func (mdb *MockDB) GetUser(user models.User) (models.User, error) {
	sameEmail := user.Email == Users[0].Email
	sameID := user.ID == Users[0].ID
	if sameEmail || sameID {
		return Users[0], nil
	}
	return models.User{}, sql.ErrNoRows
}

func (mdb *MockDB) CreateUser(user models.User) (int, error) {
	return 2, nil
}

func (mdb *MockDB) UpdateUser(user models.User) error {
	if user.ID == 0 {
		return errors.New("No user id.")
	}

	return nil
}

func (mdb *MockDB) DeleteUser(user models.User) error {
	return nil
}

func (mdb *MockDB) GetAllDevices() ([]models.AboutDevice, error) {
	return AboutDevices[:1], nil
}

func (mdb *MockDB) GetDeviceFromUser(user models.User, device models.Device) (models.Device, error) {
	rightDevice := device.AboutDevice.ID == Devices[0].AboutDevice.ID
	rightUser := user.ID == Users[0].ID
	if rightDevice && rightUser {
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

func (mdb *MockDB) CreateDeviceForUser(user models.User, aboutDevice models.AboutDevice) (int, error) {
	return 2, nil
}

func (mdb *MockDB) UpdateDevice(aboutDevice models.AboutDevice) error {
	return nil
}

func (mdb *MockDB) DeleteDevice(aboutDevice models.AboutDevice) error {
	return nil
}

func (mdb *MockDB) GetAllFeatureData(ptrToFeature interface{}) error {
	var getData = map[reflect.Type]interface{}{
		reflect.TypeOf([]models.Battery{}):    BatteryData[:1],
		reflect.TypeOf([]models.Call{}):       CallData[:1],
		reflect.TypeOf([]models.App{}):        AppData[:1],
		reflect.TypeOf([]models.Runservice{}): RunserviceData[:1],
	}
	v := reflect.ValueOf(ptrToFeature).Elem()
	v.Set(reflect.ValueOf(getData[v.Type()]))
	return nil
}

func (mdb *MockDB) GetFeatureOfDevice(aboutDevice models.AboutDevice, ptrToFeature interface{}) error {
	if aboutDevice.ID == 1 {
		var getData = map[reflect.Type]interface{}{
			reflect.TypeOf([]models.Battery{}):    BatteryData[:1],
			reflect.TypeOf([]models.Call{}):       CallData[:1],
			reflect.TypeOf([]models.App{}):        AppData[:1],
			reflect.TypeOf([]models.Runservice{}): RunserviceData[:1],
		}
		v := reflect.ValueOf(ptrToFeature).Elem()
		v.Set(reflect.ValueOf(getData[v.Type()]))
		return nil
	}
	return sql.ErrNoRows
}

func (mdb *MockDB) CreateFeatureForDevice(aboutDevice models.AboutDevice, ptrToFeature interface{}) error {
	return nil
}
