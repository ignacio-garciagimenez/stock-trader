package main

import (
	"net/http"
	"stock-trader/portfolio-service/common"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Validator = common.NewRequestValidator()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello from portfolio-service!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
