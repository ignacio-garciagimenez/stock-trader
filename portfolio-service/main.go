package main

import (
	"net/http"
	"stock-trader/portfolio-service/infrastructure"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Validator = infrastructure.NewRequestValidator()
	db, err := infrastructure.ConnectDB()
	if err != nil {
		panic("Could not connect to the database")
	}

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello from portfolio-service!")
	})
	e.POST("/portfolios", BuildOpenPortfolioFeature(db))

	e.Logger.Fatal(e.Start(":8080"))
}
