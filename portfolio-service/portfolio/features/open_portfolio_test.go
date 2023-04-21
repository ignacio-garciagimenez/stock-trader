package portfolio_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stock-trader/portfolio-service/infrastructure"
	"stock-trader/portfolio-service/portfolio"
	features "stock-trader/portfolio-service/portfolio/features"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_OpenPortfolioHandler(t *testing.T) {
	t.Run("Open portfolio with empty name", func(t *testing.T) {
		handler := features.NewOpenPortfolioHandler(&StubPortfolioRepository{})

		portfolioId, err := handler.Handle(context.Background(), features.OpenPortfolioCommand{
			Name: "    ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with short name", func(t *testing.T) {
		handler := features.NewOpenPortfolioHandler(&StubPortfolioRepository{})

		portfolioId, err := handler.Handle(context.Background(), features.OpenPortfolioCommand{
			Name: "  Name  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with long name", func(t *testing.T) {
		handler := features.NewOpenPortfolioHandler(&StubPortfolioRepository{})

		portfolioId, err := handler.Handle(context.Background(), features.OpenPortfolioCommand{
			Name: "  Name that is longer than 30 characters  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with name already provided", func(t *testing.T) {
		repo := &StubPortfolioRepository{
			save: func(ctx context.Context, p *portfolio.Portfolio) error {
				return fmt.Errorf("%w: %s", portfolio.ErrPortfolioNameAlreadyInUse, "A portfolio name")
			},
		}
		handler := features.NewOpenPortfolioHandler(repo)

		portfolioId, err := handler.Handle(context.Background(), features.OpenPortfolioCommand{
			Name: "  A portfolio name  ",
		})

		if assert.ErrorIs(t, err, portfolio.ErrPortfolioNameAlreadyInUse) {
			assert.Equal(t, `portfolio name already in use: A portfolio name`, err.Error())
			assert.Empty(t, portfolioId)
		}
	})

	t.Run("Open portfolio successfully", func(t *testing.T) {
		repo := &StubPortfolioRepository{
			save: func(ctx context.Context, p *portfolio.Portfolio) error {
				return nil
			},
		}
		handler := features.NewOpenPortfolioHandler(repo)

		portfolioId, err := handler.Handle(context.Background(), features.OpenPortfolioCommand{
			Name: "  A portfolio name  ",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, portfolioId)
		assert.Equal(t, 1, repo.callsToSave)
	})
}

func Test_OpenPortfolioEndpoint(t *testing.T) {
	t.Run("Open Portfolio Successfully", func(t *testing.T) {
		portfolioId := portfolio.PortfolioId(uuid.NewString())
		endpoint := features.NewOpenPortfolioEndpoint(&StubHandler[features.OpenPortfolioCommand, portfolio.PortfolioId]{
			call: func(ctx context.Context, command features.OpenPortfolioCommand) (portfolio.PortfolioId, error) {
				return portfolioId, nil
			},
		})

		e := echo.New()
		e.Validator = infrastructure.NewRequestValidator()
		req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, "A Portfolio name")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, endpoint.Open(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, fmt.Sprintf(`{"portfolio_id":"%s"}`, portfolioId)+"\n", rec.Body.String())
		}
	})
	t.Run("Open Portfolio With PortfolioWithSameNameAlreadyOpenedError", func(t *testing.T) {
		portfolioId := portfolio.PortfolioId(uuid.NewString())
		portfolioName := "A Portfolio name"
		endpoint := features.NewOpenPortfolioEndpoint(
			&StubHandler[features.OpenPortfolioCommand, portfolio.PortfolioId]{
				call: func(ctx context.Context, command features.OpenPortfolioCommand) (portfolio.PortfolioId, error) {
					return portfolioId, fmt.Errorf("%w: %s", portfolio.ErrPortfolioNameAlreadyInUse, portfolioName)
				},
			},
		)

		e := echo.New()
		e.Validator = infrastructure.NewRequestValidator()
		req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, portfolioName)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := endpoint.Open(c)

		if assert.Error(t, err) {
			err := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusConflict, err.Code)
			assert.Equal(t, `portfolio name already in use: A Portfolio name`, err.Message)
		}
	})
	t.Run("Open Portfolio With Unexpected Error", func(t *testing.T) {
		endpoint := features.NewOpenPortfolioEndpoint(
			&StubHandler[features.OpenPortfolioCommand, portfolio.PortfolioId]{
				call: func(ctx context.Context, command features.OpenPortfolioCommand) (portfolio.PortfolioId, error) {
					return "", errors.New("unexpected error")
				},
			},
		)

		e := echo.New()
		e.Validator = infrastructure.NewRequestValidator()
		req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, "A Portfolio name")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := endpoint.Open(c)

		if assert.Error(t, err) {
			err := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusInternalServerError, err.Code)
			assert.Equal(t, `unexpected error`, err.Message)
		}
	})

	t.Run("Open Portfolio with validation errors", func(t *testing.T) {
		tests := []struct {
			testName           string
			requestBody        string
			validationResponse string
		}{
			{
				testName:           "No Portfolio Name",
				requestBody:        "",
				validationResponse: "Name is a required field",
			},
			{
				testName:           "Portfolio Name too Short",
				requestBody:        fmt.Sprintf(`{"name":"%s"}`, "name"),
				validationResponse: "Name must be greater than 4 characters in length",
			},
			{
				testName:           "Portfolio Name too Long",
				requestBody:        fmt.Sprintf(`{"name":"%s"}`, "a name that is way too looooong"),
				validationResponse: "Name must be less than 31 characters in length",
			},
		}

		endpoint := features.NewOpenPortfolioEndpoint(nil)

		e := echo.New()
		e.Validator = infrastructure.NewRequestValidator()

		for _, tc := range tests {
			t.Run(tc.testName, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(tc.requestBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				err := endpoint.Open(c)

				if assert.Error(t, err) {
					err := err.(*echo.HTTPError)
					assert.Equal(t, http.StatusBadRequest, err.Code)
					assert.Equal(t, &infrastructure.ValidationErrorsResponse{
						Message: "there were validation errors",
						Errors: []infrastructure.FieldError{
							{
								Field: "Name",
								Error: tc.validationResponse,
							},
						},
					}, err.Message)
				}
			})

		}

	})
}

type StubHandler[K any, V any] struct {
	call func(context.Context, K) (V, error)
}

func (s *StubHandler[K, V]) Handle(ctx context.Context, command K) (V, error) {
	return s.call(ctx, command)
}

type StubPortfolioRepository struct {
	callsToSave int
	save        func(context.Context, *portfolio.Portfolio) error
	findById    func(context.Context, portfolio.PortfolioId) (*portfolio.Portfolio, error)
	findByName  func(context.Context, string) (*portfolio.Portfolio, error)
}

func (r *StubPortfolioRepository) Save(ctx context.Context, portfolio *portfolio.Portfolio) error {
	r.callsToSave++
	return r.save(ctx, portfolio)
}

func (r *StubPortfolioRepository) FindById(ctx context.Context, portfolioId portfolio.PortfolioId) (*portfolio.Portfolio, error) {
	return r.findById(ctx, portfolioId)
}

func (r *StubPortfolioRepository) FindByName(ctx context.Context, name string) (*portfolio.Portfolio, error) {
	return r.findByName(ctx, name)
}
