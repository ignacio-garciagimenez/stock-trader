package main

import (
	"stock-trader/portfolio-service/infrastructure"
	"stock-trader/portfolio-service/portfolio"
	portfolio_features "stock-trader/portfolio-service/portfolio/features"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BuildOpenPortfolioFeature(db *gorm.DB) echo.HandlerFunc {
	return infrastructure.WithTransaction(db, func(tx *gorm.DB) echo.HandlerFunc {
		return portfolio_features.NewOpenPortfolioEndpoint(
			portfolio_features.NewOpenPortfolioHandler(
				portfolio.NewPortfolioRepository(tx),
			),
		).Open
	})
}
