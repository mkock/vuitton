package vuitton

import (
	"net/url"
	"regexp"
	"strings"
)

var productRegExp = regexp.MustCompile(`nvprod[0-9a-z]*`)

// product represents a single Louis Vuitton product.
type product struct {
	URL string
}

// Valid returns true if the product URL looks valid, ie. points to louisvuitton.com and looks like a product URL.
func (p product) Valid() bool {
	if p.URL == "" ||
		!strings.HasPrefix(p.URL, "https://") ||
		!strings.Contains(p.URL, "louisvuitton") ||
		!strings.Contains(p.URL, "products") {
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
	// This must also work: https://en.louisvuitton.com/eng-nl/products/pochette-accessoires-monogram-005656.
	// Result: 005656.
	if !p.Valid() {
		return ""
	}
	match := productRegExp.FindString(p.URL)
	if match != "" {
		return match
	}
	// No nvprod ID, so let's look for a number.
	u, err := url.Parse(p.URL)
	if err == nil {
		parts := strings.Split(u.Path, "/")
		for i := range parts {
			if parts[i] == "products" && i < len(parts) && strings.Contains(parts[i+1], "-") {
				fields := strings.Split(parts[i+1], "-")
				return fields[len(fields)-1]
			}
		}
	}
	return ""
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
