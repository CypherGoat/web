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
	"net/http"
	"os"
	"path/filepath"
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
	"github.com/yuin/goldmark/renderer/html"
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

type ExchangeInfo struct {
	ImageURL  string
	RequireIP bool
}

var exchangeInfo = map[string]ExchangeInfo{
	"alfacash":     {"/exchanges/alfacash.svg", true},
	"changee":      {"/exchanges/changee.svg", true},
	"changehero":   {"/exchanges/changehero.svg", true},
	"changenow":    {"/exchanges/changenow.svg", true},
	"coincraddle":  {"/exchanges/coincraddle.png", true},
	"exch.cx":      {"/exchanges/exchcx.png", false},
	"fixedfloat":   {"/exchanges/fixedfloat-v2.svg", true},
	"majesticbank": {"/exchanges/majesticbank.png", false},
	"nanswap":      {"/exchanges/nanswap.svg", true},
	"simpleswap":   {"/exchanges/simpleswap.svg", true},
	"wizardswap":   {"/exchanges/wizardswap.png", false},
	"stealthex":    {"/exchanges/stealthex.svg", true},
	"exolix":       {"/exchanges/exolix.png", true},
	"swapuz":       {"/exchanges/swapuz.svg", false},
	"bitcoinvn":    {"/exchanges/bitcoinvn.png", false},
	"xgram":        {"/exchanges/xgram.svg", true},
	"pegasusswap":  {"/exchanges/pegasusswap.png", false},
	"godex":        {"/exchanges/godex.svg", true},
}

func EstimateHandler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")
	amountStr := c.QueryParam("amount")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing amount")
	}

	apiEstimates, err := api.FetchEstimateFromAPI(coin1, coin2, amount, false, network1, network2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
	}
	for i := range apiEstimates {

		name := strings.ToLower(apiEstimates[i].ExchangeName)

		if info, ok := exchangeInfo[name]; ok {
			apiEstimates[i].ImageURL = info.ImageURL
		}

		apiEstimates[i].Log = exchangeInfo[name].RequireIP

	}

	return views.EstimateCard(apiEstimates).Render(c.Request().Context(), c.Response())
}

func Step2Handler(c echo.Context) error {
	coin1 := c.QueryParam("coin1")
	coin2 := c.QueryParam("coin2")
	network1 := c.QueryParam("network1")
	network2 := c.QueryParam("network2")
	amountStr := c.QueryParam("amount")
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
	estimate.SendAmount = amount
	estimate.ExchangeName = partner
	estimate.ReceiveAmount = receiveamount
	estimate.Network1 = network1
	estimate.Network2 = network2

	name := strings.ToLower(estimate.ExchangeName)
	if info, ok := exchangeInfo[name]; ok {
		estimate.ImageURL = info.ImageURL
	}

	return views.AddressForm(estimate, false).Render(c.Request().Context(), c.Response())
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

	fmt.Printf("Partner: %s, Exists in map: %v, RequiresIP: %v\n", partner, partnerLower, exchangeInfo[partnerLower].RequireIP)

	affiliate := ""
	affiliateCookie, err := c.Cookie("affiliate")
	if err == nil && affiliateCookie.Value != "" {
		affiliate = affiliateCookie.Value
	}

	err, transaction := api.CreateTradeFromAPI(coin1, coin2, amount, address, partner, network1, network2, affiliate, info)
	if err != nil || transaction.Id == "" {
		var estimate api.Estimate
		estimate.Coin1 = coin1
		estimate.Coin2 = coin2
		estimate.SendAmount = amount
		estimate.ExchangeName = partner
		estimate.ReceiveAmount = amount
		estimate.Network1 = network1
		estimate.Network2 = network2

		return views.AddressForm(estimate, true).Render(c.Request().Context(), c.Response())
	}

	fmt.Printf("Redirecting to transaction: /transaction/%s\n", transaction.CGID)

	return c.Redirect(http.StatusSeeOther, "/transaction/"+transaction.CGID)
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func GetTransactionHandler(c echo.Context) error {
	id := c.Param("id")

	err, transaction := api.GetTransactionFromAPI(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error getting transaction: %v", err))
	}

	return views.TransactionPage(transaction).Render(c.Request().Context(), c.Response())

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

	sitemap += `
</urlset>`

	return c.String(http.StatusOK, sitemap)
}

func AffiliateForm(c echo.Context) error {
	if c.Request().Method == "POST" {
		email := c.FormValue("email")
		messenger := c.FormValue("messenger")
		promotion := c.FormValue("promotion")
		agree := c.FormValue("agree")

		var errors []string
		if email == "" {
			errors = append(errors, "Email is required.")
		}
		if messenger == "" {
			errors = append(errors, "Messenger username is required.")
		}
		if promotion == "" {
			errors = append(errors, "Promotion details are required.")
		}
		if agree != "on" {
			errors = append(errors, "You must agree to the terms.")
		}

		if len(errors) > 0 {
			return views.AffiliateForm(errors, "").Render(c.Request().Context(), c.Response())
		}

		// Send to API
		app := api.AffiliateApplication{
			Email:     email,
			Messenger: messenger,
			Promotion: promotion,
		}
		if err := api.SendAffiliateApplication(app); err != nil {
			errors = append(errors, "Failed to submit application. Please try again later.")
			return views.AffiliateForm(errors, "").Render(c.Request().Context(), c.Response())
		}

		return views.AffiliateForm(nil, "Thank you for applying! We'll review your application soon.").Render(c.Request().Context(), c.Response())
	}
	return views.AffiliateForm(nil, "").Render(c.Request().Context(), c.Response())
}

func BlogHandler(c echo.Context) error {
	posts, err := loadAllBlogPosts()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to load blog posts")
	}
	return render(c, views.Blog(posts))
}

func BlogPostHandler(c echo.Context) error {
	slug := c.Param("slug")
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

func loadBlogPost(slug string) (*views.BlogPostData, error) {
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
			html.WithUnsafe(), // Allow raw HTML in markdown
			html.WithXHTML(),
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
