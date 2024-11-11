package main

import (
	"git.jbennett.dev/persona-www/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    defaultH, err := handlers.New()
    if err != nil {
        panic(err)
    }

    e.Static("/static", "static")

    e.GET("/", defaultH.Get)

    e.Start("127.0.0.1:3000")

}

