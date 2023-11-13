package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func basicRoute(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/json", func(c echo.Context) error {

		data := map[string]any{
			"message": "Hello, World!",
		}

		return c.JSON(200, data)
	})

	e.GET("/html", func(c echo.Context) error {
		return c.HTML(200, "<h1>Hello, World!</h1>")
	})

	e.GET("/redirect", func(c echo.Context) error {
		return c.Redirect(302, "/")
	})

	e.GET("page1", func(c echo.Context) error {
		name := c.QueryParam("name")
		return c.String(200, "Hello, "+name)
	})

	e.GET("product/:id", func(c echo.Context) error {
		id := c.Param("id")
		// parseID, err := strconv.Atoi(id)
		// if err != nil {
		// 	return c.String(400, "Invalid ID")
		// }

		return c.String(200, "Product ID: "+id)
	})

	e.POST("/page4", func(c echo.Context) error {
		name := c.FormValue("name")
		message := c.FormValue("message")

		return c.String(200, "Name: "+name+"\nMessage: "+message)
	})

	e.Static("/static", "web")
}

type User struct {
	Name  string `json:"name" form:"name" query:"name" validate:"required"`
	Email string `json:"email" form:"email" query:"email" validate:"required,email"`
	Age   int    `json:"age" form:"age" query:"age" validate:"gte=0,lte=100"`
}

type CreateProductRequest struct {
	Name        string   `json:"name" form:"name" query:"name"`
	Description string   `json:"description" form:"description" query:"description"`
	Price       int      `json:"price" form:"price" query:"price"`
	Attributes  []string `json:"attributes" form:"attributes" query:"attributes"`
}

type UpdateProductRequest struct {
	Name        string   `json:"name" form:"name" query:"name"`
	Description string   `json:"description" form:"description" query:"description"`
	Price       int      `json:"price" form:"price" query:"price"`
	Attributes  []string `json:"attributes" form:"attributes" query:"attributes"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	fmt.Println("Hello, World!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)

		if !ok {
			report = echo.NewHTTPError(500, err.Error())
		}

		report.Code = 500

		// castedObject, ok := err.(validator.ValidationErrors)

		// if ok {

		// 	report.Code = 400

		// 	for _, err := range castedObject {
		// 		switch err.Tag() {
		// 		case "required":
		// 			report.Message = fmt.Sprintf("%s is required", err.Field())
		// 		case "email":
		// 			report.Message = fmt.Sprintf("%s is not valid email", err.Field())
		// 		}

		// 	}

		// }

		errPage := fmt.Sprintf("%d.html", report.Code)

		if err := c.File("web/" + errPage); err != nil {
			c.HTML(report.Code, "Error")
		}

		c.JSON(report.Code, report)
	}

	basicRoute(e)

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Any("/any", func(c echo.Context) error {
		u := new(User)

		err := c.Bind(u)

		if err != nil {
			return err
		}

		err = c.Validate(u)

		if err != nil {
			return err
		}

		return c.JSON(200, u)
	})

	e.Logger.Fatal(e.Start(":4000"))
}
