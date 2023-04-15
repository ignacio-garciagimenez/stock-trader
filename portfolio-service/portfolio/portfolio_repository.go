package portfolio

import (
	"context"
	"errors"
	"stock-trader/portfolio-service/common"
)

type InMemoryPortfolioRepository struct {
	*common.InMemoryBaseRepository[PortfolioId, *Portfolio]
	nameIndex map[string]*Portfolio
}

func (r *InMemoryPortfolioRepository) Save(ctx context.Context, entity *Portfolio) error {
	if err := r.InMemoryBaseRepository.Save(ctx, entity); err != nil {
		return err
	}

	r.nameIndex[entity.name] = entity
	return nil
}

func (r *InMemoryPortfolioRepository) FindByName(ctx context.Context, name string) (*Portfolio, error) {
	if portfolio, found := r.nameIndex[name]; found {
		return portfolio, nil
	}

	return nil, errors.New("entity not found")
}

type PortfolioEntity struct {
	ID string `gorm:"primaryKey;size:36;column:id"`
	name string `gorm:"index;column:name;size:30;not null"`
}
