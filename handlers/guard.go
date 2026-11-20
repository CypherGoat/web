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

package handlers

// Web UI session guard.
//
// Pages that render the swap form hand out a signed session cookie. Endpoints
// that are only meant to be called from the web UI (e.g. /estimate) require
// that cookie, so scripts can't hit them directly without first loading a
// page, and rate limits can be applied per session as well as per IP.
// Works without JavaScript: the cookie is set on the page response and sent
// back by both htmx requests and plain form submits.

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	sessionCookieName = "cg_s"
	sessionTTL        = 12 * time.Hour
	sessionCtxKey     = "cg_session"
)

var sessionSecret = loadSessionSecret()

func loadSessionSecret() []byte {
	if s := os.Getenv("CG_SESSION_SECRET"); len(s) >= 32 {
		return []byte(s)
	}
	log.Println("guard: CG_SESSION_SECRET unset or shorter than 32 chars; using a random per-process secret (sessions reset on restart)")
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func signSession(payload string) string {
	m := hmac.New(sha256.New, sessionSecret)
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil)[:16])
}

// token format: <hex id>.<unix issued>.<sig>
func makeSessionToken(id string, issued time.Time) string {
	payload := id + "." + strconv.FormatInt(issued.Unix(), 10)
	return payload + "." + signSession(payload)
}

func parseSessionToken(tok string) (id string, issued time.Time, ok bool) {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 || len(parts[0]) != 32 {
		return "", time.Time{}, false
	}
	expected := signSession(parts[0] + "." + parts[1])
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return "", time.Time{}, false
	}
	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", time.Time{}, false
	}
	issued = time.Unix(ts, 0)
	age := time.Since(issued)
	if age < -time.Minute || age > sessionTTL {
		return "", time.Time{}, false
	}
	return parts[0], issued, true
}

func newSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func isAnonymousNetwork(c echo.Context) bool {
	v, _ := c.Get("isAnonymousNetwork").(bool)
	return v
}

func setSessionCookie(c echo.Context, id string) {
	c.SetCookie(&http.Cookie{
		Name:     sessionCookieName,
		Value:    makeSessionToken(id, time.Now()),
		Path:     "/",
		MaxAge:   int(sessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   !isAnonymousNetwork(c),
		SameSite: http.SameSiteLaxMode,
	})
}

// IssueSession sets (or refreshes) the session cookie. Attach it to pages
// that render the swap form.
func IssueSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if ck, err := c.Cookie(sessionCookieName); err == nil {
			if id, issued, ok := parseSessionToken(ck.Value); ok {
				// Keep the same id so per-session limits carry over.
				if time.Since(issued) > sessionTTL/2 {
					setSessionCookie(c, id)
				}
				c.Set(sessionCtxKey, id)
				return next(c)
			}
		}
		id := newSessionID()
		setSessionCookie(c, id)
		c.Set(sessionCtxKey, id)
		return next(c)
	}
}

// RequireSession rejects requests that don't carry a valid session cookie,
// i.e. requests that didn't come from a page served by this web UI.
func RequireSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ck, err := c.Cookie(sessionCookieName)
		if err == nil {
			if id, _, ok := parseSessionToken(ck.Value); ok {
				c.Set(sessionCtxKey, id)
				return next(c)
			}
		}
		return c.String(http.StatusForbidden, "Your session has expired. Please reload the page and try again.")
	}
}

// ClientKey identifies a client for rate limiting. Onion/I2P traffic all
// arrives from localhost, so per-IP limits would lump every Tor/I2P user
// together; key those by session instead when one is available.
func ClientKey(c echo.Context) string {
	if isAnonymousNetwork(c) {
		if id, ok := c.Get(sessionCtxKey).(string); ok && id != "" {
			return "s:" + id
		}
	}
	return "ip:" + c.RealIP()
}

// SessionKey identifies the session for per-session rate limiting, falling
// back to the client key when there is no session.
func SessionKey(c echo.Context) string {
	if id, ok := c.Get(sessionCtxKey).(string); ok && id != "" {
		return "s:" + id
	}
	return ClientKey(c)
}
