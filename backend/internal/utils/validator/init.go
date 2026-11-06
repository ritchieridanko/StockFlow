package validator

import (
	"strconv"
	"unicode/utf8"
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
		return false, "Price is invalid: " + strconv.FormatInt(int64(value), 10)
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
