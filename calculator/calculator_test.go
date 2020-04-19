package calculator_test

import (
	"testing"

	"calculator/calculator"
)

func TestCalcString(t *testing.T) {
	cases := []struct {
		Input  string
		Result float64
	}{
		{"    1      +    2", 3},
		{"-4.5 * 2 + 2 * 3", -3},
		{"( 2 + 3 ) * ( 1 + 4 )", 25},
	}

	for _, tC := range cases {
		if r := calculator.CalcString(tC.Input); r != tC.Result {
			t.Errorf("expected calc %s to be %f; got %f",
				tC.Input,
				tC.Result,
				r,
			)
		}
	}
}
