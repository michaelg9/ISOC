package authentication

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/services/models"
	"github.com/urfave/negroni"
)

var users = []models.User{
	models.User{
		ID:           1,
		Email:        "user@usermail.com",
		PasswordHash: "$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK",
		APIKey:       "37e72ff927f511e688adb827ebf7e157",
	},
}

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

func (mdb *mockDB) GetDevicesFromUser(user models.User) ([]models.Device, error) {
	return []models.Device{}, nil
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
	return nil
}

func (mdb *mockDB) CreateData(device models.DeviceStored, ptrToData interface{}) error {
	return nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func TestRequireBasicAuth(t *testing.T) {
	var tests = []struct {
		basicAuth bool
		email     string
		password  string
		expected  string
	}{
		{true, "user@usermail.com", "123456", "Hello world!"},
		{true, "user@usermail.com", "1234", "Unauthorized request.\n"},
		{true, "", "", "Unauthorized request.\n"},
		{false, "", "", "Unauthorized request.\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hello", nil)
		if test.basicAuth {
			req.SetBasicAuth(test.email, test.password)
		}

		env := MiddlewareEnv{&controllers.Env{DB: &mockDB{}}}
		var handler http.Handler
		handler = negroni.New(
			negroni.HandlerFunc(env.RequireBasicAuth),
			negroni.Wrap(http.HandlerFunc(helloHandler)),
		)
		handler.ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func TestRequireTokenAuth(t *testing.T) {
	tomorrow := time.Now().Add(time.Hour * time.Duration(24)).Unix()
	yesterday := time.Now().AddDate(0, 0, -1).Unix()
	var tests = []struct {
		tokenAuth bool
		token     string
		expected  string
	}{
		{true, newToken("user@usermail.com", tomorrow), "Hello world!"},
		{true, newToken("user@usermail.com", yesterday), "Unauthorized request.\n"},
		{true, newToken("user@mail.com", tomorrow), "Unauthorized request.\n"},
		{true, "", "Unauthorized request.\n"},
		{false, "", "Unauthorized request.\n"},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hello", nil)
		if test.tokenAuth {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.token))
		}

		env := MiddlewareEnv{&controllers.Env{DB: &mockDB{}}}
		var handler http.Handler
		handler = negroni.New(
			negroni.HandlerFunc(env.RequireTokenAuth),
			negroni.Wrap(http.HandlerFunc(helloHandler)),
		)
		handler.ServeHTTP(rec, req)

		obtained := rec.Body.String()
		if test.expected != obtained {
			t.Errorf("\n...expected = %v\n...obtained = %v", test.expected, obtained)
		}
	}
}

func newToken(email string, time int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time,
	})

	tokenString, _ := token.SignedString([]byte(hmacSecret))
	return tokenString
}
