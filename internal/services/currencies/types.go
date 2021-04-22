package currencies

import "github.com/Confialink/wallet-currencies/internal/models"

const (
	TypeFiat   = "fiat"
	TypeCrypto = "crypto"
	TypeOther  = "other"
)

var knownTypes = []CurrencyType{
	TypeFiat,
	TypeCrypto,
	TypeOther,
}

type CurrencyType string

type CodeValidator interface {
	//IsValid verifies if currency code is well formed and valid
	IsValid(currencyType CurrencyType, code string) error
	//IsTypeSupported indicates if this validator is intended to check given currency type
	IsTypeSupported(currencyType CurrencyType) bool
}

func (c CurrencyType) String() string {
	return string(c)
}

//GetKnownTypes returns a list of known currency types
func GetKnownTypes() []CurrencyType {
	result := make([]CurrencyType, len(knownTypes))
	for i, knownFeed := range knownTypes {
		result[i] = knownFeed
	}
	return result
}

//IsKnownType verifies if given currency type is known
func IsKnownType(currencyType CurrencyType) bool {
	for _, knownType := range knownTypes {
		if currencyType == knownType {
			return true
		}
	}
	return false
}

//IsKnownTypeString verifies if given currency type string is known
func IsKnownTypeString(currencyType string) bool {
	return IsKnownType(CurrencyType(currencyType))
}

type Manager interface {
	AddCurrency(currency Currency) (*models.Currency, error)
	DisableCurrency(currencyCode string) (*models.Currency, error)
}

type Currency struct {
	Code          string       `json:"code" binding:"required"`
	DecimalPlaces uint8        `json:"decimalPlaces" binding:"required"`
	Type          CurrencyType `json:"type" binding:"required"`
	IsActive      bool         `json:"isActive"`
	Name          *string      `json:"name" binding:"omitempty,max=128"`
}
