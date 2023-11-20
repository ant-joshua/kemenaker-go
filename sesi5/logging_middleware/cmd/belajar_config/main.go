package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	e := echo.New()

	// viper.SetConfigType("json")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		e.Logger.Fatal(err)
		panic(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World")
	})

	viper.WatchConfig()
	viper.OnConfigChange(func(event fsnotify.Event) {
		fmt.Println("Config file changed:", event.Name)
		e.Logger.Info("Config file changed:", event.Name)
	})

	// os.Getenv()

	e.Logger.Fatal(e.Start(":" + viper.GetString("server.port")))
}
