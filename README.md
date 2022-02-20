# Vuitton

Vuitton checks product availability via the [official Louis Vuitton website](https://en.louisvuitton.com).

Vuitton is a convenient, platform-independent, zero-installation, zero-configuration product monitor for
the command line. Precompiled binaries are provided for Mac OSX, Windows and Linux.

While not guaranteed to work with every product from the Louis Vuitton website, it has been tested with a variety
of products, including bags and shoes. For shoes, it can track stock availability for specific shoe sizes, as long
as the URL contains adequate information (more information below). 

_Note:_ This package is provided _as-is_. It's a personal project foremost, and although it works perfectly with the
current version of Louis Vuitton's API, it is not guaranteed to do so at all times. Use at your own risk.


## Features

Reads product URL's (and optionally, SKU's) from a regular text file (the "P-file" / "Product-file").
Keeps track of availability and notifies you when products come in stock. Details:

* Reads product URL's from a simple text file
* Extracts product ID's from URL's for you
* Supports checking for specific SKU's when requested
* Full support for different regions/countries
* Reloads the text file periodically
* Periodically checks product availability directly against the Louis Vuitton REST API
* Keeps track of state, so will only let you know when out-of-stock products comes in stock
* Supports desktop notifications
* Will open the product in your browser when it comes in stock

Currently re-checks the P-file every 10 seconds and checks product availability every 30 seconds.
See command line flags (`./vuitton -help`) for how to change these.


## Caveats

* Will only track/monitor up to ten products (to avoid hitting the rate limiter)
* Presently only works for URL's containing product ID's prefixed with "nvprod"
* Product availability is checked periodically but not aggressively due to the API utilizing a rate limiter


## Getting Started

1. Download
   1. If you have Go 1.17+ already, a simple `go get github.com/mkock/vuitton` will do 
   2. If not, then `git clone` should do it
2. Install
   1. Run `go build -o vuitton cmd/main.go` (assuming Linux - always build to your platform) 
   2. You can also run the binaries available in the `bins` directory
3. Use
   1. From the command line, first create an empty "products.txt" file; you can use another name but then you'll need to use a flag to tell the application which name you are using.
   2. Find the products you want to monitor, on the [Louis Vuitton website](https://en.louisvuitton.com/)
   3. If the URL contains the product code (prefixed with "nvprod"), it should be usable; Add the URL to the text file (one URL per line)
   4. Start the application: `./vuitton` (run with `-help` to see available flags)
   5. It will start monitoring products:
      1. The text file can be edited while the application is running, it will reload it periodically
      2. If a product comes in stock, a desktop notification should appear; the browser should also open the URL for you
      3. Stock level is maintained as long as the application is running, so you'll only be notified when a product _comes back in stock_
      4. To quit, press CTRL+C


## Product URL's

This is an example of an acceptable product URL:

`https://en.louisvuitton.com/eng-nl/products/pocket-organiser-damier-graphite-nvprod3430052v`

It contains the product ID `nvprod3430052v`.

Another example:

`https://en.louisvuitton.com/eng-nl/products/charlie-trainers-nvprod3130266v#1A9JN8`

This product (a pair of sneakers) has the product ID `nvprod3130266v` and the SKU `1A9JN8`.
The SKU uniquely identifies the shoe size (size 8), so availability will only be checked for that size if the SKU is
included in the URL. 


## Countries

Product availability varies by country.

The default country is Denmark (DK). You may pick another via the command line flags. Supported countries are:

BE, DE, DK, ES, FI, FR, IE, IT, LU, MC, NL, AT, SE, UK, RU, US, BR, CA, MX, CN, JP, KR, HK, SG, TW, TH, AU, NZ, UA, AE, SA, KW, KW, QA

Check the available shipping countries on the [Louis Vuitton website](https://en.louisvuitton.com/) for reference.


## Stock Keeping Units

A product ID does not always uniquely identify a product. For many bags and accessories, the product ID is sufficient.
But for shoes, t-shirts and other apparel where there are different sizes available, product availability varies by size.

If you know the SKU code, you can make a small change to the P-file. For example, you've identified the URL for a pair
of shoes and wish to check product availabilty. After selecting the desired shoe size on the website, you'll see the SKU
for that shoe size written right above the product title. This is a 6-character code containing both letters and digits.

Editing the P-file, simply add a hash symbol, `#`, followed by the SKU code. In some URL's, you might automatically get
this already when you copy/paste it.

During availability checks, the algorithm will look for this particular SKU and provide more accurate results.

## Notifications

Currently, the monitor will open a product URL in your default browser when it comes in stock, and send a desktop
notification. Either of these notification types can be disabled via the command-line flags.

## Intervals

When changing any of the intervals via the command line, you can use abbreviations such as "10s" (10 seconds),
"2m" (2 minutes) and so forth.

## Finally...

If you found this application useful, give me a star on GitHub to show your appreciation.
You can also give me a mention on [Twitter](https://twitter.com/MartinKock).

Enjoy!