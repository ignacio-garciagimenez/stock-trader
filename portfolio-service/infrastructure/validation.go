package infrastructure

import (
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

type ValidationErrorsResponse struct {
	Message string       `json:"message"`
	Errors  []FieldError `json:"validation_errors"`
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type requestValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func NewRequestValidator() *requestValidator {
	trans := newTranslator()
	validator := validator.New()
	en_translations.RegisterDefaultTranslations(validator, trans)
	return &requestValidator{
		validator: validator,
		trans:     trans,
	}
}

func newTranslator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	return trans
}

func (rv *requestValidator) Validate(i any) error {
	if err := rv.validator.Struct(i); err != nil {
		resp := rv.createValidationErrorResponse(err.(validator.ValidationErrors))

		return echo.NewHTTPError(http.StatusBadRequest, resp)
	}
	return nil
}

func (rv *requestValidator) createValidationErrorResponse(validationErrors validator.ValidationErrors) *ValidationErrorsResponse {
	errorResponse := &ValidationErrorsResponse{Message: "there were validation errors"}

	for _, fieldErr := range validationErrors {
		errorResponse.Errors = append(errorResponse.Errors, FieldError{
			Field: fieldErr.Field(),
			Error: fieldErr.Translate(rv.trans),
		})
	}

	return errorResponse
}
