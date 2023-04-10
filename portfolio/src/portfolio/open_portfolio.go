package portfolio

import (
	"net/http"
	"stock-trader/portfolio-context/src/common"

	"github.com/labstack/echo/v4"
)

type OpenPortfolioHandler struct {
	portfolioRepository PortfolioRepository
}

type OpenPortfolioCommand struct {
	Name string `json:"name"`
}

type OpenPortfolioEndpoint struct {
	handler common.Handler[OpenPortfolioCommand, PortfolioId]
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

func (h *OpenPortfolioHandler) Handle(command OpenPortfolioCommand) (PortfolioId, error) {
	portfolio, err := OpenPortfolio(command.Name)
	if err != nil {
		return "", err
	}

	h.portfolioRepository.Save(portfolio)

	return portfolio.id, nil
}
