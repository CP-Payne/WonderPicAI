package validation

import "github.com/go-playground/validator/v10"

func New() *validator.Validate {
	validate := validator.New()

	err := validate.RegisterValidation("passwordcomplexity", validatePasswordComplexity)
	if err != nil {
		panic("Failed to register passwordcomplexity validator: " + err.Error())
	}

	return validate
}
