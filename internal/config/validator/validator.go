package validator

import (
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"go.uber.org/dig"
)

const (
	decimalRegexString = "^-?\\d+(\\.\\d+)?$"
)

var (
	decimalRegex = regexp.MustCompile(decimalRegexString)
)

// Validator object of validator for app
var Validator *validator.Validate

// Initialize initializes app validator
func Initialize(container *dig.Container) {
	Validator = binding.Validator.Engine().(*validator.Validate)
	if err := container.Invoke(registerValidations); err != nil {
		panic("cannot register validations: " + err.Error())
	}
}

func registerValidations() {
	_ = Validator.RegisterValidation("decimal", decimalValid)
	_ = Validator.RegisterValidation("decimalGT", decimalGreaterThan)
}

func decimalValid(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldKind := field.Kind()
	if fieldKind != reflect.String {
		return false
	}
	fieldStr := field.Interface().(string)
	return decimalRegex.MatchString(fieldStr)
}

// GreaterThan checks if passed decimal value is greater than specified.
// Validator usage: "YOUR_TAG=DECIMAL_NUMBER", e.g. "decimalGT=0"
func decimalGreaterThan(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()
	fieldStr := field.Interface().(string)
	dec, err := decimal.NewFromString(fieldStr)
	if err != nil {
		return false
	}
	decParam, _ := decimal.NewFromString(param)
	return dec.GreaterThan(decParam)
}
