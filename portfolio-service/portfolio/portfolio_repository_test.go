package portfolio_test

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-trader/portfolio-service/common"
	"stock-trader/portfolio-service/infrastructure"
	"stock-trader/portfolio-service/portfolio"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSavePortfolio(t *testing.T) {
	db, _ := infrastructure.ConnectDB()

	repo := portfolio.NewPortfolioRepository(db)

	t.Run("given a portfolio should save with domain events", func(t *testing.T) {
		portfolioName := fmt.Sprintf(`portfolio-%s-%s`, randomString(), randomString())
		portfolio, _ := portfolio.OpenPortfolio(portfolioName)
		domainEvents := portfolio.DomainEvents()

		err := repo.Save(context.Background(), portfolio)

		if assert.NoError(t, err) {
			var savedPortfolio []map[string]any
			result := db.Raw("SELECT id, name FROM portfolios WHERE id = ?", portfolio.Id()).Scan(&savedPortfolio)

			if assert.True(t, result.RowsAffected == 1) {
				assert.Equal(t, portfolioName, savedPortfolio[0]["name"])
				assert.Equal(t, string(portfolio.Id()), savedPortfolio[0]["id"])
			}
			
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
		newPortfolio, _ := portfolio.OpenPortfolio(portfolioName)

		err := repo.Save(context.Background(), newPortfolio)

		assert.NoError(t, err)

		portfolioWithSameName, _ := portfolio.OpenPortfolio(portfolioName)
		duplicateDomainEvents := portfolioWithSameName.DomainEvents()

		err = repo.Save(context.Background(), portfolioWithSameName)

		if assert.ErrorIs(t, err, portfolio.ErrPortfolioNameAlreadyInUse) {
			assert.Equal(t, fmt.Sprintf(`portfolio name already in use: %s`, portfolioName), err.Error())

			var savedPortfolio []map[string]any
			result := db.Raw("SELECT id, name FROM portfolios WHERE name = ?", portfolioName).Scan(&savedPortfolio)

			if assert.True(t, result.RowsAffected == 1) {
				assert.Equal(t, string(newPortfolio.Id()), savedPortfolio[0]["id"])
			}

			result = db.Raw("SELECT * FROM event_journal WHERE id = ?", duplicateDomainEvents[0].Id()).Scan([]any{})
			assert.True(t, result.RowsAffected == 0)
		}
	})

	t.Run("given a portfolio already saved should update it", func(t *testing.T) {
		//TODO
	})

}

func randomString() string {
	return strings.Split(uuid.NewString(), "-")[0]
}
