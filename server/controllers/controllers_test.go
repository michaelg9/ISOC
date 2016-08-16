package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/michaelg9/ISOC/server/mocks"
	"github.com/michaelg9/ISOC/server/models"
	"github.com/stretchr/testify/assert"
)

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

	env := newEnv()
	for _, test := range tests {
		testController(env.SignUp, "POST", test.url, test.expected, t)
	}
}

func TestTokenLogin(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{AccessToken: mocks.AccessToken, RefreshToken: mocks.RefreshToken})
	var tests = []struct {
		email    string
		password string
		expected string
	}{
		{mocks.Users[0].Email, "123456", string(jsonResponse)},
		{"user@mail.com", "123456", errWrongPasswordEmail + "\n"},
		{mocks.Users[0].Email, "1234", errWrongPasswordEmail + "\n"},
		{mocks.Users[0].Email, "", errNoPasswordOrEmail + "\n"},
		{"", "123456", errNoPasswordOrEmail + "\n"},
		{"", "", errNoPasswordOrEmail + "\n"},
	}

	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/login?email=%v&password=%v", test.email, test.password)
		testController(env.TokenLogin, "POST", url, test.expected, t)
	}
}

func TestToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{AccessToken: mocks.AccessToken})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{mocks.RefreshToken, string(jsonResponse)},
		//{"", errNoToken + "\n"}, left out because we have no mock session
		{"12345", errTokenInvalid + "\n"},
	}

	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/token?refreshToken=%v", test.refreshToken)
		testController(env.Token, "POST", url, test.expected, t)
	}
}

func TestRefreshToken(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.Tokens{RefreshToken: mocks.RefreshToken})
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{mocks.RefreshToken, string(jsonResponse)},
		{"", errNoToken + "\n"},
		{"12345", errTokenInvalid + "\n"},
	}

	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/refresh?refreshToken=%v", test.refreshToken)
		testController(env.RefreshToken, "POST", url, test.expected, t)
	}
}

func TestTokenLogout(t *testing.T) {
	var tests = []struct {
		refreshToken string
		expected     string
	}{
		{"123", "Success"},
		{"", errNoToken + "\n"},
		{"12345", errTokenAlreadyInvalid + "\n"},
	}

	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/auth/0.1/logout?token=%v", test.refreshToken)
		testController(env.TokenLogout, "POST", url, test.expected, t)
	}
}

func TestGetUser(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.UserResponse{
		User:    mocks.Users[0],
		Devices: mocks.Devices,
	})
	var tests = []struct {
		user     models.User // User that is logged in
		id       int
		expected string
	}{
		{mocks.Users[0], mocks.Users[0].ID, string(jsonResponse)},
		{models.User{ID: 42}, mocks.Users[0].ID, errForbidden + "\n"},
		{models.User{ID: 42}, 42, errWrongUser + "\n"},
	}

	pattern := "/data/{user}"
	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v", test.id)
		testControllerWithPattern(env.GetUser, "GET", url, pattern, test.expected, test.user, t)
	}
}

func TestGetDevice(t *testing.T) {
	jsonResponse, _ := json.Marshal(mocks.Devices[0])
	var tests = []struct {
		user     models.User
		id       int
		deviceID interface{}
		expected string
	}{
		{mocks.Users[0], mocks.Users[0].ID, mocks.AboutDevices[0].ID, string(jsonResponse)},
		{models.User{ID: 42}, mocks.Users[0].ID, mocks.AboutDevices[0].ID, errForbidden + "\n"},
		{models.User{ID: 42}, 42, mocks.AboutDevices[0].ID, errWrongDeviceOrUser + "\n"},
		{mocks.Users[0], mocks.Users[0].ID, "hello", errDeviceIDNotInt + "\n"},
		{mocks.Users[0], mocks.Users[0].ID, 25, errWrongDeviceOrUser + "\n"},
	}

	pattern := "/data/{user}/{device}"
	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v/%v", test.id, test.deviceID)
		testControllerWithPattern(env.GetDevice, "GET", url, pattern, test.expected, test.user, t)
	}
}

func TestGetFeature(t *testing.T) {
	jsonResponse, _ := json.Marshal(models.TrackedData{Battery: mocks.BatteryData[:1]})
	var tests = []struct {
		user     models.User
		id       int
		deviceID interface{}
		feature  string
		expected string
	}{
		{mocks.Users[0], mocks.Users[0].ID, mocks.AboutDevices[0].ID, "Battery", string(jsonResponse)},
		{models.User{ID: 42}, mocks.Users[0].ID, mocks.AboutDevices[0].ID, "Battery", errForbidden + "\n"},
		{models.User{ID: 42}, 42, mocks.AboutDevices[0].ID, "Battery", errWrongDeviceOrUser + "\n"},
		{mocks.Users[0], mocks.Users[0].ID, "hello", "Battery", errDeviceIDNotInt + "\n"},
		{mocks.Users[0], mocks.Users[0].ID, 25, "Battery", errWrongDeviceOrUser + "\n"},
	}

	pattern := "/data/{user}/{device}/{feature}"
	env := newEnv()
	for _, test := range tests {
		url := fmt.Sprintf("/data/%v/%v/%v", test.id, test.deviceID, test.feature)
		testControllerWithPattern(env.GetFeature, "GET", url, pattern, test.expected, test.user, t)
	}
}

/* Test helper functions */

func TestWriteResponse(t *testing.T) {
	response := models.DataOut{Device: mocks.Devices}
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
		assert.Equal(t, test.expected, obtained)
	}
}

/* Helper functions */

func testController(controller http.HandlerFunc, method, url, expected string, t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)

	http.HandlerFunc(controller).ServeHTTP(rec, req)

	obtained := rec.Body.String()
	assert.Equal(t, expected, obtained)
}

// TODO: Parameters as struct
func testControllerWithPattern(controller http.HandlerFunc, method, url, pattern, expected string, user models.User, t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	context.Set(req, UserKey, user)

	r := mux.NewRouter()
	r.HandleFunc(pattern, controller)
	r.ServeHTTP(rec, req)

	obtained := rec.Body.String()
	assert.Equal(t, expected, obtained)
}

func newEnv() Env {
	return Env{DB: &mocks.MockDB{}, Tokens: &mocks.MockTokens{}}
}
