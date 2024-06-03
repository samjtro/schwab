package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	schwabutils "github.com/samjtro/go-trade/schwab/utils"
	utils "github.com/samjtro/go-trade/utils"
)

// SearchInstrumentSimple returns instrument's simples.
// It takes on param:
func SearchInstrumentSimple(cusip string) (SimpleInstrument, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(Endpoint_searchInstrument, cusip), nil)
	utils.Check(err)
	body, err := schwabutils.Handler(req)
	utils.Check(err)
	var instrument SimpleInstrument
	err = json.Unmarshal([]byte(body), &instrument)
	utils.Check(err)
	return instrument, nil
}

// SearchInstrumentFundamental returns instrument's fundamentals.
// It takes one param:
func SearchInstrumentFundamental(symbol string) (FundamentalInstrument, error) {
	req, err := http.NewRequest("GET", Endpoint_searchInstruments, nil)
	utils.Check(err)
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("projection", "fundamental")
	req.URL.RawQuery = q.Encode()
	body, err := schwabutils.Handler(req)
	utils.Check(err)
	var instrument FundamentalInstrument
	err = json.Unmarshal([]byte(body), &instrument)
	utils.Check(err)
	return instrument, nil
}
