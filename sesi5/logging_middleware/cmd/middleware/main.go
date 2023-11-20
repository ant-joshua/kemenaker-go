package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func middlewareOne(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("middlewareOne")
		return next(c)
	}
}

func middlewareTwo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("middlewareTwo")
		return next(c)
	}
}

func middlewareOther(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("middlewareOther")
		next.ServeHTTP(w, r)
	})
}

func makeLogEntry(c echo.Context) *logrus.Entry {
	if c == nil {
		return logrus.WithFields(logrus.Fields{
			"at": time.Now().Format("2006-01-02 15:00:00"),
		})
	}

	return logrus.WithFields(logrus.Fields{
		"at":         time.Now().Format("2006-01-02 15:00:00"),
		"method":     c.Request().Method,
		"uri":        c.Request().RequestURI,
		"remote_ip":  c.Request().RemoteAddr,
		"user_agent": c.Request().UserAgent(),
	})

}

func middlewareLogrus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("middlewareLogrus")
		makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

func belajarLog(e *echo.Echo) {
	e.Use(middlewareOne)
	e.Use(middlewareTwo)
	e.Use(echo.WrapMiddleware(middlewareOther))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middlewareLogrus)

	e.GET("/", func(c echo.Context) error {
		fmt.Println("hello")
		return c.String(200, "Hello, World!")
	})
}

var (
	app     = kingpin.New("app", "A command-line application skeleton.")
	verbose = app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	// name    = kingpin.Arg("name", "Name of user.").Required().String()
	argAppName = app.Arg("name", "Application Name").Required().String()
	argPort    = app.Arg("port", "Application Port").Default("5000").Int()
)

func main() {

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		// logrus.Fatal(err)
		panic(err)
	}

	fmt.Printf("command: %s\n", command)

	fmt.Printf("%v, %s\n", *verbose, *argAppName)

	e := echo.New()

	lock := make(chan error)
	go func(lock chan error) {
		lock <- e.Start(fmt.Sprintf(":%d", *argPort))
	}(lock)

	time.Sleep(1 * time.Millisecond)
	makeLogEntry(nil).Warning("application started without ssl/tls enabled")

	err = <-lock

	if err != nil {
		makeLogEntry(nil).Panic("failed to start application")
	}

	// e.Logger.Fatal(e.Start(":4001"))
}
