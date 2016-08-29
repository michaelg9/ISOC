package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/models"
	"github.com/michaelg9/ISOC/server/routers"

	"github.com/urfave/negroni"
)

// Number of redis database we use to store the tokens
const dbNR = 0

func main() {
	time.sleep(20 * time.Second)
	db := startDB()
	sessionStore := startSessions()
	tokenstore := models.NewTokenstore(os.Getenv("REDIS_HOST")+":6379", dbNR)
	env := controllers.Env{DB: db, SessionStore: sessionStore, Tokens: tokenstore}

	router := routers.NewRouter(env)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(context.ClearHandler(router))

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, n)
}

func startDB() *models.DB {
	// Get database access variables
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PWD")
	// Create DSN
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/mobile_data?parseTime=true", user, pwd, host)
	// Panics if there is an error with the connection
	db := models.NewDB(dsn)
	return db
}

func startSessions() *sessions.CookieStore {
	// Initiate the a cookie store to save our sessions
	sessionStore := sessions.NewCookieStore([]byte("something-very-secret"))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		Secure:   false, // NOTE: Change to true if on actual server
		HttpOnly: true,
	}

	return sessionStore
}
