package portfolio

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"stock-trader/portfolio-context/src/common"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_OpenPortfolioHandler(t *testing.T) {
	t.Run("Open portfolio with empty name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "    ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with short name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  Name  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with long name", func(t *testing.T) {
		handler := &OpenPortfolioHandler{
			portfolioRepository: &InMemoryPortfolioRepository{},
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  Name that is longer than 30 characters  ",
		})

		assert.Empty(t, portfolioId)
		assert.Error(t, err)
		assert.Equal(t, "portfolio name must be between 5 and 30 characters long", err.Error())
	})

	t.Run("Open portfolio with name already provided", func(t *testing.T) {
		repo := &StubPortfolioRepository{
			save: func(p *Portfolio) error {
				return NewPortfolioWithSameNameAlreadyOpened("A portfolio name")
			},
		}
		handler := &OpenPortfolioHandler{
			portfolioRepository: repo,
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  A portfolio name  ",
		})

		if assert.Error(t, err) {
			assert.IsType(t, &PortfolioWithSameNameAlreadyOpened{}, err)
			assert.Equal(t, `a portfolio with name "A portfolio name" was already opened`, err.Error())
			assert.Empty(t, portfolioId)
		}
	})

	t.Run("Open portfolio successfully", func(t *testing.T) {
		repo := &InMemoryPortfolioRepository{
			InMemoryBaseRepository: common.NewInMemoryBaseRepository[PortfolioId, *Portfolio](),
			nameIndex:              map[string]*Portfolio{},
		}
		handler := &OpenPortfolioHandler{
			portfolioRepository: repo,
		}

		portfolioId, err := handler.Handle(OpenPortfolioCommand{
			Name: "  A portfolio name  ",
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, portfolioId)
		assert.NotNil(t, func() *Portfolio {
			portfolio, _ := repo.FindById(portfolioId)
			return portfolio
		}())
	})
}

func Test_OpenPortfolioEndpoint(t *testing.T) {
	t.Run("Open Portfolio Successfully", func(t *testing.T) {
		portfolioId := PortfolioId(uuid.NewString())
		endpoint := &OpenPortfolioEndpoint{
			handler: &StubHandler[OpenPortfolioCommand, PortfolioId]{
				call: func(command OpenPortfolioCommand) (PortfolioId, error) {
					return portfolioId, nil
				},
			},
		}

		e := echo.New()
		e.Validator = common.NewRequestValidator()
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
		portfolioId := PortfolioId(uuid.NewString())
		endpoint := &OpenPortfolioEndpoint{
			handler: &StubHandler[OpenPortfolioCommand, PortfolioId]{
				call: func(command OpenPortfolioCommand) (PortfolioId, error) {
					return portfolioId, NewPortfolioWithSameNameAlreadyOpened("A Portfolio name")
				},
			},
		}

		e := echo.New()
		e.Validator = common.NewRequestValidator()
		req := httptest.NewRequest(http.MethodPost, "/portfolios", strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, "A Portfolio name")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := endpoint.Open(c)

		if assert.Error(t, err) {
			err := err.(*echo.HTTPError)
			assert.Equal(t, http.StatusConflict, err.Code)
			assert.Equal(t, `a portfolio with name "A Portfolio name" was already opened`, err.Message)
		}
	})
	t.Run("Open Portfolio With Unexpected Error", func(t *testing.T) {
		endpoint := &OpenPortfolioEndpoint{
			handler: &StubHandler[OpenPortfolioCommand, PortfolioId]{
				call: func(command OpenPortfolioCommand) (PortfolioId, error) {
					return "", errors.New("unexpected error")
				},
			},
		}

		e := echo.New()
		e.Validator = common.NewRequestValidator()
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

		endpoint := &OpenPortfolioEndpoint{
			handler: nil,
		}

		e := echo.New()
		e.Validator = common.NewRequestValidator()

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
					assert.Equal(t, &common.ValidationErrorsResponse{
						Message: "there were validation errors",
						Errors: []common.FieldError{
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
	call func(K) (V, error)
}

func (s *StubHandler[K, V]) Handle(command K) (V, error) {
	return s.call(command)
}

type StubPortfolioRepository struct {
	save       func(*Portfolio) error
	findById   func(PortfolioId) (*Portfolio, error)
	findByName func(string) (*Portfolio, error)
}

func (r *StubPortfolioRepository) Save(portfolio *Portfolio) error {
	return r.save(portfolio)
}

func (r *StubPortfolioRepository) FindById(portfolioId PortfolioId) (*Portfolio, error) {
	return r.findById(portfolioId)
}

func (r *StubPortfolioRepository) FindByName(name string) (*Portfolio, error) {
	return r.findByName(name)
}
