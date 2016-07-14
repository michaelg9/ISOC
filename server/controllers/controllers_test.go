package controllers

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michaelg9/ISOC/server/services/models"
)

var users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	},
	models.User{
		ID:     2,
		Email:  "user@mail.com",
		APIKey: "",
	},
}

var deviceInfos = []models.DeviceStored{
	models.DeviceStored{
		ID:           1,
		Manufacturer: "Motorola",
		Model:        "Moto X (2nd Generation)",
		OS:           "Android 5.0",
	},
	models.DeviceStored{
		ID:           2,
		Manufacturer: "One Plus",
		Model:        "Three",
		OS:           "Android 6.0",
	},
}

var devices = []models.Device{
	models.Device{
		DeviceInfo: deviceInfos[0],
		Data: models.DeviceData{
			Battery: batteryData[:1],
		},
	},
}

var batteryData = []models.Battery{
	models.Battery{
		Value: 70,
		Time:  "2016-05-31 11:48:48",
	},
	models.Battery{
		Value: 71,
		Time:  "2016-05-31 11:50:31",
	},
}

/* mockDB has to implement:
 * type Datastore interface {
 * 	GetUser(user User) (fullUser User, err error)
 * 	CreateUser(user User) error
 * 	UpdateUser(user User, field string) error
 * 	DeleteUser(user User) error
 * 	GetDevicesFromUser(user User) ([]Device, error)
 * 	GetDeviceInfos(user User) ([]DeviceStored, error)
 * 	CreateDeviceForUser(user User, device DeviceStored) error
 * 	UpdateDevice(device DeviceStored, field string) error
 * 	DeleteDevice(device DeviceStored) error
 * 	GetData(device DeviceStored, ptrToData interface{}) error
 * 	CreateData(device DeviceStored, ptrToData interface{}) error
 * }
 */

type mockDB struct{}

func (mdb *mockDB) GetUser(user models.User) (models.User, error) {
	if user.Email == users[0].Email {
		return users[0], nil
	}
	return models.User{}, nil
}

func (mdb *mockDB) CreateUser(user models.User) error {
	return nil
}

func (mdb *mockDB) UpdateUser(user models.User, field string) error {
	return nil
}

func (mdb *mockDB) DeleteUser(user models.User) error {
	return nil
}

func (mdb *mockDB) GetDevicesFromUser(user models.User) ([]models.Device, error) {
	return devices[:1], nil
}

func (mdb *mockDB) GetDeviceInfos(user models.User) ([]models.DeviceStored, error) {
	return deviceInfos[:1], nil
}

func (mdb *mockDB) CreateDeviceForUser(user models.User, device models.DeviceStored) error {
	return nil
}

func (mdb *mockDB) UpdateDevice(device models.DeviceStored, field string) error {
	return nil
}

func (mdb *mockDB) DeleteDevice(device models.DeviceStored) error {
	return nil
}

func (mdb *mockDB) GetData(device models.DeviceStored, ptrToData interface{}) error {
	ptr, _ := ptrToData.(*interface{})
	*ptr = batteryData
	return nil
}

func (mdb *mockDB) CreateData(device models.DeviceStored, ptrToData interface{}) error {
	return nil
}

func TestSignUp(t *testing.T) {
	var tests = []struct {
		url      string
		expected string
	}{
		{"/signup?email=user@mail.com&password=123456", "Success"},
		{"/signup?email=user@usermail.com&password=123456", "User already exists"},
		{"/signup?email=bla@blabla", "You have to specify a password and an email.\n"},
		{"/signup?password=123456", "You have to specify a password and an email.\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", test.url, nil)

		env := Env{DB: &mockDB{}}
		http.HandlerFunc(env.SignUp).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func TestDownload(t *testing.T) {
	response := models.DataOut{Device: devices}
	jsonResponse, _ := json.Marshal(response)
	xmlResponse, _ := xml.Marshal(response)
	var tests = []struct {
		url      string
		expected string
	}{
		{"/data/0.1/q", "No API Key given.\n"},
		{"/data/0.1/q?appid=12345", string(jsonResponse)},
		{"/data/0.1/q?appid=12345&out=json", string(jsonResponse)},
		{"/data/0.1/q?appid=12345&out=xml", string(xmlResponse)},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", test.url, nil)

		env := Env{DB: &mockDB{}}
		http.HandlerFunc(env.Download).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}
