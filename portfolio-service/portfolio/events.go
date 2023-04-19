package portfolio

import (
	"stock-trader/portfolio-service/common"
)

type baseDomainEvent = common.BaseDomainEvent

type PortfolioOpened struct {
	*baseDomainEvent
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
