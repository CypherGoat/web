// Copyright (C) 2025 CypherGoat

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

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

type Estimates struct {
	Results         []Estimate
	Min             float64
	TradeValue_fiat float64
	TradeValue_btc  float64
	CGSinStable     bool `json:"CGSinStable"`
}

type Estimate struct {
	ExchangeName    string  `json:"Exchange"`
	ReceiveAmount   float64 `json:"Amount"`
	MinAmount       float64 `json:"MinAmount"`
	Network1        string
	Network2        string
	Coin1           string
	Coin2           string
	SendAmount      float64
	Address         string
	ImageURL        string
	NoTextURL       string
	KYCScore        int `json:"KYCScore"`
	Log             bool
	CGShield        bool `json:"CGShield,omitempty"`
	CoveragePercent float64
	CGSAmount       string
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
	Memo           string    `json:"Memo,omitempty"`
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

func FetchEstimateFromAPI(coin1, coin2 string, amount float64, best bool, network1, network2 string) (Estimates, error) {
	var url string
	if best {
		url = fmt.Sprintf("%s/estimate?coin1=%s&coin2=%s&amount=%f&network1=%s&network2=%s&best=true", URL, coin1, coin2, amount, network1, network2)
	} else {
		url = fmt.Sprintf("%s/estimate?coin1=%s&coin2=%s&amount=%f&network1=%s&network2=%s&best=false", URL, coin1, coin2, amount, network1, network2)
	}

	data, err := SendRequest(url)
	if err != nil {
		return Estimates{}, err
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(data, &responseMap)
	if err != nil {
		return Estimates{}, err
	}

	if errStr, ok := responseMap["error"].(string); ok {
		return Estimates{}, fmt.Errorf("%s", errStr)
	}

	type ApiResponse struct {
		Rates Estimates `json:"rates"`
	}

	var result ApiResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		return Estimates{}, err
	}

	sort.Slice(result.Rates.Results, func(i, j int) bool {
		return result.Rates.Results[i].ReceiveAmount > result.Rates.Results[j].ReceiveAmount
	})

	for i := range result.Rates.Results {
		result.Rates.Results[i].Coin1 = coin1
		result.Rates.Results[i].Coin2 = coin2
		result.Rates.Results[i].SendAmount = amount
		result.Rates.Results[i].Network1 = network1
		result.Rates.Results[i].Network2 = network2
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

	params.Add("source", "clearnet-main")

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
	// fmt.Printf("Transaction: %+v\n", transaction)
	return nil, transaction
}

func TrackTxFromAPI(t Transaction) (error, Transaction) {
	url := fmt.Sprintf("%s/transaction?provider=%s&id=%s", URL, strings.ToLower(t.Provider), t.Id)
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

type AffiliateApplication struct {
	Email     string `json:"email"`
	Messenger string `json:"messenger"`
	Promotion string `json:"promotion"`
}

func SendAffiliateApplication(app AffiliateApplication) error {
	apiURL := fmt.Sprintf("%s/affiliate/apply", URL)
	payload, err := json.Marshal(app)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	return nil
}

type AffiliateLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AffiliateLoginResponse struct {
	Token string `json:"token"`
}

type AffiliateRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type AffiliateStats struct {
	TotalClicks     int              `json:"total_referrals"`
	TotalSwaps      int              `json:"total_completed"`
	TotalCommission float64          `json:"total_earnings"`
	PendingPayout   float64          `json:"balance"`
	AffiliateID     string           `json:"affiliate_id"`
	ReferralLink    string           `json:"referral_code"`
	Commission      float64          `json:"commission"`
	RecentSwaps     []RecentSwap     `json:"recent_swaps"`
	AllTrades       []AffiliateTrade `json:"all_trades"`
}

type AffiliateTrade struct {
	ID         string  `json:"id"`
	Date       string  `json:"date"`
	Coin1      string  `json:"coin1"`
	Coin2      string  `json:"coin2"`
	Amount     float64 `json:"amount"`
	Exchange   string  `json:"exchange"`
	Status     string  `json:"status"`
	Complete   bool    `json:"complete"`
	Commission float64 `json:"commission"`
}

type RecentSwap struct {
	Date       string  `json:"date"`
	Amount     float64 `json:"amount"`
	Coin1      string  `json:"coin1"`
	Coin2      string  `json:"coin2"`
	Commission float64 `json:"commission"`
	Exchange   string  `json:"exchange"`
}

type PayoutRequest struct {
	Amount         float64 `json:"amount"`
	PayoutCurrency string  `json:"payout_currency"` // "USDC_SOL" or "XMR"
	WalletAddress  string  `json:"wallet_address"`
}

type Payout struct {
	ID             int        `json:"id"`
	AffiliateID    int        `json:"affiliate_id"`
	Amount         float64    `json:"amount"`
	PayoutCurrency string     `json:"payout_currency"`
	WalletAddress  string     `json:"wallet_address"`
	Status         string     `json:"status"`
	AdminID        string     `json:"admin_id,omitempty"`
	AdminNotes     string     `json:"admin_notes,omitempty"`
	TxHash         string     `json:"tx_hash,omitempty"`
	RequestedAt    time.Time  `json:"requested_at"`
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
}

type PayoutHistory struct {
	PendingPayouts   []Payout `json:"pending_payouts"`
	CompletedPayouts []Payout `json:"completed_payouts"`
	TotalPaidOut     float64  `json:"total_paid_out"`
	PendingAmount    float64  `json:"pending_amount"`
}

type PayoutAction struct {
	PayoutID int    `json:"payout_id"`
	Action   string `json:"action"`
	AdminID  string `json:"admin_id"`
	TxHash   string `json:"tx_hash,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

func LoginAffiliate(username, password string) (string, error) {
	apiURL := fmt.Sprintf("%s/affiliate/login", URL)

	loginReq := AffiliateLoginRequest{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(loginReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("login failed")
	}

	var loginResp AffiliateLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	return loginResp.Token, nil
}

func RegisterAffiliate(username, email, password, token string) error {
	apiURL := fmt.Sprintf("%s/affiliate/register", URL)

	registerReq := AffiliateRegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Token:    token,
	}

	payload, err := json.Marshal(registerReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed: %s", string(body))
	}

	return nil
}

func GetAffiliateStats(token string) (*AffiliateStats, error) {
	apiURL := fmt.Sprintf("%s/affiliate/stats", URL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("X-API-Key", API_KEY)
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get stats (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stats AffiliateStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	return &stats, nil
}

func RequestPayout(token string, request PayoutRequest) (*Payout, error) {
	apiURL := fmt.Sprintf("%s/affiliate/payout/request", URL)

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("payout request failed: %s", string(body))
	}

	var response struct {
		Message string `json:"message"`
		Payout  Payout `json:"payout"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Payout, nil
}

func GetPayoutHistory(token string) (*PayoutHistory, error) {
	apiURL := fmt.Sprintf("%s/affiliate/payout/history", URL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get payout history: %s", string(body))
	}

	var history PayoutHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}

	return &history, nil
}
