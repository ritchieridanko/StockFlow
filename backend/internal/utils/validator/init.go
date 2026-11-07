package validator

import (
	"strconv"
	"unicode/utf8"

	"github.com/ritchieridanko/stockflow/backend/internal/constants"
)

type Validator struct{}

func Init() *Validator {
	return &Validator{}
}

func (v *Validator) Name(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < nameMinLength {
		return false, "Name is required"
	}
	return true, ""
}

func (v *Validator) Price(value int) (bool, string) {
	if value < minPrice {
		return false, "Price is invalid: " + strconv.Itoa(value)
	}
	return true, ""
}

func (v *Validator) Quantity(value int) (bool, string) {
	if value <= minQuantity {
		return false, "Quantity is invalid: " + strconv.Itoa(value)
	}
	return true, ""
}

func (v *Validator) SKU(value string) (bool, string) {
	length := utf8.RuneCountInString(value)
	if length < skuMinLength {
		return false, "SKU is required"
	}
	return true, ""
}

func (v *Validator) StockUpdateType(value string) (bool, string) {
	if value != constants.StockUpdateTypeIn && value != constants.StockUpdateTypeOut {
		return false, "Stock update type is invalid: " + value
	}
	return true, ""
}
