package portfolio

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler[K any, V any] interface {
	Handle(K) (V, error)
}

type OpenPortfolioEndpoint struct {
	handler Handler[OpenPortfolioCommand, PortfolioId]
}

func (e *OpenPortfolioEndpoint) Open(c echo.Context) error {
	command := new(OpenPortfolioCommand)
	if err := c.Bind(command); err != nil {
		return err
	}

	portfolioId, err := e.handler.Handle(*command)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, struct {
		PortfolioId PortfolioId `json:"portfolio_id"`
	}{
		PortfolioId: portfolioId,
	})
}

type OpenPortfolioHandler struct {
	portfolioRepository PortfolioRepository
}

type OpenPortfolioCommand struct {
	Name string `json:"name"`
}

func (h *OpenPortfolioHandler) Handle(command OpenPortfolioCommand) (PortfolioId, error) {
	portfolio, err := OpenPortfolio(command.Name)
	if err != nil {
		return "", err
	}

	h.portfolioRepository.Save(portfolio)

	return portfolio.id, nil
}
