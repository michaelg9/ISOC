package authentication

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/mocks"
	"github.com/urfave/negroni"
)

/* Dummy handler */

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
		{true, mocks.Users[0].Email, "123456", "Hello world!"},
		{true, mocks.Users[0].Email, "1234", "Unauthorized request.\n"},
		{true, "", "", "Unauthorized request.\n"},
		{false, "", "", "Unauthorized request.\n"},
	}

	env := newEnv()
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hello", nil)
		if test.basicAuth {
			req.SetBasicAuth(test.email, test.password)
		}

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
	var tests = []struct {
		tokenAuth bool
		token     string
		expected  string
	}{
		{true, mocks.JWT, "Hello world!"},
		{true, "12345", "Unauthorized request.\n"},
		{true, "", "Unauthorized request.\n"},
		{false, "", "Unauthorized request.\n"},
	}

	env := newEnv()
	for _, test := range tests {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/hello", nil)
		if test.tokenAuth {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", test.token))
		}

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

func newEnv() MiddlewareEnv {
	return MiddlewareEnv{&controllers.Env{DB: &mocks.MockDB{}, Tokens: &mocks.MockTokens{}}}
}
