package util

import (
	"github.com/asaskevich/govalidator"
	"github.com/shopspring/decimal"
)

func IsDecimalPositive(i interface{}, _ interface{}) bool {
	switch v := i.(type) {
	case *decimal.Decimal:
		f, _ := v.Float64()
		return govalidator.IsPositive(f)
	default:
		return false
	}
}
