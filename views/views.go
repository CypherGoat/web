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

package views

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type Coin struct {
	Ticker  string `json:"ticker"`
	Name    string `json:"name"`
	Network string `json:"network"`
	Icon    string `json:"icon"`
}

func getCoins() []Coin {
	data, err := ioutil.ReadFile("static/coins.json")
	if err != nil {
		return nil
	}
	var coins []Coin
	json.Unmarshal(data, &coins)
	return coins
}

func formatCoinValue(coin Coin) string {
	if coin.Ticker == coin.Network {
		return coin.Ticker
	}
	return coin.Ticker + "(" + coin.Network + ")"
}

func formatCoinDisplay(coin Coin) string {
	if coin.Network == coin.Ticker {
		return coin.Name + " - " + strings.ToUpper(coin.Ticker)
	}
	return coin.Name + " - " + strings.ToUpper(coin.Ticker) + " (" + strings.ToUpper(coin.Network) + ")"
}
