package portfolio

import (
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

	repo := NewMySQLPortfolioRepository(db)

	t.Run("given a portfolio should save with domain events", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`portfolio-%s-%s`, randomString(), randomString())
		portfolio, _ := OpenPortfolio(portfolioName)
		domainEvents := portfolio.DomainEvents()

		err := repo.Save(portfolio)

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
		err := repo.Save(nil)

		if assert.Error(t, err) {
			assert.Equal(t, "portfolio cannot be nil", err.Error())
		}
	})

	t.Run("given a repeated portfolio name should return error", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`same-portfolio-name-%s`, randomString())
		portfolio, _ := OpenPortfolio(portfolioName)

		err := repo.Save(portfolio)

		assert.NoError(t, err)

		portfolioWithSameName, _ := OpenPortfolio(portfolioName)

		err = repo.Save(portfolioWithSameName)

		if assert.Error(t, err) {
			var savedEntities []portfolioEntity
			db.Where("name = ?", portfolioName).Find(&savedEntities)

			assert.True(t, len(savedEntities) == 1)
		}
	})

	t.Run("given a portfolio already saved should update it", func(t *testing.T) {

	})

}

func randomString() string {
	x := strings.Split(uuid.NewString(), "-")[0]
	return x
}
