package main

import (
	"fmt"
	"os"

	"git.jbennett.dev/persona-www/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    defaultH, err := handlers.New()
    if err != nil {
        panic(err)
    }
    
    port, exists := os.LookupEnv("PORT")
    if !exists {
        port = "80"
    }

    e.Static("/static", "static")

    e.GET("/", defaultH.Get)

    e.Start(fmt.Sprintf(":%s", port))
}

