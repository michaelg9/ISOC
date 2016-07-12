package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/michaelg9/ISOC/server/controllers"
	"github.com/michaelg9/ISOC/server/routers"
	"github.com/michaelg9/ISOC/server/services/models"

	"github.com/urfave/negroni"
)

func main() {
	db := startDB()
	sessionStore := startSessions()
	env := controllers.Env{DB: db, SessionStore: sessionStore}

	router := routers.NewRouter(env)

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(context.ClearHandler(router))

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, n)
}

func startDB() *models.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PWD")
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/mobile_data?parseTime=true", user, pwd, host)
	db, err := models.NewDB(dsn)
	if err != nil {
		panic(err)
	}

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
