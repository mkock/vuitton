package vuitton

import (
	"testing"
)

func TestValid(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{"", false},
		{"hello", false},
		{"https://www.google.com", false},
		{"http://en.louisvuitton.com/eng-nl/products/", false},
		{"https://en.louisvuitton.com/eng-nl/products/", false},
		{"https://us.louisvuitton.com/eng-nl/products/loop-bag-monogram-nvprod3190103v?country=DK", true},
		{"https://jp.louisvuitton.com/eng-nl/products/trio-messenger-bag-damier-graphite-nvprod3430073v", true},
		{"https://en.louisvuitton.com/eng-nl/products/croisillon-shawl-nvprod3390166v#M77459", true},
		{"https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v", true},
	}

	var valid bool
	for _, tt := range tests {
		p := product{URL: tt.in}
		valid = p.Valid()
		if valid != tt.valid {
			t.Errorf("expected %t, got %t", tt.valid, valid)
		}
	}
}

func TestDomain(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"hello", ""},
		{"https://en.louisvuitton.com/eng-nl/products/", ""},
		{"https://us.louisvuitton.com/eng-nl/products/loop-bag-monogram-nvprod3190103v?country=DK", "https://us.louisvuitton.com"},
		{"https://jp.louisvuitton.com/eng-nl/products/trio-messenger-bag-damier-graphite-nvprod3430073v", "https://jp.louisvuitton.com"},
		{"https://en.louisvuitton.com/eng-nl/products/croisillon-shawl-nvprod3390166v#M77459", "https://en.louisvuitton.com"},
		{"https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v", "https://en.louisvuitton.com"},
	}

	var actual string
	for _, tt := range tests {
		p := product{URL: tt.in}
		actual = p.Domain()
		if actual != tt.out {
			t.Errorf("expected %q, got %q", tt.out, actual)
		}
	}
}

func TestProductID(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"hello", ""},
		{"https://en.louisvuitton.com/eng-nl/products/", ""},
		{"https://us.louisvuitton.com/eng-nl/products/loop-bag-monogram-nvprod3190103v?country=DK", "nvprod3190103v"},
		{"https://jp.louisvuitton.com/eng-nl/products/trio-messenger-bag-damier-graphite-nvprod3430073v", "nvprod3430073v"},
		{"https://en.louisvuitton.com/eng-nl/products/croisillon-shawl-nvprod3390166v#M77459", "nvprod3390166v"},
		{"https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v", "nvprod1910068v"},
	}

	var actual string
	for _, tt := range tests {
		p := product{URL: tt.in}
		actual = p.productID()
		if actual != tt.out {
			t.Errorf("expected %q, got %q", tt.out, actual)
		}
	}
}

func TestSKU(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"hello", ""},
		{"https://en.louisvuitton.com/eng-nl/products/", ""},
		{"https://us.louisvuitton.com/eng-nl/products/loop-bag-monogram-nvprod3190103v?country=DK", ""},
		{"https://jp.louisvuitton.com/eng-nl/products/charlie-trainers-nvprod3130266v#1A9JN8", "1A9JN8"},
		{"https://en.louisvuitton.com/eng-nl/products/charlie-trainers-nvprod3130266v#1A9JNC", "1A9JNC"},
		{"https://en.louisvuitton.com/eng-nl/products/charlie-trainers-nvprod3130266v#	1A9JNC ", "1A9JNC"},
		{"https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v#", ""},
		{"https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v# ", ""},
	}

	var actual string
	for _, tt := range tests {
		p := product{URL: tt.in}
		actual = p.SKU()
		if actual != tt.out {
			t.Errorf("expected %q, got %q", tt.out, actual)
		}
	}
}
