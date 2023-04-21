package portfolio_test

import (
	"stock-trader/portfolio-service/portfolio"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenPortfolio(t *testing.T) {

	t.Run("Open Portfolio Successfully", func(t *testing.T) {
		newPortfolio, err := portfolio.OpenPortfolio("  A Really Looong Portfolio Name  ")

		assert.Nil(t, err)
		assert.NotNil(t, newPortfolio)
		assert.NotNil(t, newPortfolio.Id())
		assert.Equal(t, newPortfolio.Name(), "A Really Looong Portfolio Name")
		assert.IsType(t, portfolio.PortfolioOpened{}, newPortfolio.DomainEvents()[0])
		assert.Equal(t, string(newPortfolio.Id()), newPortfolio.DomainEvents()[0].(portfolio.PortfolioOpened).PortfolioId())
	})

	t.Run("Open Portfolio with empty name", func(t *testing.T) {
		portfolio, err := portfolio.OpenPortfolio("  ")

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

	t.Run("Open Portfolio with short name", func(t *testing.T) {
		portfolio, err := portfolio.OpenPortfolio("  Port  ")

		assert.Error(t, err, "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

	t.Run("Open Portfolio with name too long", func(t *testing.T) {
		portfolio, err := portfolio.OpenPortfolio("  This name is 31 characters long  ")

		assert.Error(t, err, "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

}
