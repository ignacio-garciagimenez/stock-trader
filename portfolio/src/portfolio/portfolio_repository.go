package portfolio

import (
	"errors"
	"stock-trader/portfolio-context/src/common"
)

type InMemoryPortfolioRepository struct {
	*common.InMemoryBaseRepository[string, *Portfolio]
	nameIndex map[string]*Portfolio
}

func (r *InMemoryPortfolioRepository) Save(entity *Portfolio) error {
	if err := r.InMemoryBaseRepository.Save(entity); err != nil {
		return err
	}

	r.nameIndex[entity.name] = entity
	return nil
}

func (r *InMemoryPortfolioRepository) FindByName(name string) (*Portfolio, error) {
	if portfolio, found := r.nameIndex[name]; found {
		return portfolio, nil
	}

	return nil, errors.New("entity not found")
}
