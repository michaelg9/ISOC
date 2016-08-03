package controllers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaelg9/ISOC/server/models"
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

/* mockDB implements the Datastore interface so we can use it to mock the DB */

type mockDB struct{}

func (mdb *mockDB) GetUser(user models.User) (models.User, error) {
	if user.Email == users[0].Email {
		return users[0], nil
	}
	return models.User{}, sql.ErrNoRows
}

func (mdb *mockDB) CreateUser(user models.User) error {
	return nil
}

func (mdb *mockDB) UpdateUser(user models.User) error {
	return nil
}

func (mdb *mockDB) DeleteUser(user models.User) error {
	return nil
}

func (mdb *mockDB) GetDevice(device models.Device) (models.Device, error) {
	if device.DeviceInfo.ID == 1 {
		return devices[0], nil
	}
	return models.Device{}, sql.ErrNoRows
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

func (mdb *mockDB) UpdateDevice(device models.DeviceStored) error {
	return nil
}

func (mdb *mockDB) DeleteDevice(device models.DeviceStored) error {
	return nil
}

func (mdb *mockDB) GetData(device models.DeviceStored, ptrToData interface{}) error {
	if device.ID == 1 {
		var getData = map[reflect.Type]interface{}{
			reflect.TypeOf([]models.Battery{}): batteryData[:1],
		}
		v := reflect.ValueOf(ptrToData).Elem()
		v.Set(reflect.ValueOf(getData[v.Type()]))
		return nil
	}
	return sql.ErrNoRows
}

func (mdb *mockDB) CreateData(device models.DeviceStored, ptrToData interface{}) error {
	return nil
}

/* mockTokens implements the TokenControl interface so it can be used to mock the token back-end */

type mockTokens struct{}

func (mTkns *mockTokens) CheckToken(tokenString string) (email string, err error) {
	if tokenString == "123" {
		return "user@usermail.com", nil
	}
	return "", errors.New("Token is invalid.")
}

func (mTkns *mockTokens) NewToken(user models.User, duration time.Duration) (string, error) {
	return "123", nil
}

func (mTkns *mockTokens) InvalidateToken(tokenString string) error {
	if tokenString != "123" {
		return errors.New("Token already invalid.")
	}
	return nil
}

/* Test controller functions */

func TestSignUp(t *testing.T) {
	var tests = []struct {
		url      string
		expected string
	}{
		{"/signup?email=user@mail.com&password=123456", "Success"},
		{"/signup?email=user@usermail.com&password=123456", "User already exists"},
		{"/signup?email=bla@blabla", errNoPasswordOrEmail + "\n"},
		{"/signup?password=123456", errNoPasswordOrEmail + "\n"},
	}

	env := Env{DB: &mockDB{}}
	for _, test := range tests {
		testController(env.SignUp, "POST", test.url, test.expected, t)
	}
}

func TestDownload(t *testing.T) {
	response := models.DataOut{Device: devices}
	jsonResponse, _ := json.Marshal(response)
	xmlResponse, _ := xml.Marshal(response)
	var tests = []struct {
		apiKey   string
		out      string
		expected string
	}{
		{"", "", errNoAPIKey + "\n"},
		{"12345", "", string(jsonResponse)},
		{"12345", "json", string(jsonResponse)},
		{"12345", "xml", string(xmlResponse)},
	}

	env := Env{DB: &mockDB{}}
	for _, test := range tests {
		var url string
		if test.out != "" {
			url = fmt.Sprintf("/data/0.1/q?appid=%v&out=%v", test.apiKey, test.out)
		} else {
			url = fmt.Sprintf("/data/0.1/q?appid=%v", test.apiKey)
		}
		testController(env.Download, "GET", url, test.expected, t)
	}
}

func TestLogin(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{AccessToken: "123", RefreshToken: "123"})
	var tests = []struct {
		email    string
		password string
		expected string
	}{
		{"user@usermail.com", "123456", string(jsonResponse)},
		{"user@mail.com", "123456", errWrongPasswordEmail + "\n"},
		{"user@usermail.com", "1234", errWrongPasswordEmail + "\n"},
		{"user@usermail.com", "", errNoPasswordOrEmail + "\n"},
		{"", "123456", errNoPasswordOrEmail + "\n"},
		{"", "", errNoPasswordOrEmail + "\n"},
	}

	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/login?email=%v&password=%v", test.email, test.password)
		testController(env.Login, "POST", url, test.expected, t)
	}
}

func TestToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{AccessToken: "123"})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", string(jsonResponse)},
		{"", errNoToken + "\n"},
		{"12345", errTokenInvalid + "\n"},
	}

	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/token?refreshToken=%v", test.refreshToken)
		testController(env.Token, "POST", url, test.expected, t)
	}
}

func TestRefreshToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{RefreshToken: "123"})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", string(jsonResponse)},
		{"", errNoToken + "\n"},
		{"12345", errTokenInvalid + "\n"},
	}

	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/refresh?refreshToken=%v", test.refreshToken)
		testController(env.RefreshToken, "POST", url, test.expected, t)
	}
}

func TestLogoutToken(t *testing.T) {
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", "Success"},
		{"", errNoToken + "\n"},
		{"12345", errTokenAlreadyInvalid + "\n"},
	}

	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/logout?token=%v", test.refreshToken)
		testController(env.LogoutToken, "POST", url, test.expected, t)
	}
}

func TestUser(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.UserResponse{
		User:    users[0],
		Devices: devices,
	})
	var tests = []struct {
		email    string
		apiKey   string
		expected string
	}{
		{users[0].Email, users[0].APIKey, string(jsonResponse)},
		{"user@mail.com", users[0].APIKey, errWrongUser + "\n"},
		{users[0].Email, "123", errWrongAPIKey + "\n"},
	}

	pattern := "/data/{email}"
	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v?appid=%v", test.email, test.apiKey)
		testControllerWithPattern(env.User, "GET", url, pattern, test.expected, t)
	}
}

func TestDevice(t *testing.T) {
	jsonResponse, _ := json.Marshal(devices[0])
	var tests = []struct {
		email    string
		apiKey   string
		deviceID interface{}
		expected string
	}{
		{users[0].Email, users[0].APIKey, deviceInfos[0].ID, string(jsonResponse)},
		{"user@mail.com", users[0].APIKey, deviceInfos[0].ID, errWrongUser + "\n"},
		{users[0].Email, "123", deviceInfos[0].ID, errWrongAPIKey + "\n"},
		{users[0].Email, users[0].APIKey, "hello", errDeviceIDNotInt + "\n"},
		{users[0].Email, users[0].APIKey, "25", errWrongDevice + "\n"},
	}

	pattern := "/data/{email}/{device}"
	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v/%v?appid=%v", test.email, test.deviceID, test.apiKey)
		testControllerWithPattern(env.Device, "GET", url, pattern, test.expected, t)
	}
}

func TestFeature(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.DeviceData{Battery: batteryData[:1]})
	var tests = []struct {
		email    string
		apiKey   string
		deviceID interface{}
		feature  string
		expected string
	}{
		{users[0].Email, users[0].APIKey, deviceInfos[0].ID, "Battery", string(jsonResponse)},
		{users[0].Email, "1234", deviceInfos[0].ID, "Battery", errWrongAPIKey + "\n"},
		{"user@mail.com", users[0].APIKey, deviceInfos[0].ID, "Battery", errWrongUser + "\n"},
		{users[0].Email, users[0].APIKey, "hello", "Battery", errDeviceIDNotInt + "\n"},
		{users[0].Email, users[0].APIKey, "25", "Battery", errWrongDevice + "\n"},
	}

	pattern := "/data/{email}/{device}/{feature}"
	env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v/%v/%v?appid=%v", test.email, test.deviceID, test.feature, test.apiKey)
		testControllerWithPattern(env.Feature, "GET", url, pattern, test.expected, t)
	}
}

/* Test helper functions */

func TestWriteResponse(t *testing.T) {
	response := models.DataOut{Device: devices}
	jsonResponse, _ := json.Marshal(response)
	xmlResponse, _ := xml.Marshal(response)
	var tests = []struct {
		format   string
		response interface{}
		expected string
	}{
		{"", response, string(jsonResponse)},
		{"json", response, string(jsonResponse)},
		{"xml", response, string(xmlResponse)},
		{"yml", response, ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		_ = writeResponse(rec, test.format, test.response)
		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

/* Helper functions */

func testController(controller http.HandlerFunc, method, url, expected string, t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)

	http.HandlerFunc(controller).ServeHTTP(rec, req)

	obtained := rec.Body.String()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}

func testControllerWithPattern(controller http.HandlerFunc, method, url, pattern, expected string, t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)

	r := mux.NewRouter()
	r.HandleFunc(pattern, controller)
	r.ServeHTTP(rec, req)

	obtained := rec.Body.String()
	if expected != obtained {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, obtained)
	}
}
