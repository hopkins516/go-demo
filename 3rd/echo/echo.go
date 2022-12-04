// Copyright © 2022 hops. All rights reserved.

package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello ECHO pkg.")
}

func EchoDemo() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":1234"))
}

func main() {
	EchoDemo()
}
