package portfolio

import (
	"context"
	"errors"
	"fmt"
	"stock-trader/portfolio-service/common"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type portfolioEntity struct {
	Id   string `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (portfolioEntity) TableName() string {
	return "portfolios"
}

func mapPortfolioEntity(portfolio *Portfolio) (*portfolioEntity, error) {
	if portfolio == nil {
		return nil, errors.New("portfolio cannot be nil")
	}

	return &portfolioEntity{
		Id:   string(portfolio.id),
		Name: portfolio.name,
	}, nil
}

type mySQLPortfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &mySQLPortfolioRepository{
		db: db,
	}
}

func (r *mySQLPortfolioRepository) FindById(ctx context.Context, portfolioId PortfolioId) (*Portfolio, error) {
	return nil, nil
}

func (r *mySQLPortfolioRepository) Save(ctx context.Context, portfolio *Portfolio) error {
	mappedEntity, err := mapPortfolioEntity(portfolio)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", mappedEntity.Id).First(&portfolioEntity{}); result.RowsAffected == 0 {
			if result = tx.Create(mappedEntity); result.RowsAffected == 0 {
				if result.Error != nil {
					if isDuplicate := isDuplicatePortfolioNameError(result.Error); isDuplicate {
						return fmt.Errorf("%w: %s", ErrPortfolioNameAlreadyInUse, portfolio.name)
					}
					return result.Error
				}
			}

		} else {
			if result = tx.Save(mappedEntity); result.Error != nil {
				return result.Error
			}
		}

		for _, domainEvent := range portfolio.domainEvents {
			event := &common.DomainEventEntity{
				Id:        domainEvent.Id(),
				Timestamp: domainEvent.Timestamp(),
				Name:      domainEvent.Name(),
				EventData: datatypes.JSONMap(domainEvent.EventData()),
			}

			if err := tx.Create(event).Error; err != nil {
				return err
			}
		}

		portfolio.ClearDomainEvents()

		return nil
	})

}

func (r *mySQLPortfolioRepository) FindByName(ctx context.Context, portfolioName string) (*Portfolio, error) {
	return nil, nil
}

func isDuplicatePortfolioNameError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") &&
		strings.Contains(err.Error(), "portfolios.idx_name")
}
