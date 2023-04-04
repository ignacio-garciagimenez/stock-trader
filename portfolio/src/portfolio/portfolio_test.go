package portfolio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenPortfolio(t *testing.T) {

	t.Run("Open Portfolio Successfully", func(t *testing.T) {
		portfolio, err := OpenPortfolio("  A Really Looong Portfolio Name  ")

		assert.Nil(t, err)
		assert.NotNil(t, portfolio)
		assert.NotNil(t, portfolio.Id())
		assert.Equal(t, portfolio.Name(), "A Really Looong Portfolio Name")
		assert.IsType(t, PortfolioOpened{}, portfolio.domainEvents[0])
		assert.Equal(t, portfolio.Id(), portfolio.domainEvents[0].(PortfolioOpened).PortfolioId())
	})

	t.Run("Open Portfolio with empty", func(t *testing.T) {
		portfolio, err := OpenPortfolio("  ")

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

	t.Run("Open Portfolio with short name", func(t *testing.T) {
		portfolio, err := OpenPortfolio("  Port  ")

		assert.Error(t, err, "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

	t.Run("Open Portfolio with name too long", func(t *testing.T) {
		portfolio, err := OpenPortfolio("  This name is 31 characters long  ")

		assert.Error(t, err, "portfolio name must be between 5 and 30 characters long")
		assert.Nil(t, portfolio)
	})

}
