// Copyright (c) 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package monautil_test

import (
	"math/big"
	"testing"

	. "github.com/monasuite/monautil"

	//"github.com/shopspring/decimal"
)

func TestAmountCreation(t *testing.T) {
	var bigAmt Amount
	var smallAmt Amount
	bigAmt, _ = bigAmt.SetString("10512000000000001")
	smallAmt, _ = smallAmt.SetString("-10512000000000001")
	tests := []struct {
		name     string
		amount   string
		valid    bool
		expected Amount
	}{
		// Positive tests.
		{
			name:     "zero",
			amount:   string("0"),
			valid:    true,
			expected: Amount(*big.NewInt(0)),
		},
		{
			name:     "max producible",
			amount:   string("105.12e6"),
			valid:    true,
			expected: Amount(*big.NewInt(MaxSatoshi)),
		},
		{
			name:     "min producible",
			amount:   string("-105.12e6"),
			valid:    true,
			expected: Amount(*big.NewInt(-MaxSatoshi)),
		},
		{
			name:     "exceeds max producible",
			amount:   string("105120000.00000001"),
			valid:    true,
			expected: bigAmt,
		},
		{
			name:     "exceeds min producible",
			amount:   string("-105120000.00000001"),
			valid:    true,
			expected: smallAmt,
		},
		{
			name:     "one hundred",
			amount:   string("100"),
			valid:    true,
			expected: Amount(*big.NewInt(100 * SatoshiPerBitcoin)),
		},
		{
			name:     "fraction",
			amount:   string("0.01234567"),
			valid:    true,
			expected: Amount(*big.NewInt(1234567)),
		},
		{
			name:     "rounding up",
			amount:   string("54.999999999999943157"),
			valid:    true,
			expected: Amount(*big.NewInt(55 * SatoshiPerBitcoin)),
		},
		{
			name:     "rounding down",
			amount:   string("55.000000000000056843"),
			valid:    true,
			expected: Amount(*big.NewInt(55 * SatoshiPerBitcoin)),
		},

		// Negative tests.
		// {
		// big.float can't convert NaN.
		// 	name:   "not-a-number",
		// 	amount: big.NewFloat(math.NaN()),
		// 	valid:  false,
		// },
		{
			name:   "-infinity",
			amount: string("105120000000000000001"),
			valid:  false,
		},
		{
			name:   "+infinity",
			amount: string("-105120000000000000001"),
			valid:  false,
		},
	}

	for _, test := range tests {
		a, err := NewAmount(test.amount)
		switch {
		case test.valid && err != nil:
			t.Errorf("%v: Positive test Amount creation failed with: %v", test.name, err)
			continue
		case !test.valid && err == nil:
			t.Errorf("%v: Negative test Amount creation succeeded (value %v) when should fail", test.name, a)
			continue
		}

		if a.Cmp(test.expected) != 0 {
			t.Errorf("%v: Created amount %v does not match expected %v", test.name, a, test.expected)
			continue
		}
	}
}

func TestAmountUnitConversions(t *testing.T) {
	tests := []struct {
		name      string
		amount    Amount
		unit      AmountUnit
		converted string
		s         string
	}{
		{
			name:      "MMONA",
			amount:    Amount(*big.NewInt(MaxSatoshi)),
			unit:      AmountMegaBTC,
			converted: string("105.12"),
			s:         "105.12 MMONA",
		},
		{
			name:      "kMONA",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountKiloBTC,
			converted: string("444.333222111"),
			s:         "444.333222111 kMONA",
		},
		{
			name:      "MONA",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountBTC,
			converted: string("444333.222111"),
			s:         "444333.222111 MONA",
		},
		{
			name:      "mMONA",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountMilliBTC,
			converted: string("444333222.111"),
			s:         "444333222.111 mMONA",
		},
		{

			name:      "μMONA",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountMicroBTC,
			converted: string("444333222111"),
			s:         "444333222111 μMONA",
		},
		{

			name:      "Watanabe",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountSatoshi,
			converted: string("44433322211100"),
			s:         "44433322211100 Watanabe",
		},
		{

			name:      "non-standard unit",
			amount:    Amount(*big.NewInt(44433322211100)),
			unit:      AmountUnit(-1),
			converted: string("4443332.22111"),
			s:         "4443332.22111 1e-1 MONA",
		},
	}

	for _, test := range tests {
		f := test.amount.ToUnit(test.unit)
		// float64->big.float is a margin of error. So......
		//decimalconvertedtest, _ := decimal.NewFromString(test.converted)
		if f != test.converted {
			t.Errorf("%v: converted value %v does not match expected %v", test.name, f, test.converted)
			continue
		}

		s := test.amount.Format(test.unit)
		if s != test.s {
			t.Errorf("%v: format '%v' does not match expected '%v'", test.name, s, test.s)
			continue
		}

		// Verify that Amount.ToBTC works as advertised.
		f1 := test.amount.ToUnit(AmountBTC)
		f2 := test.amount.ToBTC()
		if f1 != f2 {
			t.Errorf("%v: ToBTC does not match ToUnit(AmountBTC): %v != %v", test.name, f1, f2)
		}

		// Verify that Amount.String works as advertised.
		s1 := test.amount.Format(AmountBTC)
		s2 := test.amount.String()
		if s1 != s2 {
			t.Errorf("%v: String does not match Format(AmountBitcoin): %v != %v", test.name, s1, s2)
		}
	}
}

func TestAmountMulF64(t *testing.T) {
	var amt Amount
	tests := []struct {
		name string
		amt  Amount
		mul  float64
		res  Amount
	}{
		{
			name: "Multiply 0.1 BTC by 2",
			amt:  amt.SetFloat(100e5), // 0.1 BTC
			mul:  2,
			res:  amt.SetFloat(200e5), // 0.2 BTC
		},
		{
			name: "Multiply 0.2 BTC by 0.02",
			amt:  amt.SetFloat(200e5), // 0.2 BTC
			mul:  1.02,
			res:  amt.SetFloat(204e5), // 0.204 BTC
		},
		{
			name: "Multiply 0.1 BTC by -2",
			amt:  amt.SetFloat(100e5), // 0.1 BTC
			mul:  -2,
			res:  amt.SetFloat(-200e5), // -0.2 BTC
		},
		{
			name: "Multiply 0.2 BTC by -0.02",
			amt:  amt.SetInt(200e5), // 0.2 BTC
			mul:  -1.02,
			res:  amt.SetInt(-204e5), // -0.204 BTC
		},
		{
			name: "Multiply -0.1 BTC by 2",
			amt:  amt.SetInt(-100e5), // -0.1 BTC
			mul:  2,
			res:  amt.SetInt(-200e5), // -0.2 BTC
		},
		{
			name: "Multiply -0.2 BTC by 0.02",
			amt:  amt.SetInt(-200e5), // -0.2 BTC
			mul:  1.02,
			res:  amt.SetInt(-204e5), // -0.204 BTC
		},
		{
			name: "Multiply -0.1 BTC by -2",
			amt:  amt.SetInt(-100e5), // -0.1 BTC
			mul:  -2,
			res:  amt.SetInt(200e5), // 0.2 BTC
		},
		{
			name: "Multiply -0.2 BTC by -0.02",
			amt:  amt.SetInt(-200e5), // -0.2 BTC
			mul:  -1.02,
			res:  amt.SetInt(204e5), // 0.204 BTC
		},
		{
			name: "Round down",
			amt:  amt.SetInt(49), // 49 Satoshis
			mul:  0.01,
			res:  amt.SetInt(0),
		},
		{
			name: "Round up",
			amt:  amt.SetInt(50), // 50 Satoshis
			mul:  0.01,
			res:  amt.SetInt(1), // 1 Satoshi
		},
		{
			name: "Multiply by 0.",
			amt:  amt.SetInt(1e8), // 1 BTC
			mul:  0,
			res:  amt.SetInt(0), // 0 BTC
		},
		{
			name: "Multiply 1 by 0.5.",
			amt:  amt.SetInt(1), // 1 Satoshi
			mul:  0.5,
			res:  amt.SetInt(1), // 1 Satoshi
		},
		{
			name: "Multiply 100 by 66%.",
			amt:  amt.SetInt(100), // 100 Satoshis
			mul:  0.66,
			res:  amt.SetInt(66), // 66 Satoshis
		},
		{
			name: "Multiply 100 by 66.6%.",
			amt:  amt.SetInt(100), // 100 Satoshis
			mul:  0.666,
			res:  amt.SetInt(67), // 67 Satoshis
		},
		{
			name: "Multiply 100 by 2/3.",
			amt:  amt.SetInt(100), // 100 Satoshis
			mul:  2.0 / 3,
			res:  amt.SetInt(67), // 67 Satoshis
		},
	}

	for _, test := range tests {
		a := test.amt.MulF64(test.mul) 
		if a.Cmp(test.res) != 0 {
			t.Errorf("%v: expected %v got %v", test.name, test.res, a)
		}
	}
}
