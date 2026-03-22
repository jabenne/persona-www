package main

import (
	"fmt"
	"os"

	"git.jbennett.dev/persona-www/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	lanyardID := os.Getenv("LANYARD_ID")

	defaultH, err := handlers.New(lanyardID)
	if err != nil {
		panic(err)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "3030"
	}

	e.Static("/static", "static")

	e.GET("/", defaultH.Get)
	e.GET("/presence", defaultH.GetPresence)

	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
