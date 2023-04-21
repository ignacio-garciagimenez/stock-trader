package infrastructure

import (
	"database/sql"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestWithTransaction(t *testing.T) {
	mock := &MockGormUnitOfWork{}
	handlerCalled := false

	handler := WithTransaction(mock, func(db *gorm.DB) echo.HandlerFunc {
		return func(c echo.Context) error {
			handlerCalled = true
			return nil
		}
	})

	assert.False(t, handlerCalled)
	assert.False(t, mock.called)

	err := handler(nil)

	if assert.NoError(t, err) {
		assert.True(t, handlerCalled)
		assert.True(t, mock.called)
	}

	
}

type MockGormUnitOfWork struct{
	called bool
}

func (m *MockGormUnitOfWork) Transaction(f func(tx *gorm.DB) error, _ ...*sql.TxOptions) error {
	m.called = true
	return f(nil)
}
