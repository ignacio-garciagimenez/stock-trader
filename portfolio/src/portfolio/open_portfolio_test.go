package portfolio

import (
	"stock-trader/portfolio-context/src/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OpenPortfolio(t *testing.T) {
	t.Run("Open portfolio with empty name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Open(OpenPortfolioCommand{
			name: "    ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with short name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Open(OpenPortfolioCommand{
			name: "  Name  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with long name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Open(OpenPortfolioCommand{
			name: "  Name that is longer than 30 characters  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio successfully", func(t *testing.T) {
		repo := &InMemoryPortfolioRepository{
			InMemoryBaseRepository: common.NewInMemoryBaseRepository[PortfolioId, *Portfolio](),
			nameIndex: map[string]*Portfolio{},
		}
		handler := &OpenPortfolioHandler{
			portfolioRepository: repo,
		}

		portfolioId, err := handler.Open(OpenPortfolioCommand{
			name: "  A portfolio name  ",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, portfolioId)
		assert.NotNil(t, func () *Portfolio {
			portfolio, _ := repo.FindById(portfolioId)
			return portfolio
		}())
	})
}
