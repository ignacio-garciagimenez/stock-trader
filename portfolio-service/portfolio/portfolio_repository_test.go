package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-trader/portfolio-service/common"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSavePortfolio(t *testing.T) {
	db, _ := common.ConnectDB()

	repo := NewPortfolioRepository(db)

	t.Run("given a portfolio should save with domain events", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`portfolio-%s-%s`, randomString(), randomString())
		portfolio, _ := OpenPortfolio(portfolioName)
		domainEvents := portfolio.DomainEvents()

		err := repo.Save(context.Background(), portfolio)

		if assert.NoError(t, err) {
			var savedEntity portfolioEntity
			db.Raw("SELECT id, name FROM portfolios WHERE id = ?", portfolio.id).Scan(&savedEntity)

			assert.Equal(t, portfolioName, savedEntity.Name)
			assert.Equal(t, string(portfolio.id), savedEntity.Id)

			for _, event := range domainEvents {
				var savedEvent common.DomainEventEntity
				db.Raw("SELECT id, timestamp, name, event_data FROM event_journal WHERE id = ?", event.Id()).Scan(&savedEvent)
				assert.Equal(t, event.Id(), savedEvent.Id)
				assert.Equal(t, event.Timestamp().Format(time.RFC3339), savedEvent.Timestamp.Format(time.RFC3339))
				assert.Equal(t, event.Name(), savedEvent.Name)

				b, _ := savedEvent.EventData.MarshalJSON()
				var eventData map[string]any
				json.Unmarshal(b, &eventData)
				assert.Equal(t, event.EventData(), eventData)
			}
		}
	})

	t.Run("given a nil portfolio should return error", func(t *testing.T) {
		err := repo.Save(context.Background(), nil)

		if assert.Error(t, err) {
			assert.Equal(t, "portfolio cannot be nil", err.Error())
		}
	})

	t.Run("given a repeated portfolio name should return error and not save domain events", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`same-portfolio-name-%s`, randomString())
		portfolio, _ := OpenPortfolio(portfolioName)

		err := repo.Save(context.Background(), portfolio)

		assert.NoError(t, err)

		portfolioWithSameName, _ := OpenPortfolio(portfolioName)

		err = repo.Save(context.Background(), portfolioWithSameName)

		if assert.Error(t, err) {
			customErr, ok := err.(*PortfolioWithSameNameAlreadyOpened)
			assert.True(t, ok)
			assert.Equal(t, fmt.Sprintf(`a portfolio with name "%s" was already opened`, portfolioName), customErr.Error())

			var savedEntities []portfolioEntity
			db.Where("name = ?", portfolioName).Find(&savedEntities)

			assert.True(t, len(savedEntities) == 1)
			assert.Equal(t, string(portfolio.id), savedEntities[0].Id)

			result := db.Raw("SELECT * FROM event_journal WHERE id = ?", portfolioWithSameName.domainEvents[0].Id()).Scan([]any{})
			assert.True(t, result.RowsAffected == 0)
		}
	})

	t.Run("given a portfolio already saved should update it", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`same-portfolio-name-%s`, randomString())
		portfolio, _ := OpenPortfolio(portfolioName)

		err := repo.Save(context.Background(), portfolio)

		assert.NoError(t, err)

		portfolio.name = fmt.Sprintf(`another-name-%s`, randomString())

		err = repo.Save(context.Background(), portfolio)

		if assert.NoError(t, err) {
			var savedEntities []portfolioEntity
			db.Where("id = ?", portfolio.id).Find(&savedEntities)

			assert.True(t, len(savedEntities) == 1)
			assert.Equal(t, string(portfolio.id), savedEntities[0].Id)
			assert.Equal(t, portfolio.name, savedEntities[0].Name)
		}
	})

}

func randomString() string {
	return strings.Split(uuid.NewString(), "-")[0]
}
