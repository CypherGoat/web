---
title: "Swap Crypto Directly Inside Claude"
description: CypherGoat now works inside Claude. Check live rates, compare exchanges, and initiate swaps without opening a browser.
date: 2026-04-15T12:00:00+02:00
image: "/blog/images/mcp.jpg"
draft: false
author: "4rkal"
---

CypherGoat now has an MCP server.

This means you can check live swap rates, compare exchanges, and initiate swaps directly inside Claude — without opening a browser.

## **What MCP Is**

MCP (Model Context Protocol) is an open standard that lets AI models call external tools natively.

Instead of describing what you want and waiting for a link, the AI calls the tool directly and returns the result inline.

It is supported by Claude Desktop, Claude.ai, and any MCP-compatible client.

## **What the CypherGoat MCP Server Does**

The server exposes CypherGoat's core functionality as callable tools:

- `get_best_rate` — returns the top exchange and estimated output for a given pair and amount
- `get_all_rates` — returns the full ranked list across all partnered exchanges
- `create_swap` — initiates a swap and returns the deposit address
- `get_swap_status` — returns the current status of a swap
- `list_supported_pairs` — returns all available trading pairs

No API key required. No signup.

## **How to Add It**

**Claude.ai:** Go to Settings → Connectors → Add custom connector → paste `https://mcp.cyphergoat.com/mcp`, and connect.

**Claude Desktop:** Add the following to your config file and restart:

```json
"cyphergoat": {
  "type": "url",
  "url": "https://mcp.cyphergoat.com/mcp"
}
```

You can watch a full demo including setup guide here: [Swap Crypto Directly in Claude AI (CypherGoat MCP)](https://www.youtube.com/watch?v=9sRwtqcKBMc)

## **Who This Is For**

The MCP server is most useful for developers building crypto-focused AI agents, or anyone who wants to query rates and initiate swaps without leaving their AI setup.

The use cases are niche but if you're already building with Claude, it fills a gap nothing else does.

The server is live at [mcp.cyphergoat.com](https://mcp.cyphergoat.com).
