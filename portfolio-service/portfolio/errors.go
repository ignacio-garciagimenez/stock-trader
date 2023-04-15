package portfolio

import "fmt"

type PortfolioWithSameNameAlreadyOpened struct {
	portfolioName string
}

func (e PortfolioWithSameNameAlreadyOpened) Error() string {
	return fmt.Sprintf(`a portfolio with name "%s" was already opened`, e.portfolioName)
}

func NewPortfolioWithSameNameAlreadyOpened(name string) error {
	return &PortfolioWithSameNameAlreadyOpened{
		portfolioName: name,
	}
}
