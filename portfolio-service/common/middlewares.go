package common

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GormUnitOfWork interface {
	Transaction(func (tx *gorm.DB) error, ...*sql.TxOptions) error
}

func WithTransaction(uow GormUnitOfWork, handlerBuilder func(db *gorm.DB) echo.HandlerFunc) echo.HandlerFunc {
	return func (context echo.Context) error {
		return uow.Transaction(func(tx *gorm.DB) error {
			return handlerBuilder(tx)(context)
		})
	}
}