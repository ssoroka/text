// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package currency

import (
	"testing"

	"golang.org/x/text/language"
)

func TestParseISO(t *testing.T) {
	testCases := []struct {
		in  string
		out string
		ok  bool
	}{
		{"USD", "USD", true},
		{"xxx", "XXX", true},
		{"xts", "XTS", true},
		{"XX", "XXX", false},
		{"XXXX", "XXX", false},
		{"", "XXX", false},       // not well-formed
		{"UUU", "XXX", false},    // unknown
		{"\u22A9", "XXX", false}, // non-ASCII, printable

		{"aaa", "XXX", false},
		{"zzz", "XXX", false},
		{"000", "XXX", false},
		{"999", "XXX", false},
		{"---", "XXX", false},
		{"\x00\x00\x00", "XXX", false},
		{"\xff\xff\xff", "XXX", false},
	}
	for i, tc := range testCases {
		if x, err := ParseISO(tc.in); x.String() != tc.out || err == nil != tc.ok {
			t.Errorf("%d:%s: was %s, %v; want %s, %v", i, tc.in, x, err == nil, tc.out, tc.ok)
		}
	}
}

func TestFromRegion(t *testing.T) {
	testCases := []struct {
		region, currency string
		ok               bool
	}{
		{"NL", "EUR", true},
		{"BE", "EUR", true},
		{"AG", "XCD", true},
		{"CH", "CHF", true},
		{"CU", "CUP", true},   // first of multiple
		{"DG", "USD", true},   // does not have M49 code
		{"150", "XXX", false}, // implicit false
		{"CP", "XXX", false},  // explicit false in CLDR
		{"CS", "XXX", false},  // all expired
		{"ZZ", "XXX", false},  // none match
	}
	for _, tc := range testCases {
		cur, ok := FromRegion(language.MustParseRegion(tc.region))
		if cur.String() != tc.currency || ok != tc.ok {
			t.Errorf("%s: got %v, %v; want %v, %v", tc.region, cur, ok, tc.currency, tc.ok)
		}
	}
}

func TestFromTag(t *testing.T) {
	testCases := []struct {
		tag, currency string
		conf          language.Confidence
	}{
		{"nl", "EUR", language.Low},      // nl also spoken outside Euro land.
		{"nl-BE", "EUR", language.Exact}, // region is known
		{"pt", "BRL", language.Low},
		{"en", "USD", language.Low},
		{"en-u-cu-eur", "EUR", language.Exact},
		{"tlh", "XXX", language.No}, // Klingon has no country.
		{"es-419", "XXX", language.No},
		{"und", "USD", language.Low},
	}
	for _, tc := range testCases {
		cur, conf := FromTag(language.MustParse(tc.tag))
		if cur.String() != tc.currency || conf != tc.conf {
			t.Errorf("%s: got %v, %v; want %v, %v", tc.tag, cur, conf, tc.currency, tc.conf)
		}
	}
}

var (
	czk = MustParseISO("CZK")
	zwr = MustParseISO("ZWR")
)

func TestTable(t *testing.T) {
	for i := 4; i < len(currency); i += 4 {
		if a, b := currency[i-4:i-1], currency[i:i+3]; a >= b {
			t.Errorf("currency unordered at element %d: %s >= %s", i, a, b)
		}
	}
	// First currency has index 1, last is numCurrencies.
	if c := currency.Elem(1)[:3]; c != "ADP" {
		t.Errorf("first was %c; want ADP", c)
	}
	if c := currency.Elem(numCurrencies)[:3]; c != "ZWR" {
		t.Errorf("last was %c; want ZWR", c)
	}
}

func TestKindRounding(t *testing.T) {
	testCases := []struct {
		kind  Kind
		cur   Currency
		scale int
		inc   int
	}{
		{Standard, USD, 2, 1},
		{Standard, CHF, 2, 1},
		{Cash, CHF, 2, 5},
		{Standard, TWD, 2, 1},
		{Cash, TWD, 0, 1},
		{Standard, czk, 2, 1},
		{Cash, czk, 0, 1},
		{Standard, zwr, 2, 1},
		{Cash, zwr, 0, 1},
	}
	for i, tc := range testCases {
		if scale, inc := tc.kind.Rounding(tc.cur); scale != tc.scale && inc != tc.inc {
			t.Errorf("%d: got %d, %d; want %d, %d", i, scale, inc, tc.scale, tc.inc)
		}
	}
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		USD.String()
	}
}
