package vuitton

import (
	"regexp"
	"strings"
)

var productRegExp = regexp.MustCompile(`nvprod[0-9a-z]*`)

// product represents a single Louis Vuitton product.
type product struct {
	URL string
}

// Valid returns true if the product URL looks valid, ie. points to louisvuitton.com and contains a product ID.
func (p product) Valid() bool {
	if p.URL == "" ||
		!strings.HasPrefix(p.URL, "https://") ||
		!strings.Contains(p.URL, "louisvuitton") ||
		!strings.Contains(p.URL, "nvprod") {
		return false
	}
	return true
}

// Domain returns the base domain of the product's URL.
// If the product URL is invalid, an empty string is returned.
func (p product) Domain() string {
	if !p.Valid() {
		return ""
	}
	parts := strings.SplitAfter(p.URL, ".com")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// productID returns the product ID for the product identified by its URL.
// productID does this be extracting the product code (usually prefixed with "nvprod") from the URL.
// If the product URL is invalid, or doesn't contain a product ID, an empty string is returned.
func (p product) productID() string {
	// Example URL: https://en.louisvuitton.com/eng-nl/products/ecorce-rousse-perfumed-candle-nvprod1910068v.
	// Result: nvprod1910068v.
	if !p.Valid() {
		return ""
	}
	return productRegExp.FindString(p.URL)
}

// SKU returns the SKU for the product identified by its URL.
// SKU does this by extracting the SKU (usually prefixed with "#") from the URL.
// If the product URL is invalid, or doesn't contain an SKU, an empty string is returned.
func (p product) SKU() string {
	// Example URL: https://en.louisvuitton.com/eng-nl/products/charlie-trainers-nvprod3130266v#1A9JN8.
	// Result: 1A9JN8.
	if !p.Valid() || !strings.ContainsRune(p.URL, '#') {
		return ""
	}
	parts := strings.Split(p.URL, "#")
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
