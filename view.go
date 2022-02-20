package vuitton

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	title      = "⥮ Louis Vuitton Product Monitor"
	pIDPadding = 20
)

// output replaces the product table with the provided message.
func (m *MainLoop) output(msg string) {
	if msg == "" {
		return
	}
	m.ViewPort.Clear()
	m.ViewPort.Update(title + "\n\n" + msg + "\n")
}

// updateView renders a table with the current product list and stock levels.
// If a message has been set on MainLoop.message, it will be rendered below the table.
func (m *MainLoop) updateView() {
	m.Lock()
	defer m.Unlock()

	b := strings.Builder{}

	inStock := func(val bool) string {
		if val {
			return "Yes"
		}
		return "No"
	}

	// Render header info.
	b.WriteString(title + "\n\n")
	h := tablewriter.NewWriter(&b)
	h.SetBorder(false)
	h.SetAlignment(tablewriter.ALIGN_LEFT)
	h.Append([]string{"Region", m.Country.Code()})
	h.Append([]string{"Product file", m.PFileName})
	h.Append([]string{"Products found", strconv.Itoa(len(m.products))})
	h.Append([]string{"Check interval", m.AvailabilityInterval.String()})
	h.Render()
	b.WriteString("\n")

	// Render stock level table.
	t := tablewriter.NewWriter(&b)
	t.SetHeader([]string{"Product", "SKU", "In stock?"})
	pIDs := make([]string, 0, len(m.products))
	for pID := range m.products {
		pIDs = append(pIDs, pID)
	}
	sort.Strings(pIDs)
	for _, pID := range pIDs {
		stockLevel := m.products[pID]
		if !stockLevel.product.Valid() {
			pID = "Invalid product URL!"
		}
		t.Append([]string{fmt.Sprintf("%-*s", pIDPadding, pID), stockLevel.product.SKU(), inStock(stockLevel.inStock)})
	}
	t.Render()

	// Render message.
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.message + "\n")
	}

	m.ViewPort.Clear()
	m.ViewPort.Update(b.String())
}
