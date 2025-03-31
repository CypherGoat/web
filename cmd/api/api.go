// Copyright (C) 2025 CypherGoat
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

var URL = getEnv("CYPHERGOAT_API_URL", "https://api.cyphergoat.com")

var API_KEY = getEnv("CYPHERGOAT_API_KEY", "")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Estimate struct {
	ExchangeName  string  `json:"Exchange"`
	ReceiveAmount float64 `json:"Amount"`
	MinAmount     float64 `json:"MinAmount"`
	Network1      string
	Network2      string
	Coin1         string
	Coin2         string
	SendAmount    float64
	Address       string
	ImageURL      string
	KYCScore      int `json:"KYCScore"`
}

type Info struct {
	IP        string
	UserAgent string
	LangList  string
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}
type Transaction struct {
	Coin1          string    `json:"Coin1,omitempty"`
	Coin2          string    `json:"Coin2,omitempty"`
	Network1       string    `json:"Network1,omitempty"`
	Network2       string    `json:"Network2,omitempty"`
	Address        string    `json:"Address,omitempty"`
	EstimateAmount float64   `json:"EstimateAmount,omitempty"`
	Provider       string    `json:"Provider,omitempty"`
	Id             string    `json:"Id,omitempty"`
	SendAmount     float64   `json:"SendAmount,omitempty"`
	Track          string    `json:"Track,omitempty"`
	Status         string    `json:"Status,omitempty"`
	KYC            string    `json:"KYC,omitempty"`
	Token          string    `json:"Token,omitempty"`
	Done           bool      `json:"Done,omitempty"`
	CGID           string    `json:"CGID,omitempty"`
	CreatedAt      time.Time `json:"CreatedAt,omitempty"`
	Affiliate      string    `json:"Affiliate,omitempty"`
}

func SendRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(data, &responseMap)
	if err != nil {
		return nil, err
	}

	if errStr, ok := responseMap["error"].(string); ok {
		return nil, fmt.Errorf("%s", errStr)
	}

	return data, nil
}

func FetchEstimateFromAPI(coin1, coin2 string, amount float64, best bool, network1, network2 string) ([]Estimate, error) {
	var url string
	if best {
		url = fmt.Sprintf("%s/estimate?coin1=%s&coin2=%s&amount=%f&network1=%s&network2=%s&best=true", URL, coin1, coin2, amount, network1, network2)
	} else {
		url = fmt.Sprintf("%s/estimate?coin1=%s&coin2=%s&amount=%f&network1=%s&network2=%s&best=false", URL, coin1, coin2, amount, network1, network2)
	}

	data, err := SendRequest(url)
	if err != nil {
		return nil, err
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(data, &responseMap)
	if err != nil {
		return nil, err
	}

	if errStr, ok := responseMap["error"].(string); ok {
		return nil, fmt.Errorf("%s", errStr)
	}

	type ApiResponse struct {
		Rates []Estimate `json:"rates"`
	}

	var result ApiResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	sort.Slice(result.Rates, func(i, j int) bool {
		return result.Rates[i].ReceiveAmount > result.Rates[j].ReceiveAmount
	})

	for i := range result.Rates {
		result.Rates[i].Coin1 = coin1
		result.Rates[i].Coin2 = coin2
		result.Rates[i].SendAmount = amount
		result.Rates[i].Network1 = network1
		result.Rates[i].Network2 = network2
	}

	return result.Rates, nil
}

func CreateTradeFromAPI(coin1, coin2 string, amount float64, address, partner string, network1, network2, affiliate string, info Info) (error, Transaction) {
	params := url.Values{}
	params.Add("coin1", coin1)
	params.Add("coin2", coin2)
	params.Add("amount", fmt.Sprintf("%f", amount))
	params.Add("partner", partner)
	params.Add("address", address)
	params.Add("network1", network1)
	params.Add("network2", network2)
	params.Add("affiliate", affiliate)
	params.Add("ip", info.IP)
	params.Add("useragent", info.UserAgent)
	params.Add("lang", info.LangList)

	requestURL := fmt.Sprintf("%s/swap?%s", URL, params.Encode())

	data, err := SendRequest(requestURL)
	if err != nil {
		return err, Transaction{}
	}

	var result TransactionResponse
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err, Transaction{}
	}

	transaction := result.Transaction
	fmt.Printf("Transaction: %+v\n", transaction)
	return nil, transaction
}

func TrackTxFromAPI(t Transaction) (error, Transaction) {
	url := fmt.Sprintf("%s/transaction?id=%s", URL, strings.ToLower(t.Provider), t.Id)
	data, err := SendRequest(url)
	if err != nil {
		return err, t
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(data, &responseMap)
	if err != nil {
		return err, t
	}

	status, ok := responseMap["status"].(string)
	if !ok {
		return fmt.Errorf("status field is missing or not a string"), Transaction{}

	}
	t.Status = status

	return nil, t

}

func GetTransactionFromAPI(id string) (error, Transaction) {
	url := fmt.Sprintf("%s/transaction?id=%s", URL, id)
	data, err := SendRequest(url)
	if err != nil {
		return err, Transaction{}
	}

	var result map[string]Transaction
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err, Transaction{}
	}

	transaction := result["transaction"]
	return nil, transaction
}
