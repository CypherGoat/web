package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Cap (capjs.js.org) proof-of-work captcha.
//
//	CAP_ENDPOINT  public widget URL including the site key, e.g. https://cap.example.com/<site-key>/
//	CAP_SECRET    the site key's secret (server-side only)
//
// Tokens are verified against <CAP_ENDPOINT>siteverify. If either variable is
// missing, verification fails closed so signup is never silently unprotected.

var capHTTP = &http.Client{Timeout: 5 * time.Second}

// CapEndpoint returns the widget endpoint with a trailing slash, or "" if unset.
func CapEndpoint() string {
	ep := strings.TrimSpace(os.Getenv("CAP_ENDPOINT"))
	if ep == "" {
		return ""
	}
	return strings.TrimRight(ep, "/") + "/"
}

// verifyCap checks a solved-captcha token with the Cap server.
func verifyCap(ctx context.Context, token string) bool {
	endpoint := CapEndpoint()
	secret := os.Getenv("CAP_SECRET")
	if endpoint == "" || secret == "" {
		log.Printf("captcha: CAP_ENDPOINT/CAP_SECRET not configured; rejecting")
		return false
	}
	if token == "" || len(token) > 2048 {
		return false
	}

	payload, _ := json.Marshal(map[string]string{"secret": secret, "response": token})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"siteverify", bytes.NewReader(payload))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := capHTTP.Do(req)
	if err != nil {
		log.Printf("captcha: siteverify request failed: %v", err)
		return false
	}
	defer resp.Body.Close()

	var out struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&out); err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK && out.Success
}
