// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package monautil

const (
	// SatoshiPerBitcent is the number of satoshi in one monacoin cent.
	SatoshiPerBitcent = 1e6

	// SatoshiPerBitcoin is the number of satoshi in one monacoin (1 BTC).
	SatoshiPerBitcoin = 1e8

	// MaxSatoshi is the maximum transaction amount allowed in satoshi.
	MaxSatoshi = 10512e4 * SatoshiPerBitcoin
)
