package vuitton

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// lvURL is the URL that provides product availability on Louis Vuitton's API. Needs a country code and a product ID.
const lvURL = "https://api.louisvuitton.com/api/%s/catalog/availability/%s"

// avail describes the product availability for a single product, identified by its SKUID.
type avail struct {
	SKUID   string `json:"skuId"`
	Exists  bool   `json:"exists"`
	InStock bool   `json:"inStock"`
}

// response describes the part of the API response that we wish to process.
type response struct {
	SKUAvailability []avail `json:"skuAvailability"`
}

// availability checks product availability for the given product ID.
func (m *MainLoop) availability(p product) (inStock bool, err error) {
	c := m.Client
	c.Timeout = m.RequestTimeout

	pID := p.productID()
	if pID == "" {
		return false, errors.New("invalid URL or no product ID")
	}

	url := fmt.Sprintf(lvURL, m.Country.Code(), pID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	setCommonHeaders(req)
	req.Header.Add("origin", p.Domain())
	req.Header.Add("referer", p.Domain())

	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}

	if resp.Body == nil {
		return false, errors.New("empty response")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("request unsuccessful, status code is %d", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var skus response
	err = json.Unmarshal(bytes, &skus)
	if err != nil {
		return false, err
	}

	mySKU := p.SKU()
	switch {
	case len(skus.SKUAvailability) == 0:
		return false, errors.New("no SKU's available")
	case len(skus.SKUAvailability) == 1 || mySKU == "":
		return skus.SKUAvailability[0].InStock, nil
	default:
		for _, sku := range skus.SKUAvailability {
			if sku.SKUID == mySKU {
				return sku.InStock, nil
			}
		}
		return false, nil
	}
}

// setCommonHeaders adds common headers to the availability request.
func setCommonHeaders(req *http.Request) {
	req.Header.Add("authority", "api.louisvuitton.com")
	req.Header.Add("sec-ch-ua", `"Chromium";v="96", "Opera";v="82", ";Not A Brand";v="99"`)
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("dnt", `1`)
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36 OPR/82.0.4227.23")
	req.Header.Add("sec-ch-ua-platform", "Linux")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("accept-language", "en-US,en;q=0.9,da;q=0.8")
}
