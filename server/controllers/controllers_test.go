package controllers

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	return models.User{}, nil
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
	ptr, _ := ptrToData.(*interface{})
	*ptr = batteryData
	return nil
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
// TODO: create URL later in testbody

func TestSignUp(t *testing.T) {
	var tests = []struct {
		url      string
		expected string
	}{
		{"/signup?email=user@mail.com&password=123456", "Success"},
		{"/signup?email=user@usermail.com&password=123456", "User already exists"},
		{"/signup?email=bla@blabla", errMissingPasswordOrEmail + "\n"},
		{"/signup?password=123456", errMissingPasswordOrEmail + "\n"},
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
		{"/data/0.1/q", errNoAPIKey + "\n"},
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
		{"user@usermail.com", "", errMissingPasswordOrEmail + "\n"},
		{"", "123456", errMissingPasswordOrEmail + "\n"},
		{"", "", errMissingPasswordOrEmail + "\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		url := fmt.Sprintf("/auth/0.1/login?email=%v&password=%v", test.email, test.password)
		req, _ := http.NewRequest("POST", url, nil)

		env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
		http.HandlerFunc(env.Login).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func TestToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{AccessToken: "123"})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", string(jsonResponse)},
		{"", errMissingToken + "\n"},
		{"12345", "Token is invalid." + "\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		url := fmt.Sprintf("/auth/0.1/token?refreshToken=%v", test.refreshToken)
		req, _ := http.NewRequest("POST", url, nil)

		env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
		http.HandlerFunc(env.Token).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func TestRefreshToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{RefreshToken: "123"})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", string(jsonResponse)},
		{"", errMissingToken + "\n"},
		{"12345", "Token is invalid." + "\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		url := fmt.Sprintf("/auth/0.1/refresh?refreshToken=%v", test.refreshToken)
		req, _ := http.NewRequest("POST", url, nil)

		env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
		http.HandlerFunc(env.RefreshToken).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func TestLogoutToken(t *testing.T) {
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", "Success"},
		{"", errMissingToken + "\n"},
		{"12345", "Token already invalid." + "\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		url := fmt.Sprintf("/auth/0.1/logout?token=%v", test.refreshToken)
		req, _ := http.NewRequest("POST", url, nil)

		env := Env{DB: &mockDB{}, Tokens: &mockTokens{}}
		http.HandlerFunc(env.LogoutToken).ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
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
