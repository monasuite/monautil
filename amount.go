// Copyright (c) 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package monautil

import (
	"errors"
	//"math"
	"math/big"
	"strconv"
	//"fmt"

	"github.com/shopspring/decimal"
)

// AmountUnit describes a method of converting an Amount to something
// other than the base unit of a monacoin.  The value of the AmountUnit
// is the exponent component of the decadic multiple to convert from
// an amount in monacoin to an amount counted in units.
type AmountUnit int

// These constants define various units used when describing a monacoin
// monetary amount.
const (
	AmountMegaBTC  AmountUnit = 6
	AmountKiloBTC  AmountUnit = 3
	AmountBTC      AmountUnit = 0
	AmountMilliBTC AmountUnit = -3
	AmountMicroBTC AmountUnit = -6
	AmountSatoshi  AmountUnit = -8
)

// msat max
// int64 is overflow \(^o^)/
const (
	MaxMsat string = "105120000000000000000"
	MinMsat string = "-105120000000000000000"
)

// String returns the unit as a string.  For recognized units, the SI
// prefix is used, or "Satoshi" for the base unit.  For all unrecognized
// units, "1eN BTC" is returned, where N is the AmountUnit.
func (u AmountUnit) String() string {
	switch u {
	case AmountMegaBTC:
		return "MMONA"
	case AmountKiloBTC:
		return "kMONA"
	case AmountBTC:
		return "MONA"
	case AmountMilliBTC:
		return "mMONA"
	case AmountMicroBTC:
		return "μMONA"
	case AmountSatoshi:
		return "Watanabe"
	default:
		return "1e" + strconv.FormatInt(int64(u), 10) + " MONA"
	}
}

// Amount represents the base monacoin monetary unit (colloquially referred
// to as a `Satoshi').  A single Amount is equal to 1e-8 of a monacoin.
type Amount big.Int

func (a *Amount) Add(amt1 Amount, amt2 Amount) Amount {
	return Amount(*a.Int().Add(amt1.Int(),amt2.Int()))
}

func (a *Amount) Sub(amt1 Amount, amt2 Amount) Amount {
	return Amount(*a.Int().Sub(amt1.Int(),amt2.Int()))
}

func (a *Amount) Mul(amt1 Amount, amt2 Amount) Amount {
	return Amount(*a.Int().Mul(amt1.Int(),amt2.Int()))
}

func (a *Amount) Cmp(amt2 Amount) int {
	return a.Int().Cmp(amt2.Int())
}

func (a *Amount) SetFloat(f float64) Amount {
	b := decimal.NewFromFloat(f)
	b.RoundBank(0)
	return Amount(*b.BigInt())
}

func (a *Amount) SetInt(i int64) Amount {
	return Amount(*big.NewInt(i))
}

func (a *Amount) SetBigFloat(bF *big.Float) Amount {
	b, _ := bF.Float64()
	c := decimal.NewFromFloat(b)
	c.RoundBank(0)
	return Amount(*c.BigInt())
}

func (a *Amount) SetDecimal(d decimal.Decimal) Amount {
	return Amount(*d.BigInt())
}

func (a *Amount) SetString(st string) (Amount, error) {
	d, err := decimal.NewFromString(st)
	if err != nil {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}
	// for mwat
	maxMsat, _ := decimal.NewFromString(MaxMsat)
	minMsat, _ := decimal.NewFromString(MinMsat)

	if d.Cmp(maxMsat) == 1 {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}
	if d.Cmp(minMsat) == -1 {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}

	d.RoundBank(0)
	return a.SetDecimal(d), nil
}

// round converts a floating point number, which may or may not be representable
// as an integer, to the Amount integer type by rounding to the nearest integer.
// This is performed by adding or subtracting 0.5 depending on the sign, and
// relying on integer truncation to round the value to the nearest Amount.
func round(bF *big.Float) Amount {
	if bF.Cmp(big.NewFloat(0)) == -1 {
		bF.Sub(bF, big.NewFloat(0.5))
	} else {
		bF.Add(bF, big.NewFloat(0.5))
	}
	var amt Amount
	return amt.SetBigFloat(bF)
}
// NewAmount creates an Amount from a floating point value representing
// some value in monacoin.  NewAmount errors if f is NaN or +-Infinity, but
// does not check that the amount is within the total amount of monacoin
// producible as f may not refer to an amount at a single moment in time.
//
// NewAmount is for specifically for converting BTC to Satoshi.
// For creating a new Amount with an int64 value which denotes a quantity of Satoshi,
// do a simple type conversion from type int64 to Amount.
// See GoDoc for example: http://godoc.org/github.com/monasuite/monautil#example-Amount
func NewAmount(st string) (Amount, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is +-Infinity.
	d, err := decimal.NewFromString(st)
	if err != nil {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}
	// for mwat(Inf+-)
	if d.Cmp(decimal.NewFromInt(105.12e6 * 1e3)) == 1 {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}
	if d.Cmp(decimal.NewFromInt(-105.12e6 * 1e3)) == -1 {
		return Amount(*big.NewInt(0)), errors.New("invalid monacoin amount")
	}

	decimalSatoshi := decimal.NewFromInt(SatoshiPerBitcoin)
	d2 := d.Mul(decimalSatoshi)
	s := d2.StringFixedBank(0)
	d3, _ := decimal.NewFromString(s)
	var amt Amount
	return amt.SetDecimal(d3), nil
}

// return big.Int from Amount
func (a *Amount) Int() *big.Int {
	return (*big.Int)(a)
}

// return decimal from Amount. when big.float->decimal converted, it is error......
func (a *Amount) Decimal() decimal.Decimal {
	return decimal.NewFromBigInt(a.Int(), 0)
}

// ToUnit converts a monetary amount counted in monacoin base units to a
// decimal value representing an amount of monacoin.
func (a *Amount) ToDecimalUnit(u AmountUnit) decimal.Decimal {
	b := decimal.NewFromBigInt(a.Int(),int32(-(u+8)))
	return b
}

// ToUnit converts a monetary amount counted in monacoin base units to a
// string value representing an amount of monacoin.
func (a *Amount) ToUnit(u AmountUnit) string {
	return a.ToDecimalUnit(u).String()
}

// ToBTC is the equivalent of calling ToUnit with AmountBTC.
func (a *Amount) ToBTC() string {
	return a.ToUnit(AmountBTC)
}

// Format formats a monetary amount counted in monacoin base units as a
// string for a given unit.  The conversion will succeed for any unit,
// however, known units will be formated with an appended label describing
// the units with SI notation, or "Satoshi" for the base unit.
func (a *Amount) Format(u AmountUnit) string {
	units := " " + u.String()
	return a.ToUnit(u) + units
}

// String is the equivalent of calling Format with AmountBTC.
func (a *Amount) String() string {
	return a.Format(AmountBTC)
}

// MulF64 multiplies an Amount by a floating point value.  While this is not
// an operation that must typically be done by a full node or wallet, it is
// useful for services that build on top of monacoin (for example, calculating
// a fee by multiplying by a percentage).
func (a *Amount) MulF64(f float64) Amount {
	x := big.NewFloat(f).SetPrec(64)
	y := big.NewFloat(1)
	y.SetInt(a.Int())
	z := big.NewFloat(1)
	z.Mul(x,y)
	return round(z)
}
