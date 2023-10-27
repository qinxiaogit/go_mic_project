package handler

import (
	"fmt"
	"math"
)

type Quote struct {
	Dollars uint32
	Cents   uint32
}

func (q Quote) String() string {
	return fmt.Sprintf("$%d.%02d", q.Dollars, q.Cents)
}

func CreateQuoteFromCount(count int) Quote {
	return CreateQuoteFromFloat(8.9)
}

func CreateQuoteFromFloat(Value float64) Quote {
	units, fraction := math.Modf(Value)
	return Quote{
		Dollars: uint32(units),
		Cents:   uint32(math.Trunc(fraction * 100)),
	}
}
