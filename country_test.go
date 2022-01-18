package vuitton

import "testing"

func TestCountryValid(t *testing.T) {
	tests := []struct {
		in    string
		valid bool
	}{
		{"", false},
		{"hello", false},
		{"dk", true},
		{"DK", true},
		{"Dk", true},
		{"SG", true},
		{"xx", false},
	}

	var actual bool
	for _, tt := range tests {
		c := Country(tt.in)
		actual = c.Valid()
		if actual != tt.valid {
			t.Errorf("%s: expected %t, got %t", tt.in, tt.valid, actual)
		}
	}
}
