package portfolio

type OpenPortfolioHandler struct {
	portfolioRepository PortfolioRepository
}

type OpenPortfolioCommand struct {
	name string
}

func (h *OpenPortfolioHandler) Open(command OpenPortfolioCommand) (string, error) {
	portfolio, err := OpenPortfolio(command.name)
	if err != nil {
		return "", err
	}

	return portfolio.id, nil
}
