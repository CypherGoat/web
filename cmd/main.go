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
	"net/http"
	"strings"
	"time"

	"github.com/CypherGoat/web/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skip2/go-qrcode"
)

func AffiliateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		affiliate := c.QueryParam("ref")
		if affiliate != "" {
			cookie := new(http.Cookie)
			cookie.Name = "affiliate"
			cookie.Value = affiliate
			cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
			cookie.Path = "/"
			c.SetCookie(cookie)
		}
		return next(c)
	}
}

func main() {
	e := echo.New()

	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
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

	e.GET("/affiliate", handlers.AffiliateHandler)

	e.GET("/health", handlers.HealthHandler)
	e.GET("/affiliate/terms", handlers.AffiliateTerms)

	e.GET("/swap/:coin", handlers.CoinExchangeScreenHandler)

	e.GET("/robots.txt", handlers.RobotsHandler)

	e.GET("/sitemap.xml", handlers.SitemapHandler)

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
