package trigonometry

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_sin(t *testing.T) {
	scenarios := map[int64]float64{
		0:   0,
		15:  0.258819045,
		30:  0.5,
		45:  0.707106781,
		60:  0.866025404,
		75:  0.965925826,
		90:  1,
		105: 0.965925826,
		120: 0.866025404,
		135: 0.707106781,
		150: 0.5,
		165: 0.258819045,
		180: 0,
		195: -0.258819045,
		210: -0.5,
		225: -0.707106781,
		240: -0.866025404,
		255: -0.965925826,
		270: -1,
		285: -0.965925826,
		300: -0.866025404,
		315: -0.707106781,
		330: -0.5,
		345: -0.258819045,
		360: 0,
		375: 0.258819045,
		390: 0.5,
	}

	for degree, expected := range scenarios {
		res := Sin(decimal.New(degree, 0))
		diff := res.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"degree   : %d\n"+
				"expected : %f\n"+
				"get      : %s",
				degree, expected, res.String())
		}
	}
}

func Test_cos(t *testing.T) {
	scenarios := map[int64]float64{
		0:   1,
		15:  0.965925826,
		30:  0.866025404,
		45:  0.707106781,
		60:  0.5,
		75:  0.258819045,
		90:  0,
		105: -0.258819045,
		120: -0.5,
		135: -0.707106781,
		150: -0.866025404,
		165: -0.965925826,
		180: -1,
		195: -0.965925826,
		210: -0.866025404,
		225: -0.707106781,
		240: -0.5,
		255: -0.258819045,
		270: 0,
		285: 0.258819045,
		300: 0.5,
		315: 0.707106781,
		330: 0.866025404,
		345: 0.965925826,
		360: 1,
		375: 0.965925826,
		390: 0.866025404,
	}

	for degree, expected := range scenarios {
		res := Cos(decimal.New(degree, 0))
		diff := res.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"degree   : %d\n"+
				"expected : %f\n"+
				"get      : %s",
				degree, expected, res.String())
		}
	}
}

func Test_tan(t *testing.T) {
	scenarios := map[int64]float64{
		0:   0,
		15:  0.267949192,
		30:  0.577350269,
		45:  1,
		60:  1.732050808,
		75:  3.732050808,
		105: -3.732050808,
		120: -1.732050808,
		135: -1,
		150: -0.577350269,
		165: -0.267949192,
		180: 0,
		195: 0.267949192,
		210: 0.577350269,
		225: 1,
		240: 1.732050808,
		255: 3.732050808,
		285: -3.732050808,
		300: -1.732050808,
		315: -1,
		330: -0.577350269,
		345: -0.267949192,
		360: 0,
		375: 0.267949192,
		390: 0.577350269,
	}

	for degree, expected := range scenarios {
		res := Tan(decimal.New(degree, 0))
		diff := res.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"degree   : %d\n"+
				"expected : %f\n"+
				"get      : %s",
				degree, expected, res.String())
		}
	}
}

func Test_acos(t *testing.T) {
	scenarios := map[float64]float64{
		1:            0,
		0.965925826:  15,
		0.866025404:  30,
		0.707106781:  45,
		0.5:          60,
		0.258819045:  75,
		0:            90,
		-0.258819045: 105,
		-0.5:         120,
		-0.707106781: 135,
		-0.866025404: 150,
		-0.965925826: 165,
		-1:           180,
	}

	for cosValue, expected := range scenarios {
		res := Acos(decimal.NewFromFloat(cosValue))
		diff := res.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"cos value : %f\n"+
				"expected  : %f\n"+
				"get       : %s",
				cosValue, expected, res.String())
		}
	}
}

func Test_acot(t *testing.T) {
	scenarios := map[float64]float64{3.732050808: 15,
		1.732050808:  30,
		1:            45,
		0.577350269:  60,
		0.267949192:  75,
		0:            90,
		-0.267949192: 105,
		-0.577350269: 120,
		-1:           135,
		-1.732050808: 150,
		-3.732050808: 165,
	}

	for cotValue, expected := range scenarios {
		res := Acot(decimal.NewFromFloat(cotValue))
		diff := res.Sub(decimal.NewFromFloat(expected))

		if !diff.Round(3).Equal(decimal.Zero) {
			t.Errorf("\n"+
				"cot value : %f\n"+
				"expected  : %f\n"+
				"get       : %s",
				cotValue, expected, res.String())
		}
	}
}
