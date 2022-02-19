package vuitton

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// stockLevel keeps track of stock levels across reloads.
type stockLevel struct {
	product   product
	inStock   bool
	updatedAt time.Time
}

// MainLoop is the loop that has a dual purpose:
// 1. Reload the products_sample.txt file when it changes
// 2. Periodically check product availability
type MainLoop struct {
	Country              Country
	AvailabilityInterval time.Duration
	RequestTimeout       time.Duration
	Client               *http.Client
	PFileName            string
	PFileInterval        time.Duration
	Notification         func(title, msg string)
	OpenBrowser          bool

	sync.Mutex // Protects the field(s) below.
	products   map[string]stockLevel
}

// Run starts the main loop. Run does not exit until interrupted by a SIGINT, or if an unrecoverable error occurs.
func (m *MainLoop) Run() error {
	// Init stock levels.
	m.products = make(map[string]stockLevel)

	// Set up OS interrupts.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up P-file monitor.
	lastRead := time.Time{}
	pFileTicker := time.NewTicker(m.PFileInterval)
	pFileFunc := func() {
		if !m.PFileModifiedSince(lastRead) {
			return
		}
		lastRead = time.Now()
		ps, err := m.ReadPFile()
		if err != nil {
			output(err.Error())
			return
		}
		if len(ps) == 0 {
			output("No products to monitor, please update your products text file")
			return
		}
		m.Lock()
		{
			// Update products.
			var pID string
			var pCached stockLevel
			var ok bool
			for _, p := range ps {
				pID = p.productID()
				if pID == "" {
					output("Failed to determine product ID for one of the URL's, does it include a product code?")
					continue
				}
				pCached, ok = m.products[pID]
				if ok {
					pCached.updatedAt = lastRead
				} else {
					m.products[pID] = stockLevel{
						product:   p,
						inStock:   false,
						updatedAt: lastRead,
					}
				}
			}
			// Un-cache removed products (they will have updatedAt < lastRead).
			for pID, pCached = range m.products {
				if pCached.updatedAt.Before(lastRead) {
					delete(m.products, pID)
				}
			}
		}
		m.Unlock()
		output(fmt.Sprintf("Reloaded %d products", len(ps)))
	}

	// Set up availability checks.
	availTicker := time.NewTicker(m.AvailabilityInterval)
	availFunc := func() {
		m.Lock()
		defer m.Unlock()
		if len(m.products) == 0 {
			return
		}
		// Perform availability checks.
		// TODO: We keep the lock during availability checks, this can be improved.
		output("Checking product availability")
		for pID, lvl := range m.products {
			inStock, err := m.availability(lvl.product)
			if err != nil {
				output(err.Error())
				continue
			}
			if inStock && !lvl.inStock {
				output(fmt.Sprintf("Product %q seems to be in stock!", pID))
				if m.Notification != nil {
					m.Notification("Vuitton Monitor", fmt.Sprintf("Product %q is in stock!", pID))
				}
				if m.OpenBrowser {
					m.browseTo(lvl.product.URL)
				}
			}
			lvl.inStock = inStock
			m.products[pID] = lvl
		}
	}

	// Intro.
	output(fmt.Sprintf("Region: %s", m.Country.Code()))
	output(fmt.Sprintf("Products file: %s", m.PFileName))

	// Load products once before entering the loop.
	pFileFunc()

	// Main loop.
	for {
		select {
		case <-sigs:
			fmt.Println(" Bye")
			return nil // TODO(mkock) Proper shutdown!
		case <-availTicker.C:
			go availFunc()
		case <-pFileTicker.C:
			go pFileFunc()
		}
	}
}

// browseTo opens a browser with the given url.
// If opening of the browser fails, then calling browseTo is a no-op.
func (m *MainLoop) browseTo(url string) {
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	default:
	}
}
