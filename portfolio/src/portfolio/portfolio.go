package portfolio

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"stock-trader/portfolio-context/src/common"
)

func OpenPortfolio(name string) (*Portfolio, error) {
	trimmedName := strings.TrimSpace(name)

	if len(trimmedName) <= 4 || len(trimmedName) >= 31 {
		return nil, errors.New("portfolio name must be between 5 and 30 characters long")
	}

	portfolio := &Portfolio{
		id:   uuid.NewString(),
		name: trimmedName,
	}

	portfolio.domainEvents = append(portfolio.domainEvents, PortfolioOpened{
		BaseDomainEvent: *common.NewBaseDomainEvent("portfolio-opened"),
		portfolioId:     portfolio.Id(),
	})

	return portfolio, nil

}

type Portfolio struct {
	domainEvents []common.DomainEvent
	id           string
	name         string
}

func (p Portfolio) Id() string {
	return p.id
}

func (p Portfolio) Name() string {
	return p.name
}

type PortfolioOpened struct {
	common.BaseDomainEvent
	portfolioId string
}

func (p PortfolioOpened) PortfolioId() string {
	return p.portfolioId
}

func (p PortfolioOpened) EventData() map[string]any {
	return map[string]any{
		"portfolioId": p.portfolioId,
	}
}
