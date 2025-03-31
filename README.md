# 🐐 CypherGoat - instant exchange aggregator

This is the official repository of [cyphergoat.com](https://cyphergoat.com) the open source instant exchange aggregator. This repository houses the official web-ui (the js free UI will be released soon)

All of the code is licensed under AGPL-v3

## Installation
Requirements:
- [Go](https://go.dev)
- [Air](https://github.com/air-verse/air) for hot reloading
- [Templ](https://github.com/a-h/templ) for templates

Clone the repository

`git clone https://github.com/CypherGoat/web`

`cd web`

Install dependencies

`go mod tidy`

Setup your environment:

`export CYPHERGOAT_API_KEY="your_api_key_here"`

`export CYPHERGOAT_API_URL="api.cyphergoat.com"` (optional)

OR via .env

`echo "CYPHERGOAT_API_KEY=your_api_key_here" > .env`

> **Note:** The API key is only required for live API calls. You can still audit, review and modify the code without an API key.

Start the development server

`air`

If you head to http://localhost:4200 you should see cyphergoat up and running


## API intergration
CypherGoat requires an API key for fetching exchange rates and executing trades. The API provides endpoints for
- Fetching estimates
- Creating transactions
- Tracking transaction status

See our API documentation at [api.cyphergoat.com](https://api.cyphergoat.com)