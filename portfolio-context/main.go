package main

import (
	"net/http"
	"stock-trader/portfolio-context/common"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Validator = common.NewRequestValidator()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World From inside container 12! ")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
