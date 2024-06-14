package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Data map[string]float64 `json:"data"`
}

func main() {
	baseURL := "https://api.freecurrencyapi.com/v1/latest"
	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		fmt.Fprintln(os.Stderr, "Error. No API key found in environment valiables. Set it in the API_KEY variable.")
		os.Exit(1)
	}
	baseCurrency := os.Args[2]
	targetCurrencies := os.Args[1]

	apiKeyParam := "apikey=" + apiKey
	baseCurrencyParam := "base_currency=" + baseCurrency
	currenciesParam := "currencies=" + targetCurrencies
	url := baseURL + "?" + strings.Join([]string{apiKeyParam, baseCurrencyParam, currenciesParam}, "&")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if status := resp.StatusCode; status != 200 {
		fmt.Fprintf(os.Stderr, "HTTP response code: %d\n", status)
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", url, err)
		os.Exit(1)
	}

	currencyList := strings.Split(targetCurrencies, ",")
	for _, currency := range currencyList {
		var response Response
		if err = json.Unmarshal(data, &response); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshalling %s: %v\n", url, err)
			os.Exit(1)
		}

		value, _ := response.Data[currency]

		fmt.Printf("Value %s: %f %s\n", currency, value, baseCurrency)
	}
}
