package utils

import (
	"github.com/elliotchance/pie/v2"
)

func FindFirst[T any](ss []T, fn func(value T) bool) *T {
	index := pie.FindFirstUsing(ss, fn)

	if index == -1 {
		return nil
	}

	return &ss[index]
}
