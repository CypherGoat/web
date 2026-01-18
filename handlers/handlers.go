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

package handlers

import (
	"encoding/json"
	"fmt"
	"html"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CypherGoat/web/views"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"

	"github.com/CypherGoat/web/cmd/api"

	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

var (
	transactions = make(map[string]api.Transaction)
	mu           sync.Mutex
)

func IndexHandler(c echo.Context) error {
	return views.IndexTempl().Render(c.Request().Context(), c.Response())
}

func CGShieldTermsHandler(c echo.Context) error {
	return views.CGShield().Render(c.Request().Context(), c.Response())
}

func ContactHandler(c echo.Context) error {
	return views.Contact().Render(c.Request().Context(), c.Response())
}

func AffiliateTerms(c echo.Context) error {
	return views.AffiliateTerms().Render(c.Request().Context(), c.Response())
}

func AboutHandler(c echo.Context) error {
	return views.About().Render(c.Request().Context(), c.Response())

}

func TermsHandler(c echo.Context) error {
	return views.Terms().Render(c.Request().Context(), c.Response())

}

func PrivacyPolicyHandler(c echo.Context) error {
	return views.Privacy().Render(c.Request().Context(), c.Response())

}

func AffiliateHandler(c echo.Context) error {
	return views.Affiliate().Render(c.Request().Context(), c.Response())
}

func NotFoundHandler(c echo.Context) error {
	return views.NotFoundPage().Render(c.Request().Context(), c.Response())
}

func ExchangesHandler(c echo.Context) error {
	return views.Exchanges().Render(c.Request().Context(), c.Response())
}

var blockedExchangesForAnonymousNetworks = map[string]bool{
	"stealthex":  true,
	"changenow":  true,
	"fixedfloat": true,
	"simpleswap": true,
	"exolix":     true,
	"nanswap":    true,
}

func isExchangeBlocked(exchangeName string, isAnonymousNetwork bool) bool {
	if !isAnonymousNetwork {
		return false
	}
	return blockedExchangesForAnonymousNetworks[strings.ToLower(exchangeName)]
}

func ExchangeDetailHandler(c echo.Context) error {
	shortCode := c.Param("shortcode")

	var exchange views.Exchange
	found := false
	for _, ex := range views.ExchangesList {
		if ex.ShortCode == shortCode {
			exchange = ex
			found = true
			break
		}
	}

	if !found {
		return c.Redirect(http.StatusSeeOther, "/exchanges")
	}

	return views.ExchangeDetail(exchange).Render(c.Request().Context(), c.Response())
}

type ExchangeInfo struct {
	ImageURL  string
	NoTextURL string
	RequireIP bool
}

var exchangeInfo = map[string]ExchangeInfo{
	"alfacash":     {"/exchanges/alfacash.svg", "/exchanges/no-text/alfacash.jpg", true},
	"changee":      {"/exchanges/changee.svg", "", true},
	"changehero":   {"/exchanges/changehero.svg", "/exchanges/no-text/changehero.jpeg", true},
	"changenow":    {"/exchanges/changenow.svg", "", true},
	"coincraddle":  {"/exchanges/coincraddle.png", "", true},
	"exch.cx":      {"/exchanges/exchcx.png", "", false},
	"fixedfloat":   {"/exchanges/fixedfloat-v2.svg", "/exchanges/no-text/fixedfloat.svg", true},
	"majesticbank": {"/exchanges/majesticbank.png", "", false},
	"nanswap":      {"/exchanges/nanswap.svg", "", true},
	"simpleswap":   {"/exchanges/simpleswap.svg", "/exchanges/no-text/simpleswap.svg", true},
	"wizardswap":   {"/exchanges/wizardswap.png", "", false},
	"stealthex":    {"/exchanges/stealthex.svg", "/exchanges/no-text/stealthex.png", true},
	"exolix":       {"/exchanges/exolix.png", "/exchanges/no-text/exolix.png", true},
	"swapuz":       {"/exchanges/swapuz.svg", "/exchanges/no-text/swapuz.png", false},
	"bitcoinvn":    {"/exchanges/bitcoinvn.png", "/exchanges/no-text/bitcoinvn.png", false},
	"xgram":        {"/exchanges/xgram.svg", "", true},
	"pegasusswap":  {"/exchanges/pegasusswap.png", "", false},
	"godex":        {"/exchanges/godex.svg", "/exchanges/no-text/godex.jpg", true},
	"letsexchange": {"/exchanges/letsexchange.svg", "/exchanges/no-text/letsexchange.svg", true},
	"quickex":      {"/exchanges/quickex.png", "", false},
	"silkbyte":     {"/exchanges/silkbyte.png", "", false},
	"etzswap":      {"/exchanges/etz.png", "", false},
	"thorchain":    {"/exchanges/thorchain.png", "/exchanges/no-text/thorchain.png", false},
}

func parseCoinValue(value string) (string, string) {
	re := regexp.MustCompile(`^([^\(]+)(?:\(([^)]+)\))?$`)
	matches := re.FindStringSubmatch(value)
	if len(matches) == 3 {
		ticker := matches[1]
		network := matches[2]
		if network == "" {
			network = ticker
		}
		return ticker, network
	}
	return value, value
}

func EstimateHandler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")

	coin1 = strings.ToLower(coin1)
	coin2 = strings.ToLower(coin2)

	amountStr := c.QueryParam("amount")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")
	mode := c.QueryParam("mode")

	nojs := c.QueryParam("nojs")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing amount")
	}

	if nojs == "true" {

		sorting := c.QueryParam("sort")

		query := url.Values{}
		for key, values := range c.QueryParams() {
			if key == "sort" {
				continue
			}
			for _, v := range values {
				query.Add(key, v)
			}
		}
		baseQuery := query.Encode()

		if sorting != "kyc" {
			sorting = "rate"
		}

		amountFloat := amount
		if err != nil {
			return views.EstimateTempl(err.Error(), []api.Estimate{}, "", baseQuery).Render(c.Request().Context(), c.Response())
		}
		if amountFloat == 0 {
			return views.EstimateTempl("amount can't be zero", []api.Estimate{}, "", baseQuery).Render(c.Request().Context(), c.Response())
		}

		coin1Ticker, network1 := parseCoinValue(coin1)
		coin2Ticker, network2 := parseCoinValue(coin2)

		estimates, err := api.FetchEstimateFromAPI(coin1Ticker, coin2Ticker, amountFloat, false, network1, network2)
		if err != nil {
			return views.EstimateTempl(err.Error(), []api.Estimate{}, "", baseQuery).Render(c.Request().Context(), c.Response())
		}

		isAnonymousNetwork, _ := c.Get("isAnonymousNetwork").(bool)

		value_btc := estimates.TradeValue_btc
		value_usd := estimates.TradeValue_fiat

		for i := range estimates.Results {
			// fmt.Println(apiEstimates[i].ExchangeName)

			name := strings.ToLower(estimates.Results[i].ExchangeName)

			if info, ok := exchangeInfo[name]; ok {
				estimates.Results[i].ImageURL = info.ImageURL
				estimates.Results[i].NoTextURL = info.NoTextURL
				for _, ex := range views.ExchangesList {
					if strings.ToLower(ex.ShortCode) == name {
						estimates.Results[i].CGShield = ex.CGShield
						if ex.CGShield {
							amount := ex.CGSAmountFloat
							fiat := ex.CGSinStable
							if fiat {
								estimates.Results[i].CoveragePercent = math.Min(math.Max(amount/value_usd*100, 0), 100)
							} else {
								estimates.Results[i].CoveragePercent = math.Min(math.Max(amount/value_btc*100, 0), 100)
							}
						}
						break
					}
				}

			} else {
				estimates.Results[i].ImageURL = ""
			}

			estimates.Results[i].Log = exchangeInfo[name].RequireIP
			estimates.Results[i].Blocked = isExchangeBlocked(estimates.Results[i].ExchangeName, isAnonymousNetwork)

		}

		if sorting == "kyc" {
			sort.SliceStable(estimates.Results, func(i, j int) bool {
				if estimates.Results[i].KYCScore == estimates.Results[j].KYCScore {
					return estimates.Results[i].ReceiveAmount > estimates.Results[j].ReceiveAmount
				}
				return estimates.Results[i].KYCScore < estimates.Results[j].KYCScore
			})
		}

		return views.EstimateTempl("", estimates.Results, sorting, baseQuery).Render(c.Request().Context(), c.Response())
	}

	var estimates api.Estimates

	if mode == "pay" {
		estimates, err = api.FetchPaymentEstimateFromAPI(coin1, coin2, amount, network1, network2)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
		}
	} else {
		estimates, err = api.FetchEstimateFromAPI(coin1, coin2, amount, false, network1, network2)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
		}
	}

	isAnonymousNetwork, _ := c.Get("isAnonymousNetwork").(bool)

	value_btc := estimates.TradeValue_btc
	value_usd := estimates.TradeValue_fiat
	for i := range estimates.Results {

		name := strings.ToLower(estimates.Results[i].ExchangeName)

		if info, ok := exchangeInfo[name]; ok {
			estimates.Results[i].ImageURL = info.ImageURL
			estimates.Results[i].NoTextURL = info.NoTextURL
			for _, ex := range views.ExchangesList {
				if strings.ToLower(ex.ShortCode) == name {
					estimates.Results[i].CGShield = ex.CGShield
					if ex.CGShield {
						amount := ex.CGSAmountFloat
						fiat := ex.CGSinStable
						if fiat {
							estimates.Results[i].CoveragePercent = math.Min(math.Max(amount/value_usd*100, 0), 100)
						} else {
							estimates.Results[i].CoveragePercent = math.Min(math.Max(amount/value_btc*100, 0), 100)
						}
					}
					break
				}
			}
		}

		estimates.Results[i].Log = exchangeInfo[name].RequireIP
		estimates.Results[i].Blocked = isExchangeBlocked(estimates.Results[i].ExchangeName, isAnonymousNetwork)

	}

	return views.EstimateCard(estimates, mode).Render(c.Request().Context(), c.Response())
}

func Step2Handler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")
	amountStr := c.QueryParam("amount")
	mode := c.QueryParam("mode")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing amount")
	}
	partner := c.QueryParam("partner")
	receiveamountStr := c.QueryParam("receiveamount")
	receiveamount, err := strconv.ParseFloat(receiveamountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing receive amount")
	}

	var estimate api.Estimate
	estimate.Coin1 = coin1
	estimate.Coin2 = coin2
	estimate.ExchangeName = partner
	estimate.Network1 = network1
	estimate.Network2 = network2

	// In payment mode, amount is what user wants to receive, receiveamount is what they need to send
	if mode == "pay" {
		estimate.SendAmount = receiveamount // What user needs to send
		estimate.ReceiveAmount = amount     // What user wants to receive
	} else {
		estimate.SendAmount = amount           // What user wants to send
		estimate.ReceiveAmount = receiveamount // What user will receive
	}

	for _, ex := range views.ExchangesList {
		if strings.EqualFold(ex.ShortCode, partner) {
			estimate.CGSAmount = ex.CGAmount
			break
		}
	}
	name := strings.ToLower(estimate.ExchangeName)
	if info, ok := exchangeInfo[name]; ok {
		estimate.ImageURL = info.ImageURL
	}

	return views.AddressForm(estimate, false, mode).Render(c.Request().Context(), c.Response())
}

func Step3Handler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")
	amountStr := c.QueryParam("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "error parsing amount")
	}
	partner := c.QueryParam("partner")
	address := c.QueryParam("address")
	mode := c.QueryParam("mode")

	info := api.Info{
		IP:        "",
		UserAgent: "",
		LangList:  "",
	}
	partnerLower := strings.ToLower(partner)

	if exchangeData, ok := exchangeInfo[partnerLower]; ok && exchangeData.RequireIP {
		ip := c.RealIP()
		userAgent := c.Request().Header.Get("User-Agent")
		langList := c.Request().Header.Get("Accept-Language")

		info = api.Info{
			IP:        ip,
			UserAgent: userAgent,
			LangList:  langList,
		}
	}

	affiliate := ""
	affiliateCookie, err := c.Cookie("affiliate")
	if err == nil && affiliateCookie.Value != "" {
		affiliate = affiliateCookie.Value
	}

	// Get anonymous network status from middleware
	isAnonymousNetwork, _ := c.Get("isAnonymousNetwork").(bool)
	source := "clearnet-main"
	if isAnonymousNetwork {
		source = "anonet"
	}

	var transaction api.Transaction

	// Check if it's payment mode and call the appropriate API
	if mode == "pay" {
		// In payment mode, 'amount' is what the user wants to receive
		err, transaction = api.CreatePaymentFromAPI(coin1, coin2, amount, address, partner, network1, network2, affiliate, info, source)
	} else {
		// In swap mode, 'amount' is what the user wants to send
		err, transaction = api.CreateTradeFromAPI(coin1, coin2, amount, address, partner, network1, network2, affiliate, info, source)
	}

	if err != nil || transaction.Id == "" {
		var estimate api.Estimate
		estimate.Coin1 = coin1
		estimate.Coin2 = coin2
		estimate.SendAmount = amount
		estimate.ExchangeName = partner
		estimate.ReceiveAmount = amount
		estimate.Network1 = network1
		estimate.Network2 = network2

		return views.AddressForm(estimate, true, mode).Render(c.Request().Context(), c.Response())
	}

	return c.Redirect(http.StatusSeeOther, "/transaction/"+transaction.CGID)
}

func CGPayHandler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	network1 := c.QueryParam("network1")
	address := c.QueryParam("address")

	amount_str := c.QueryParam("amount")
	amount, err := strconv.ParseFloat(amount_str, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "error parsing amount")
	}

	return views.CGPayForm(coin1, network1, amount, address, "").Render(c.Request().Context(), c.Response())
}

func CGPayCreateHandler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")
	amountStr := c.QueryParam("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "error parsing amount")
	}
	partner := c.QueryParam("partner")
	address := c.QueryParam("address")

	info := api.Info{
		IP:        "",
		UserAgent: "",
		LangList:  "",
	}
	partnerLower := strings.ToLower(partner)

	if exchangeData, ok := exchangeInfo[partnerLower]; ok && exchangeData.RequireIP {
		ip := c.RealIP()
		userAgent := c.Request().Header.Get("User-Agent")
		langList := c.Request().Header.Get("Accept-Language")

		info = api.Info{
			IP:        ip,
			UserAgent: userAgent,
			LangList:  langList,
		}
	}

	affiliate := ""
	affiliateCookie, err := c.Cookie("affiliate")
	if err == nil && affiliateCookie.Value != "" {
		affiliate = affiliateCookie.Value
	}

	// Get anonymous network status from middleware
	isAnonymousNetwork, _ := c.Get("isAnonymousNetwork").(bool)
	source := "clearnet-main"
	if isAnonymousNetwork {
		source = "anonet"
	}

	err, transaction := api.CreateQuickPaymentFromAPI(coin1, coin2, amount, address, network1, network2, affiliate, info, source)

	if err != nil {
		fmt.Printf("API Error: %v\n", err)
		return views.CGPayForm(coin1, network1, amount, address, err.Error()).Render(c.Request().Context(), c.Response())
	}

	if transaction.Id == "" && transaction.CGID == "" {
		fmt.Printf("Transaction missing both Id and CGID fields\n")
		return views.CGPayForm(coin1, network1, amount, address, "Transaction creation failed - no transaction ID returned").Render(c.Request().Context(), c.Response())
	}

	transactionID := transaction.CGID
	if transactionID == "" {
		transactionID = transaction.Id
	}

	if transactionID == "" {
		fmt.Printf("Both transaction.CGID and transaction.Id are empty\n")
		return views.CGPayForm(coin1, network1, amount, address, "Transaction creation failed - invalid transaction ID").Render(c.Request().Context(), c.Response())
	}

	return c.Redirect(http.StatusSeeOther, "/transaction/"+transactionID+"?cg-pay=true")
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func GetTransactionHandler(c echo.Context) error {
	id := c.Param("id")
	cgPay := false
	if v := c.QueryParam("cg-pay"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cgPay = b
		}
	}

	err, transaction := api.GetTransactionFromAPI(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error getting transaction: %v", err))
	}

	return views.TransactionPage(transaction, cgPay).Render(c.Request().Context(), c.Response())

}

func getCoinName(ticker string) string {
	defaultName := strings.ToUpper(ticker)

	jsonFile, err := os.ReadFile("static/coins.json")
	if err != nil {
		fmt.Printf("Error reading coins.json: %v\n", err)
		return defaultName
	}

	type Coin struct {
		Ticker  string `json:"ticker"`
		Name    string `json:"name"`
		Network string `json:"network"`
	}

	var coins []Coin
	if err := json.Unmarshal(jsonFile, &coins); err != nil {
		fmt.Printf("Error parsing coins.json: %v\n", err)
		return defaultName
	}

	for _, coin := range coins {
		if strings.EqualFold(coin.Ticker, ticker) &&
			(coin.Network == ticker || coin.Network == "") {
			return coin.Name
		}
	}

	for _, coin := range coins {
		if strings.EqualFold(coin.Ticker, ticker) {
			return coin.Name
		}
	}

	return defaultName
}

func CoinExchangeScreenHandler(c echo.Context) error {
	coin := c.Param("coin")

	coinName, exists := coinExists(coin)
	if !exists {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return views.CoinExchangeScreen(coin, coinName).Render(c.Request().Context(), c.Response())
}

func coinExists(ticker string) (string, bool) {
	jsonFile, err := os.ReadFile("static/coins.json")
	if err != nil {
		fmt.Printf("Error reading coins.json: %v\n", err)
		return "", false
	}

	type Coin struct {
		Ticker  string `json:"ticker"`
		Name    string `json:"name"`
		Network string `json:"network"`
	}

	var coins []Coin
	if err := json.Unmarshal(jsonFile, &coins); err != nil {
		fmt.Printf("Error parsing coins.json: %v\n", err)
		return "", false
	}

	for _, coin := range coins {
		if strings.EqualFold(coin.Ticker, ticker) &&
			(coin.Network == ticker || coin.Network == "") {
			return coin.Name, true
		}
	}

	for _, coin := range coins {
		if strings.EqualFold(coin.Ticker, ticker) {
			return coin.Name, true
		}
	}

	return "", false
}

func HealthHandler(c echo.Context) error {
	c.String(http.StatusOK, "OK")
	return nil
}

func RobotsHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/plain")

	robots := `User-agent: *
Allow: /
Sitemap: https://cyphergoat.com/sitemap.xml
`

	return c.String(http.StatusOK, robots)
}
func SitemapHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/xml")

	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	// Add static pages
	staticPages := []struct {
		Path       string
		ChangeFreq string
		Priority   string
	}{
		{"/", "daily", "1.0"},
		{"/about", "monthly", "0.8"},
		{"/terms", "monthly", "0.7"},
		{"/privacy", "monthly", "0.7"},
		{"/contact", "monthly", "0.8"},
		{"/affiliate", "weekly", "0.9"},
		{"/affiliate/terms", "monthly", "0.7"},
		{"/blog", "daily", "0.9"},
		{"/this-week-in-monero", "weekly", "0.9"},
	}

	baseURL := "https://cyphergoat.com"

	for _, page := range staticPages {
		sitemap += fmt.Sprintf(`
  <url>
    <loc>%s%s</loc>
    <changefreq>%s</changefreq>
    <priority>%s</priority>
  </url>`, baseURL, page.Path, page.ChangeFreq, page.Priority)
	}

	jsonFile, err := os.ReadFile("static/coins.json")
	if err == nil {
		type Coin struct {
			Ticker  string `json:"ticker"`
			Name    string `json:"name"`
			Network string `json:"network"`
		}

		var coins []Coin
		if err := json.Unmarshal(jsonFile, &coins); err == nil {
			uniqueTickers := make(map[string]bool)

			for _, coin := range coins {
				ticker := strings.ToLower(coin.Ticker)
				if _, exists := uniqueTickers[ticker]; !exists {
					uniqueTickers[ticker] = true

					sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/swap/%s</loc>
    <changefreq>daily</changefreq>
    <priority>0.8</priority>
  </url>`, baseURL, ticker)
				}
			}
		}
	}

	posts, err := loadAllBlogPosts()
	if err == nil {
		for _, post := range posts {
			if strings.HasPrefix(post.Slug, "twim-") {
				issueSlug := strings.Replace(post.Slug, "twim-", "issue-", 1)
				sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/this-week-in-monero/%s</loc>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>`, baseURL, issueSlug)
			} else {
				sitemap += fmt.Sprintf(`
  <url>
    <loc>%s/blog/%s</loc>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>`, baseURL, post.Slug)
			}
		}
	}

	sitemap += `
</urlset>`

	return c.String(http.StatusOK, sitemap)
}

// func AffiliateForm(c echo.Context) error {
// 	if c.Request().Method == "POST" {
// 		email := c.FormValue("email")
// 		messenger := c.FormValue("messenger")
// 		promotion := c.FormValue("promotion")
// 		agree := c.FormValue("agree")

// 		var errors []string
// 		if email == "" {
// 			errors = append(errors, "Email is required.")
// 		}
// 		if messenger == "" {
// 			errors = append(errors, "Messenger username is required.")
// 		}
// 		if promotion == "" {
// 			errors = append(errors, "Promotion details are required.")
// 		}
// 		if agree != "on" {
// 			errors = append(errors, "You must agree to the terms.")
// 		}

// 		if len(errors) > 0 {
// 			return views.AffiliateForm(errors, "").Render(c.Request().Context(), c.Response())
// 		}

// 		// Send to API
// 		app := api.AffiliateApplication{
// 			Email:     email,
// 			Messenger: messenger,
// 			Promotion: promotion,
// 		}
// 		if err := api.SendAffiliateApplication(app); err != nil {
// 			errors = append(errors, "Failed to submit application. Please try again later.")
// 			return views.AffiliateForm(errors, "").Render(c.Request().Context(), c.Response())
// 		}

// 		return views.AffiliateForm(nil, "Thank you for applying! We'll review your application soon.").Render(c.Request().Context(), c.Response())
// 	}
// 	return views.AffiliateForm(nil, "").Render(c.Request().Context(), c.Response())
// }

func BlogHandler(c echo.Context) error {
	posts, err := loadAllBlogPosts()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to load blog posts")
	}

	// Filter out TWIM posts from regular blog
	regularPosts := []*views.BlogPostData{}
	for _, post := range posts {
		if !strings.HasPrefix(post.Slug, "twim-") {
			regularPosts = append(regularPosts, post)
		}
	}

	return render(c, views.Blog(regularPosts))
}

func TWIMHandler(c echo.Context) error {
	posts, err := loadAllBlogPosts()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to load blog posts")
	}

	// Filter blog posts to only show TWIM articles
	twimPosts := []*views.BlogPostData{}
	for _, post := range posts {
		if strings.HasPrefix(post.Slug, "twim-") {
			post.IssueSlug = strings.Replace(post.Slug, "twim-", "issue-", 1)
			twimPosts = append(twimPosts, post)
		}
	}

	return render(c, views.TWIMIndex(twimPosts))
}

func TWIMPostHandler(c echo.Context) error {
	slug := c.Param("slug")

	var fileSlug string
	if strings.HasPrefix(slug, "issue-") {
		fileSlug = strings.Replace(slug, "issue-", "twim-", 1)
	} else {
		// Redirect old twim-N URLs to issue-N
		if strings.HasPrefix(slug, "twim-") {
			issueSlug := strings.Replace(slug, "twim-", "issue-", 1)
			return c.Redirect(http.StatusMovedPermanently, "/this-week-in-monero/"+issueSlug)
		}
		return c.Redirect(http.StatusMovedPermanently, "/this-week-in-monero")
	}

	post, err := loadBlogPost(fileSlug)
	if err != nil || post == nil {
		return c.Redirect(http.StatusMovedPermanently, "/this-week-in-monero")
	}

	return render(c, views.TWIMPost(post))
}

func BlogPostHandler(c echo.Context) error {
	slug := c.Param("slug")

	// Redirect TWIM posts to new location
	if strings.HasPrefix(slug, "twim-") {
		return c.Redirect(http.StatusMovedPermanently, "/this-week-in-monero/"+slug)
	}

	post, err := loadBlogPost(slug)
	if err != nil || post == nil {
		return echo.NewHTTPError(404, "Blog post not found")
	}
	return render(c, views.BlogPost(post))
}

func TWIMRedirectHandler(c echo.Context) error {
	slug := c.Param("slug")

	if strings.HasPrefix(slug, "twim-") {
		issueSlug := strings.Replace(slug, "twim-", "issue-", 1)
		return c.Redirect(http.StatusMovedPermanently, "/this-week-in-monero/"+issueSlug)
	}

	post, err := loadBlogPost(slug)
	if err != nil || post == nil {
		return echo.NewHTTPError(404, "Blog post not found")
	}
	return render(c, views.BlogPost(post))
}

func loadAllBlogPosts() ([]*views.BlogPostData, error) {
	var posts []*views.BlogPostData

	blogDir := "content/blog"
	entries, err := os.ReadDir(blogDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		slug := strings.TrimSuffix(entry.Name(), ".md")
		post, err := loadBlogPostFromFile(filepath.Join(blogDir, entry.Name()))
		if err != nil {
			continue
		}

		if post.Draft {
			continue
		}

		post.Slug = slug
		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", posts[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", posts[j].Date)

		if errI != nil || errJ != nil {
			return posts[i].Date > posts[j].Date
		}

		return dateI.After(dateJ)
	})

	return posts, nil
}

func TWIMRSSHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "application/rss+xml")

	posts, err := loadAllBlogPosts()
	if err != nil {
		return err
	}

	var twimPosts []*views.BlogPostData
	for _, post := range posts {
		if strings.HasPrefix(post.Slug, "twim-") && !post.Draft {
			twimPosts = append(twimPosts, post)
		}
	}

	sort.Slice(twimPosts, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", twimPosts[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", twimPosts[j].Date)

		if errI != nil || errJ != nil {
			return twimPosts[i].Date > twimPosts[j].Date
		}

		return dateI.After(dateJ)
	})

	if len(twimPosts) > 10 {
		twimPosts = twimPosts[:10]
	}

	rss := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>This Week In Monero - CypherGoat</title>
    <link>https://cyphergoat.com/this-week-in-monero</link>
    <description>Weekly insights and updates from the Monero ecosystem. Stay informed on the latest developments, news and analysis in the world of private digital currency.</description>
    <language>en-us</language>
    <managingEditor>recanman@cyphergoat.com (recanman)</managingEditor>
    <webMaster>support@cyphergoat.com (CypherGoat)</webMaster>
    <lastBuildDate>` + time.Now().Format(time.RFC1123Z) + `</lastBuildDate>
    <atom:link href="https://cyphergoat.com/this-week-in-monero/rss.xml" rel="self" type="application/rss+xml" />
    <image>
      <url>https://cyphergoat.com/static/icons/cg-notext-light.png</url>
      <title>This Week In Monero - CypherGoat</title>
      <link>https://cyphergoat.com/this-week-in-monero</link>
    </image>`

	for _, post := range twimPosts {
		issueSlug := strings.Replace(post.Slug, "twim-", "issue-", 1)
		itemURL := fmt.Sprintf("https://cyphergoat.com/this-week-in-monero/%s", issueSlug)

		description := post.Description
		if description == "" && len(post.HTMLContent) > 200 {
			description = post.HTMLContent[:200] + "..."
		}

		pubDate := post.Date
		if parsedDate, err := time.Parse("2006-01-02", post.Date); err == nil {
			pubDate = parsedDate.Format(time.RFC1123Z)
		}

		rss += fmt.Sprintf(`
    <item>
      <title>%s</title>
      <link>%s</link>
      <guid>%s</guid>
      <pubDate>%s</pubDate>
      <author>%s</author>
      <description><![CDATA[%s]]></description>
    </item>`,
			html.EscapeString(post.Title),
			itemURL,
			itemURL,
			pubDate,
			html.EscapeString(post.Author),
			description)
	}

	rss += `
  </channel>
</rss>`

	return c.String(http.StatusOK, rss)
}

func loadBlogPost(slug string) (*views.BlogPostData, error) {
	// Validate slug to ensure it's a simple filename without path traversal
	if strings.Contains(slug, "/") || strings.Contains(slug, "\\") || strings.Contains(slug, "..") {
		return nil, fmt.Errorf("invalid blog post slug")
	}

	validSlug := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(slug)
	if !validSlug {
		return nil, fmt.Errorf("invalid blog post slug")
	}

	filename := filepath.Join("content/blog", slug+".md")
	post, err := loadBlogPostFromFile(filename)
	if err != nil {
		return nil, err
	}

	if post.Draft {
		return nil, fmt.Errorf("post is draft")
	}

	post.Slug = slug
	return post, nil
}

func loadBlogPostFromFile(filename string) (*views.BlogPostData, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(content), "---")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid markdown format: missing frontmatter")
	}

	var post views.BlogPostData
	err = yaml.Unmarshal([]byte(parts[1]), &post)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	markdownContent := strings.TrimSpace(strings.Join(parts[2:], "---"))
	post.Content = markdownContent

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithUnsafe(), // Allow raw HTML in markdown
			goldmarkhtml.WithXHTML(),
		),
	)

	var buf strings.Builder
	err = md.Convert([]byte(markdownContent), &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}
	post.HTMLContent = buf.String()

	return &post, nil
}

func ChallengeHandler(c echo.Context) error {
	data := views.FormData{
		Nonce: c.FormValue("nonce"),
		Email: c.FormValue("email"),
		Name:  c.FormValue("name"),
		Hp:    c.FormValue("hp_field"),
	}

	// Honeypot check
	if data.Hp != "" {
		return c.String(http.StatusBadRequest, "Bot detected")
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return render(c, views.JSChallenge(data))
}

func AffiliateLoginHandler(c echo.Context) error {
	if c.Request().Method == "POST" {
		username := c.FormValue("username")
		password := c.FormValue("password")

		var errors []string
		if username == "" {
			errors = append(errors, "Username is required.")
		}
		if password == "" {
			errors = append(errors, "Password is required.")
		}

		if len(errors) > 0 {
			return views.AffiliateLogin(errors, "").Render(c.Request().Context(), c.Response())
		}

		token, err := api.LoginAffiliate(username, password)
		if err != nil {
			errors = append(errors, "Invalid credentials.")
			return views.AffiliateLogin(errors, "").Render(c.Request().Context(), c.Response())
		}

		cookie := new(http.Cookie)
		cookie.Name = "affiliate_token"
		cookie.Value = token
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.MaxAge = 24 * 60 * 60 * 7 // 7 days
		c.SetCookie(cookie)

		return c.Redirect(http.StatusSeeOther, "/affiliate/dashboard")
	}
	return views.AffiliateLogin(nil, "").Render(c.Request().Context(), c.Response())
}

func AffiliateRegisterHandler(c echo.Context) error {
	if c.Request().Method == "POST" {
		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")
		confirmPassword := c.FormValue("confirm_password")

		var errors []string
		if username == "" {
			errors = append(errors, "Username is required.")
		}
		if email == "" {
			errors = append(errors, "Email is required.")
		}
		if password == "" {
			errors = append(errors, "Password is required.")
		}
		if confirmPassword != password {
			errors = append(errors, "Passwords do not match.")
		}
		if len(password) < 8 {
			errors = append(errors, "Password must be at least 8 characters long.")
		}

		if len(errors) > 0 {
			return views.AffiliateRegister(errors, "", "").Render(c.Request().Context(), c.Response())
		}

		err := api.RegisterAffiliate(username, email, password, "")
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				errors = append(errors, "Registration endpoint not available. Please contact support.")
			} else {
				errors = append(errors, "Registration failed: "+err.Error())
			}
			return views.AffiliateRegister(errors, "", "").Render(c.Request().Context(), c.Response())
		}

		return views.AffiliateRegister(nil, "Registration successful! You can now log in.", "").Render(c.Request().Context(), c.Response())
	}

	return views.AffiliateRegister(nil, "", "").Render(c.Request().Context(), c.Response())
}

func AffiliateDashboardHandler(c echo.Context) error {
	cookie, err := c.Cookie("affiliate_token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/affiliate/login")
	}

	stats, err := api.GetAffiliateStats(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/affiliate/login")
	}

	return views.AffiliateDashboard(stats).Render(c.Request().Context(), c.Response())
}

func AffiliateLogoutHandler(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "affiliate_token"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/affiliate/login")
}

func AffiliatePayoutHandler(c echo.Context) error {
	cookie, err := c.Cookie("affiliate_token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/affiliate/login")
	}

	if c.Request().Method == "POST" {
		amountStr := c.FormValue("amount")

		var errors []string

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount < 200.00 {
			errors = append(errors, "Minimum payout amount is $200.00")
		}

		payoutCurrency := c.FormValue("payout_currency")
		walletAddress := c.FormValue("wallet_address")

		if payoutCurrency != "USDC_SOL" && payoutCurrency != "XMR" {
			errors = append(errors, "Please select a valid payout currency")
		}

		if payoutCurrency == "USDC_SOL" {
			if len(walletAddress) < 32 || len(walletAddress) > 44 {
				errors = append(errors, "Invalid Solana address format (32-44 characters)")
			}
		} else if payoutCurrency == "XMR" {
			if len(walletAddress) < 95 || len(walletAddress) > 106 {
				errors = append(errors, "Invalid Monero address format (95-106 characters)")
			}
		}

		if len(errors) > 0 {
			stats, _ := api.GetAffiliateStats(cookie.Value)
			return views.AffiliatePayout(errors, "", stats).Render(c.Request().Context(), c.Response())
		}

		request := api.PayoutRequest{
			Amount:         amount,
			PayoutCurrency: payoutCurrency,
			WalletAddress:  walletAddress,
		}

		_, err = api.RequestPayout(cookie.Value, request)
		if err != nil {
			errors = append(errors, "Payout request failed: "+err.Error())
			stats, _ := api.GetAffiliateStats(cookie.Value)
			return views.AffiliatePayout(errors, "", stats).Render(c.Request().Context(), c.Response())
		}

		stats, _ := api.GetAffiliateStats(cookie.Value)
		return views.AffiliatePayout(nil, "Payout request submitted successfully!", stats).Render(c.Request().Context(), c.Response())
	}

	stats, err := api.GetAffiliateStats(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/affiliate/login")
	}

	return views.AffiliatePayout(nil, "", stats).Render(c.Request().Context(), c.Response())
}

func AffiliatePayoutHistoryHandler(c echo.Context) error {
	cookie, err := c.Cookie("affiliate_token")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/affiliate/login")
	}

	history, err := api.GetPayoutHistory(cookie.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load payout history")
	}

	return views.AffiliatePayoutHistory(history).Render(c.Request().Context(), c.Response())
}
