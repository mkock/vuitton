package vuitton

import "strings"

// countryMap acts as an allow-list and a conversion chart.
var countryMap = map[string]string{
	"AE":  "ara-ae",
	"AT":  "eng-nl",
	"AU":  "eng-au",
	"BE":  "eng-nl",
	"BR":  "por-br",
	"CA":  "eng-ca",
	"CN":  "zhs-cn",
	"DE":  "deu-de",
	"DK":  "eng-nl",
	"ES":  "esp-es",
	"FI":  "eng-nl",
	"FR":  "fra-fr",
	"HK":  "eng-hk",
	"IE":  "eng-nl",
	"IT":  "ita-it",
	"JP":  "jpn-jp",
	"KR":  "kor-kr",
	"KW":  "eng-ae",
	"LU":  "eng-nl",
	"MC":  "eng-nl",
	"MX":  "esp-mx",
	"NL":  "eng-nl",
	"NZ":  "eng-sg",
	"QA":  "eng-ae",
	"RU":  "rus-ru",
	"SA":  "eng-ae",
	"SE":  "eng-nl",
	"SG":  "eng-sg",
	"TH":  "tha-th",
	"TW":  "zht-tw",
	"UAE": "eng-ae",
	"UK":  "eng-gb",
	"US":  "eng-us",
}

// Country represents a two-letter country code.
type Country string

// Valid returns true if the country code is on the allow-list defined by countryMap.
func (c Country) Valid() bool {
	if len(c) != 2 {
		return false
	}
	_, ok := countryMap[strings.ToUpper(string(c))]
	return ok
}

// Code returns the API country/language combo that corresponds to the country code.
func (c Country) Code() string {
	if code, ok := countryMap[strings.ToUpper(string(c))]; ok {
		return code
	}
	return ""
}
