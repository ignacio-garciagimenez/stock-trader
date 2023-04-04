package portfolio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"stock-trader/portfolio-context/src/common"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_OpenPortfolioHandler(t *testing.T) {
	t.Run("Open portfolio with empty name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "    ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with short name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  Name  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with long name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  Name that is longer than 30 characters  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio successfully", func(t *testing.T) {
		repo := &InMemoryPortfolioRepository{
			InMemoryBaseRepository: common.NewInMemoryBaseRepository[PortfolioId, *Portfolio](),
			nameIndex:              map[string]*Portfolio{},
		}
		handler := &OpenPortfolioHandler{
			portfolioRepository: repo,
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  A portfolio name  ",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, portfolioId)
		assert.NotNil(t, func() *Portfolio {
			portfolio, _ := repo.FindById(portfolioId)
			return portfolio
		}())
	})
}

func Test_OpenPortfolioEndpoint(t *testing.T) {
	t.Run("Open Portfolio Successfully", func(t *testing.T) {
		portfolioId := PortfolioId(uuid.NewString())
		endpoint := &OpenPortfolioEndpoint{
			handler: &StubHandler[OpenPortfolioCommand, PortfolioId]{
				call: func(command OpenPortfolioCommand) (PortfolioId, error) {
					return portfolioId, nil
				},
			},
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, "A Portfolio name")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, endpoint.Open(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, fmt.Sprintf(`{"portfolio_id":"%s"}`, portfolioId) + "\n", rec.Body.String())
		}
	})
}

type StubHandler[K any, V any] struct {
	call func(K) (V, error)
}

func (s *StubHandler[K, V]) Handle(command K) (V, error) {
	return s.call(command)
}
