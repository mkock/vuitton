package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"vuitton"

	"github.com/atomicgo/cursor"
	"github.com/gen2brain/beeep"
)

const minIntervalsSeconds = 2

var (
	pFileName            string
	countryCode          string
	openBrowser          bool
	notify               bool
	pFileInterval        time.Duration
	availabilityInterval time.Duration
)

// init handles CLI flags.
func init() {
	flag.StringVar(&countryCode, "country", "dk", "country code to check availability for, two letters, any case")
	flag.StringVar(&pFileName, "filename", "products.txt", "name of file to load product URLs from")
	flag.BoolVar(&openBrowser, "browser", true, "attempt to open the product URL in your browser when it comes in stock")
	flag.BoolVar(&notify, "notify", true, "attempt to notify via desktop notification when product comes in stock")
	flag.DurationVar(&pFileInterval, "pfilecheck", 10*time.Second, "interval between reloads of 'p-file' (product-file)")
	flag.DurationVar(&availabilityInterval, "availabilitycheck", 30*time.Second, "interval between product availability checks")
	flag.Parse()

}

// desktopNotification attempts to push a desktop notification when a product comes in stock.
func desktopNotification(title, msg string) {
	if notify {
		_ = beeep.Notify(title, msg, "")
	}
}

// main runs the main loop, which has two purposes:
// 1) check for updates to the product file
// 2) check product availability
// CTRL+C will interrupt the loop.
func main() {
	// Resolve flags etc.
	countryCode := strings.ToUpper(countryCode)
	exitCode := 0

	// Validate country.
	country := vuitton.Country(countryCode)
	if !country.Valid() {
		msg := fmt.Sprintln("Invalid country. Acceptable values are: BE, DE, DK, ES, FI, FR, IE, IT, LU, MC, NL, AT, SE, UK, RU, US, BR, CA, MX, CN, JP, KR, HK, SG, TW, TH, AU, NZ, UA, AE, SA, KW, KW, QA")
		printErrorUsageAndExit(1, msg)
	}

	// Validate durations.
	if availabilityInterval.Seconds() < minIntervalsSeconds {
		msg := fmt.Sprintf("Invalid duration for availability interval, must be at least %d seconds\n", minIntervalsSeconds)
		printErrorUsageAndExit(2, msg)
	}
	if pFileInterval.Seconds() < minIntervalsSeconds {
		msg := fmt.Sprintf("Invalid duration for p-file interval, must be at least %d seconds\n", minIntervalsSeconds)
		printErrorUsageAndExit(3, msg)
	}

	// Check if p-file exists.
	info, err := os.Stat(pFileName)
	if err != nil {
		msg := "Unable to read p-file, please create one, make sure it's readable and then use the 'filename' argument to point to the file's location\n"
		printErrorUsageAndExit(4, msg)
	}
	if !info.Mode().IsRegular() {
		msg := "p-file does not look like it's a regular text file\n"
		printErrorUsageAndExit(5, msg)
	}

	view := cursor.NewArea()
	m := vuitton.MainLoop{
		ViewPort:             &view,
		Country:              country,
		AvailabilityInterval: availabilityInterval,
		RequestTimeout:       5 * time.Second,
		Client:               &http.Client{},
		PFileName:            pFileName,
		OpenBrowser:          openBrowser,
		Notification:         desktopNotification,
		PFileInterval:        pFileInterval,
	}
	err = m.Run()
	if err != nil {
		fmt.Println("Error:", err.Error())
		exitCode = 6
	}
	os.Exit(exitCode)
}

func printErrorUsageAndExit(exitCode int, msg string) {
	fmt.Println("Error:", msg)
	flag.Usage()
	os.Exit(exitCode)
}
