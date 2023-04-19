package portfolio

import (
	"context"
	"errors"
	"strings"

	"stock-trader/portfolio-service/common"

	"github.com/google/uuid"
)

type PortfolioRepository interface {
	common.Repository[PortfolioId, *Portfolio]
	FindByName(context.Context, string) (*Portfolio, error)
}

type PortfolioId string

type Portfolio struct {
	domainEvents []common.DomainEvent
	id           PortfolioId
	name         string
}

func OpenPortfolio(name string) (*Portfolio, error) {
	trimmedName := strings.TrimSpace(name)

	if len(trimmedName) <= 4 || len(trimmedName) >= 31 {
		return nil, errors.New("portfolio name must be between 5 and 30 characters long")
	}

	portfolio := &Portfolio{
		id:   PortfolioId(uuid.NewString()),
		name: trimmedName,
	}

	portfolio.domainEvents = append(portfolio.domainEvents, PortfolioOpened{
		baseDomainEvent: common.NewBaseDomainEvent("portfolio-opened"),
		portfolioId:     string(portfolio.id),
	},)

	return portfolio, nil

}

func (p Portfolio) Id() PortfolioId {
	return p.id
}

func (p Portfolio) Name() string {
	return p.name
}

func (p Portfolio) DomainEvents() []common.DomainEvent {
	output := []common.DomainEvent{}
	for _, value := range p.domainEvents {
		output = append(output, value)
	}
	return output
}

func (p *Portfolio) ClearDomainEvents() {
	p.domainEvents = []common.DomainEvent{}
}
