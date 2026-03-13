package utils

import (
	"fmt"

	"github.com/NETWAYS/go-check"
)

func MakeThreashold(lower *float64, upper *float64) *check.Threshold {
	if lower != nil && upper != nil {
		return &check.Threshold{Lower: *lower, Upper: *upper}
	} else {
		return nil
	}
}

func FormatFloat(value *float64) string {
	if value != nil {
		return fmt.Sprintf("%.0f", *value)
	} else {
		return "nil"
	}
}
