package main

import (
	"fmt"
	"os"

	"git.jbennett.dev/persona-www/handlers"
	"git.jbennett.dev/persona-www/services/lanyard"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	lanyardOpts := make([]lanyard.Option, 0, 2)
	lanyardURL, exists := os.LookupEnv("LANYARD_URL")
	if exists {
		lanyardOpts = append(lanyardOpts, lanyard.WithBaseURL(lanyardURL))
	}

	lanyardID, exists := os.LookupEnv("LANYARD_ID")
	if exists {
		lanyardOpts = append(lanyardOpts, lanyard.WithDefaultID(lanyardID))
	}

	lanyardSvc := lanyard.New(lanyardOpts...)

	defaultH, err := handlers.New(lanyardSvc)
	if err != nil {
		panic(err)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	e.Static("/static", "static")

	e.GET("/", defaultH.Get)
	e.GET("/presence", defaultH.GetPresence)

	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}
}
