package portfolio

import (
	"context"
	"errors"
	"stock-trader/portfolio-service/common"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type inMemoryBaseRepository = common.InMemoryBaseRepository[PortfolioId, *Portfolio]

type InMemoryPortfolioRepository struct {
	*inMemoryBaseRepository
	nameIndex map[string]*Portfolio
}

func NewInMemoryPortfolioRepository() *InMemoryPortfolioRepository {
	return &InMemoryPortfolioRepository{
		inMemoryBaseRepository: common.NewInMemoryBaseRepository[PortfolioId, *Portfolio](),
		nameIndex:              map[string]*Portfolio{},
	}
}

func (r *InMemoryPortfolioRepository) Save(ctx context.Context, entity *Portfolio) error {
	if err := r.inMemoryBaseRepository.Save(ctx, entity); err != nil {
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

func NewMySQLPortfolioRepository(db *gorm.DB) *mySQLPortfolioRepository {
	return &mySQLPortfolioRepository{
		db: db,
	}
}

func (r *mySQLPortfolioRepository) FindById(portfolioId PortfolioId) (*Portfolio, error) {
	return nil, nil
}

func (r *mySQLPortfolioRepository) Save(portfolio *Portfolio) error {
	mappedEntity, err := mapPortfolioEntity(portfolio)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Where("id = ?", mappedEntity.Id).FirstOrCreate(&mappedEntity); result.RowsAffected == 0 {
			if result.Error != nil {
				return result.Error
			}

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

func (r *mySQLPortfolioRepository) findByName(portfolioName string) (*Portfolio, error) {
	return nil, nil
}
