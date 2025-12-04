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

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/CypherGoat/web/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skip2/go-qrcode"
)

func CacheMiddleware(duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			maxAge := int(duration.Seconds())
			c.Response().Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			c.Response().Header().Set("Expires", time.Now().Add(duration).Format(http.TimeFormat))
			return next(c)
		}
	}
}

func AffiliateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		affiliate := c.QueryParam("ref")
		if affiliate != "" {
			c.Set("affiliate", affiliate)
			cookie := new(http.Cookie)
			cookie.Name = "affiliate"
			cookie.Value = affiliate
			cookie.Path = "/"
			cookie.HttpOnly = true
			cookie.Secure = true
			c.SetCookie(cookie)

		}
		return next(c)
	}
}

func AllowedHostsMiddleware(allowedHosts []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := c.Request().Host

			if colonIndex := strings.LastIndex(host, ":"); colonIndex != -1 {
				host = host[:colonIndex]
			}

			for _, allowedHost := range allowedHosts {
				if host == allowedHost {
					return next(c)
				}
			}

			return c.String(http.StatusForbidden, "Access denied: Host not allowed")
		}
	}
}

func main() {
	e := echo.New()

	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	allowedHosts := []string{
		"cyphergoat.com",
		"www.cyphergoat.com",
		"localhost",
		"127.0.0.1",
		"cyphergmw4huw7jzhat3misfm5jj2m4nvafockqbj7i5rrlec6mobdid.onion",
		"cjorx2hmv3jxs5sshxhw73nmkiayupim6gfye3sfts2pzmjgk5rq.b32.i2p",
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(AllowedHostsMiddleware(allowedHosts))
	e.Use(allowPayEmbedding)
	e.Use(AffiliateMiddleware)

	e.Use(removeTrailingSlash)

	e.Static("/static", "static")
	e.Static("/icons", "static/icons")
	e.Static("/exchanges", "static/exchanges")

	e.GET("/", handlers.IndexHandler)
	e.GET("/estimate", handlers.EstimateHandler)
	e.GET("/step2", handlers.Step2Handler)
	e.GET("/step3", handlers.Step3Handler)
	e.GET("/generate", generateQRCodeHandler)
	e.GET("/transaction/:id", handlers.GetTransactionHandler)
	e.GET("/about", handlers.AboutHandler)
	e.GET("/terms", handlers.TermsHandler)
	e.GET("/contact", handlers.ContactHandler)

	e.GET("/privacy", handlers.PrivacyPolicyHandler)

	e.GET("/health", handlers.HealthHandler)
	e.GET("/affiliate/terms", handlers.AffiliateTerms)

	e.GET("/affiliate/payout", handlers.AffiliatePayoutHandler)
	e.POST("/affiliate/payout", handlers.AffiliatePayoutHandler)
	e.GET("/affiliate/payout/history", handlers.AffiliatePayoutHistoryHandler)

	e.GET("/swap/:coin", handlers.CoinExchangeScreenHandler)

	e.GET("/robots.txt", handlers.RobotsHandler)

	e.GET("/sitemap.xml", handlers.SitemapHandler)

	e.GET("/exchanges", handlers.ExchangesHandler)
	e.GET("/exchange/:shortcode", handlers.ExchangeDetailHandler)

	e.GET("/cyphergoat-shield", handlers.CGShieldTermsHandler)

	e.POST("/challenge", handlers.ChallengeHandler)

	e.GET("/blog", handlers.BlogHandler, CacheMiddleware(1*time.Hour))
	e.GET("/blog/:slug", handlers.TWIMRedirectHandler, CacheMiddleware(1*time.Hour))

	e.GET("/this-week-in-monero", handlers.TWIMHandler, CacheMiddleware(1*time.Hour))
	e.GET("/this-week-in-monero/:slug", handlers.TWIMPostHandler, CacheMiddleware(1*time.Hour))
	e.GET("/this-week-in-monero/rss.xml", handlers.TWIMRSSHandler, CacheMiddleware(1*time.Hour))

	e.Static("/blog/images", "static/blog")

	e.GET("/affiliate", handlers.AffiliateHandler)

	e.GET("/pay", handlers.CGPayHandler)
	e.GET("/pay/create", handlers.CGPayCreateHandler)

	e.GET("/affiliate/login", handlers.AffiliateLoginHandler)
	e.POST("/affiliate/login", handlers.AffiliateLoginHandler)
	e.GET("/affiliate/register", handlers.AffiliateRegisterHandler)
	e.POST("/affiliate/register", handlers.AffiliateRegisterHandler)
	e.GET("/affiliate/dashboard", handlers.AffiliateDashboardHandler)
	e.GET("/affiliate/logout", handlers.AffiliateLogoutHandler)

	e.GET("/affiliate/terms", handlers.AffiliateTerms)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == http.StatusNotFound {
			handlers.NotFoundHandler(c)
		} else {
			c.String(code, err.Error())
		}
	}
	e.Logger.Fatal(e.Start(":4200"))
}

func generateQRCodeHandler(c echo.Context) error {
	text := c.QueryParam("text")
	if text == "" {
		return c.String(http.StatusBadRequest, "Text query parameter is required")
	}
	png, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate QR code")
	}
	return c.Blob(http.StatusOK, "image/png", png)
}

func removeTrailingSlash(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path

		if path == "/" {
			return next(c)
		}

		if strings.HasSuffix(path, "/") {
			targetPath := strings.TrimSuffix(path, "/")
			query := c.Request().URL.RawQuery
			if query != "" {
				targetPath = targetPath + "?" + query
			}
			return c.Redirect(http.StatusMovedPermanently, targetPath)
		}

		return next(c)
	}
}

func allowPayEmbedding(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if path == "/pay" || strings.HasPrefix(path, "/pay/") {
			c.Response().Header().Set("X-Frame-Options", "ALLOWALL")
		}
		return next(c)
	}
}
