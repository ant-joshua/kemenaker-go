package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=0,lte=80"`
}

func nonEchoHandler(e *echo.Echo) {
	var nativeHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}

	e.GET("/native", echo.WrapHandler(http.HandlerFunc(nativeHandler)))

	var httpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	e.GET("/http", echo.WrapHandler(httpHandler))

	var echoHandler = echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))

	}))

	e.GET("/echo", echoHandler)
}

func basicFunction(e *echo.Echo) {
	// Method string,
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World")
	})

	// Method HTML
	e.GET("/html", func(c echo.Context) error {
		return c.HTML(200, "<h1>Hello World</h1>")
	})

	// Method Redirect
	e.GET("/redirect", func(c echo.Context) error {
		return c.Redirect(302, "/")
	})

	// Method JSON
	e.GET("/json", func(c echo.Context) error {
		data := map[string]any{
			"message": "Hello World",
			"counter": 2,
		}

		return c.JSON(200, data)
	})

	// Parsing Request
	e.GET("/page1", func(c echo.Context) error {
		name := c.QueryParam("name")

		if name == "" {
			return c.String(400, "Name is required")
		}

		data := fmt.Sprintf("Hello %s", name)

		return c.String(200, data)
	})

	e.GET("/page2/:name", func(c echo.Context) error {
		name := c.Param("name")

		data := fmt.Sprintf("Hello %s", name)

		return c.String(200, data)
	})

	e.GET("/page3/:name/*", func(c echo.Context) error {
		name := c.Param("name")
		other := c.Param("*")

		data := fmt.Sprintf("Hello %s %s", name, other)

		return c.String(200, data)
	})

	e.POST("/page4", func(c echo.Context) error {
		name := c.FormValue("name")
		message := c.FormValue("message")

		data := fmt.Sprintf("Hello %s %s", name, message)

		return c.String(200, data)
	})
}

func validatorRoute(e *echo.Echo) {
	e.POST("/users", func(c echo.Context) error {
		//var user = new(User)
		var user = &User{}

		if err := c.Bind(user); err != nil {
			return err
		}

		if err := c.Validate(user); err != nil {
			return err
		}

		return c.JSON(200, user)
	})
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	basicFunction(e)
	nonEchoHandler(e)
	validatorRoute(e)

	err := e.Start(":4000")
	e.Logger.Fatal(err)
}
