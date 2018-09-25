// Copyright (c) 2013-2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58_test

import (
	"testing"

	"github.com/wakiyamap/monautil/base58"
)

var checkEncodingStringTests = []struct {
	version []byte
	in      string
	out     string
}{
	{[]byte{20}, "", "3MNQE1X"},
	{[]byte{20}, " ", "B2Kr6dBE"},
	{[]byte{20}, "-", "B3jv1Aft"},
	{[]byte{20}, "0", "B482yuaX"},
	{[]byte{20}, "1", "B4CmeGAC"},
	{[]byte{20}, "-1", "mM7eUf6kB"},
	{[]byte{20}, "11", "mP7BMTDVH"},
	{[]byte{20}, "abc", "4QiVtDjUdeq"},
	{[]byte{20}, "1234598760", "ZmNb8uQn5zvnUohNCEPP"},
	{[]byte{20}, "abcdefghijklmnopqrstuvwxyz", "K2RYDcKfupxwXdWhSAxQPCeiULntKm63UXyx5MvEH2"},
	{[]byte{20}, "00000000000000000000000000000000000000000000000000000000000000", "bi1EWXwJay2udZVxLJozuTb8Meg4W9c6xnmJaRDjg6pri5MBAxb9XwrpQXbtnqEoRV5U2pixnFfwyXC8tRAVC8XxnjK"},
}

func TestBase58Check(t *testing.T) {
	for x, test := range checkEncodingStringTests {
		// test encoding
		if res := base58.CheckEncode([]byte(test.in), test.version); res != test.out {
			t.Errorf("CheckEncode test #%d failed: got %s, want: %s", x, res, test.out)
		}

		// test decoding
		res, version, err := base58.CheckDecode(test.out, 1)
		if err != nil {
			t.Errorf("CheckDecode test #%d failed with err: %v", x, err)
		} else if !bytes.Equal(version, test.version) {
			t.Errorf("CheckDecode test #%d failed: got version: %d want: %d", x, version, test.version)
		} else if string(res) != test.in {
			t.Errorf("CheckDecode test #%d failed: got: %s want: %s", x, res, test.in)
		}
	}

	// test the two decoding failure cases
	// case 1: checksum error
	_, _, err := base58.CheckDecode("3MNQE1Y", 1)
	if err != base58.ErrChecksum {
		t.Error("Checkdecode test failed, expected ErrChecksum")
	}
	// case 2: invalid formats (string lengths below 5 mean the version byte and/or the checksum
	// bytes are missing).
	testString := ""
	for len := 0; len < 4; len++ {
		// make a string of length `len`
		_, _, err = base58.CheckDecode(testString, 1)
		if err != base58.ErrInvalidFormat {
			t.Error("Checkdecode test failed, expected ErrInvalidFormat")
		}
	}

}
