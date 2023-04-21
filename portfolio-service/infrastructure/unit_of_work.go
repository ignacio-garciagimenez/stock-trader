package infrastructure

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GormUnitOfWork interface {
	Transaction(func (tx *gorm.DB) error, ...*sql.TxOptions) error
}


type FeatureBuilder = func (db *gorm.DB) echo.HandlerFunc

func WithTransaction(uow GormUnitOfWork, builderFunc FeatureBuilder) echo.HandlerFunc {
	return func (context echo.Context) error {
		return uow.Transaction(func(tx *gorm.DB) error {
			return builderFunc(tx)(context)
		})
	}
}