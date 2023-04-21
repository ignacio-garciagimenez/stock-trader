package main

import (
	"net/http"
	"stock-trader/portfolio-service/common"
	"stock-trader/portfolio-service/portfolio"
	portfolio_features "stock-trader/portfolio-service/portfolio/features"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	e.Validator = common.NewRequestValidator()
	db, err := common.ConnectDB()
	if err != nil {
		panic("Could not connect to the database")
	}

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello from portfolio-service!")
	})
	e.POST("/portfolios", common.WithTransaction(db, func(tx *gorm.DB) echo.HandlerFunc {
		return portfolio_features.NewPorfolioEndpoint(
			portfolio_features.NewOpenPortfolioHandler(
				portfolio.NewPortfolioRepository(tx),
			),
		).Open
	}))

	e.Logger.Fatal(e.Start(":8080"))
}
