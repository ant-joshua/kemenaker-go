package main

import (
	"fmt"
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type M map[string]interface{}

func main() {
	tmpl := template.Must(template.ParseGlob("./*.html"))

	const csrfTokenHeader = "X-CSRF-Token"
	const CSRFKey = "csrf"

	e := echo.New()

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:" + csrfTokenHeader,
		ContextKey:  CSRFKey,
	}))

	e.GET("/", func(c echo.Context) error {
		data := make(M)
		data[CSRFKey] = c.Get(CSRFKey)

		return tmpl.Execute(c.Response(), data)
	})

	e.POST("/sayhello", func(c echo.Context) error {
		data := make(M)

		if err := c.Bind(&data); err != nil {
			return err
		}

		message := fmt.Sprintf("Hello %s", data["name"])

		return c.String(200, message)
	})

	e.Logger.Fatal(e.Start(":9000"))
}
