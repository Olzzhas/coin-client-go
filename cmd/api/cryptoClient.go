package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiUrl = "https://api.coingecko.com/api/v3/coins/"

type CoinData struct {
	Id         string `json:"id"`
	Symbol     string `json:"symbol"`
	MarketData struct {
		CurrentPrice map[string]float64 `json:"current_price"`
	} `json:"market_data"`
}

type CoinResponse struct {
	Coin     string  `json:"coin"`
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

func (app *application) getCoinData(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	currency := params.Get("currency")
	coin := params.Get("coin")

	url := fmt.Sprintf("%s/%s", apiUrl, coin)

	response, err := http.Get(url)
	if err != nil {
		_ = fmt.Errorf("")
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var coinData CoinData
	err = decoder.Decode(&coinData)
	if err != nil {
		_ = fmt.Errorf("")
	}

	var coinPrice float64

	if value, ok := coinData.MarketData.CurrentPrice[currency]; ok {
		coinPrice = value
	}

	responseData := new(CoinResponse)
	responseData.Coin = coin
	responseData.Currency = currency
	responseData.Value = coinPrice

	err = app.writeJSON(w, http.StatusOK, responseData, nil)
	if err != nil {
		_ = fmt.Errorf("error is occured while writing json: %v", err)
	}
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
