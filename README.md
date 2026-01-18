# 🐐 CypherGoat - instant exchange aggregator

This is the official repository of [cyphergoat.com](https://cyphergoat.com) the open source instant exchange aggregator. This repository houses the official web-ui (the js free UI can be found [here](https://github.com/CypherGoat/nojs))

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

## Add support for a coin
To add support for a new coin, it must first be available on at least one of our partnered exchanges. You can view the list of supported exchanges [here](https://cyphergoat.com/#our-partners).

Once confirmed, follow these steps:

1. **Edit the `coins.json` File**  
   Add your coin to the [`coins.json`](https://github.com/CypherGoat/web/blob/main/static/coins.json) file. Use the following format as a template:


```json
   {
     "ticker": "xmr",
     "name": "Monero",
     "network": "xmr",
     "isStable": false,
     "icon": "icons/xmr.webp",
     "min": 0
   }
```

> Note: you should leave the min as 0 and it will be completed by our automated system.

2. **Add the Coin Icon**

Upload the coin’s logo (preferably from coingecko) to the [icons folder](https://github.com/CypherGoat/web/tree/main/static/icons)

## API intergration
CypherGoat requires an API key for fetching exchange rates and executing trades. The API provides endpoints for
- Fetching estimates
- Creating transactions
- Tracking transaction status

See our API documentation at [api.cyphergoat.com](https://api.cyphergoat.com)

You can request an api key sending us an email to [support@cyphergoat.com](mailto:support@cyphergoat.com)
