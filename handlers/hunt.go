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
// along with this program, if not, see <http://www.gnu.org/licenses/>.

package handlers

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/CypherGoat/web/views"
)

type HuntEntry struct {
	Prize string `json:"prize"` // "200xmr" | "100xmr" | "50xmr" | "20xmr" | "shirt"
	Label string `json:"label"`
}

type HuntClaim struct {
	Code      string `json:"code"`
	Contact   string `json:"contact"` // XMR address for XMR prizes, contact info for shirts
	ClaimedAt string `json:"claimed_at"`
}

var (
	huntMu         sync.RWMutex
	huntEntries    map[string]HuntEntry
	huntClaimed    map[string]HuntClaim
	huntClaimsFile string
)

// InitHunt loads codes from GOAT_HUNT_FILE and persisted claims from the
// adjacent _claims.json file. Safe to call when the env var is not set.
func InitHunt() {
	huntFile := os.Getenv("GOAT_HUNT_FILE")
	if huntFile == "" {
		return
	}

	if i := strings.LastIndex(huntFile, "."); i >= 0 {
		huntClaimsFile = huntFile[:i] + "_claims.json"
	} else {
		huntClaimsFile = huntFile + "_claims.json"
	}

	data, err := os.ReadFile(huntFile)
	if err != nil {
		return
	}
	// Unmarshal via RawMessage so non-entry fields (e.g. readme comments) are skipped
	// rather than causing the entire parse to fail.
	raw := make(map[string]json.RawMessage)
	if json.Unmarshal(data, &raw) != nil {
		return
	}
	normalised := make(map[string]HuntEntry, len(raw))
	for k, v := range raw {
		var e HuntEntry
		if json.Unmarshal(v, &e) == nil && e.Prize != "" {
			normalised[strings.ToUpper(k)] = e
		}
	}

	claimed := make(map[string]HuntClaim)
	if cdata, err := os.ReadFile(huntClaimsFile); err == nil {
		var list []HuntClaim
		if json.Unmarshal(cdata, &list) == nil {
			for _, c := range list {
				claimed[c.Code] = c
			}
		}
	}

	huntMu.Lock()
	huntEntries = normalised
	huntClaimed = claimed
	huntMu.Unlock()
}

func lookupHuntCode(raw string) (HuntEntry, bool) {
	huntMu.RLock()
	defer huntMu.RUnlock()
	if huntEntries == nil {
		return HuntEntry{}, false
	}
	e, ok := huntEntries[strings.ToUpper(strings.TrimSpace(raw))]
	return e, ok
}

func huntCodeClaimed(raw string) bool {
	huntMu.RLock()
	defer huntMu.RUnlock()
	_, ok := huntClaimed[strings.ToUpper(strings.TrimSpace(raw))]
	return ok
}

// HuntPrizeStatus returns claimed/total counts keyed by prize type.
func HuntPrizeStatus() map[string]views.PrizeStatus {
	huntMu.RLock()
	defer huntMu.RUnlock()
	out := make(map[string]views.PrizeStatus)
	for code, e := range huntEntries {
		s := out[e.Prize]
		s.Total++
		if _, claimed := huntClaimed[code]; claimed {
			s.Claimed++
		}
		out[e.Prize] = s
	}
	return out
}

var errAlreadyClaimed = errors.New("already claimed")

func recordHuntClaim(raw, contact string) error {
	huntMu.Lock()
	defer huntMu.Unlock()
	code := strings.ToUpper(strings.TrimSpace(raw))
	if _, ok := huntClaimed[code]; ok {
		return errAlreadyClaimed
	}
	huntClaimed[code] = HuntClaim{
		Code:      code,
		Contact:   contact,
		ClaimedAt: time.Now().UTC().Format(time.RFC3339),
	}
	var all []HuntClaim
	for _, c := range huntClaimed {
		all = append(all, c)
	}
	data, _ := json.MarshalIndent(all, "", "  ")
	return os.WriteFile(huntClaimsFile, data, 0644)
}
