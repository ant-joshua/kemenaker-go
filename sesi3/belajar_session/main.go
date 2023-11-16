package main

import (
	"fmt"
	"net/http"

	"github.com/antonlindstrom/pgstore"
	"github.com/globalsign/mgo"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "gopkg.in/mgo.v2"
)

type User struct {
	username string
	password string
}

func cookieStore(e *echo.Echo) {
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/get", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		user := new(User)
		user.username = "admin"
		user.password = "admin"

		sess.Values["foo"] = "bar"
		sess.Values["accessToken"] = user.username
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	})

	e.GET("/get2", func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		if len(sess.Values) == 0 {
			return c.String(http.StatusOK, "No session")
		}

		fmt.Printf("%+v\n", sess.Values["user"])

		return c.String(http.StatusOK, sess.Values["foo"].(string))
	})

	e.GET("/delete", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusTemporaryRedirect, "/get")
	})
}

func newMongoStore() *mongostore.MongoStore {
	mgoSession, err := mgo.Dial("mongodb+srv://hacktiv8:f8MJobSgK22juw75@cluster0.rjmmzo6.mongodb.net/belajar_session")
	if err != nil {
		panic(err)
	}

	defer mgoSession.Close()

	dbCollection := mgoSession.DB("belajar_session").C("sessions")
	maxAge := 86400 * 7
	ensureTTL := true
	authKey := []byte("secret")
	encryptionKey := []byte("secret")

	store := mongostore.NewMongoStore(dbCollection, maxAge, ensureTTL, authKey, encryptionKey)

	return store

}

func newPostgreStore() *pgstore.PGStore {
	url := "postgres://postgres:postgres@localhost:5432/belajar_session?sslmode=disable"
	authKey := []byte("secret")
	encryptionKey := []byte("secret")

	store, err := pgstore.NewPGStore(url, authKey, encryptionKey)

	if err != nil {
		panic(err)
	}

	return store
}

func main() {
	e := echo.New()

	store := newPostgreStore()

	e.Use(echo.WrapMiddleware(context.ClearHandler))
	// e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "https://google.com", "https://hacktiv8.com", "https://kode.id", "https://www.kode.id"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	e.GET("/set-mongo", func(c echo.Context) error {
		sess, _ := store.Get(c.Request(), "postgresSession")

		user := new(User)
		user.username = "admin"
		user.password = "admin"

		sess.Values["foo"] = "bar"
		sess.Values["username"] = user.username
		sess.Save(c.Request(), c.Response())

		return c.NoContent(http.StatusOK)
	})

	e.GET("/get-mongo", func(c echo.Context) error {
		sess, _ := store.Get(c.Request(), "postgresSession")

		if len(sess.Values) == 0 {
			return c.String(http.StatusOK, "No session")
		}

		fmt.Printf("%+v\n", sess.Values["user"])

		return c.String(http.StatusOK, sess.Values["foo"].(string))
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "halo",
		})
	})

	e.Logger.Fatal(e.Start(":8001"))

}
