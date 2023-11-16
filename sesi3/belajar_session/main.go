package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type User struct {
	username string
	password string
}

func main() {
	e := echo.New()

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

	// e.GET("/delete", func(c echo.Context) error {
	// 	sess, _ := session.Get("session", c)
	// 	sess.Options.MaxAge = -1
	// 	sess.Save(c.Request(), c.Response())

	// 	return c.Redirect(http.StatusTemporaryRedirect, "/get")
	// })

	e.Logger.Fatal(e.Start(":8001"))

}
