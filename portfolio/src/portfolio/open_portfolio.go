package portfolio

type OpenPortfolioHandler struct {
	portfolioRepository PortfolioRepository
}

type OpenPortfolioCommand struct {
	name string
}

func (h *OpenPortfolioHandler) Open(command OpenPortfolioCommand) (PortfolioId, error) {
	portfolio, err := OpenPortfolio(command.name)
	if err != nil {
		return "", err
	}

	h.portfolioRepository.Save(portfolio)

	return portfolio.id, nil
}
