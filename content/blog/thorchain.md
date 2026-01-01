---
title: "THORChain Swaps Are Now Available on CypherGoat"
description: CypherGoat now supports THORChain-based swaps as an optional routing path. Here’s what that actually means, how it works, and when it makes sense to use it.
date: 2025-12-30T21:00:00+02:00
image: "/blog/images/thor.png"
draft: false
author: "4rkal"
---

We’ve added **THORChain-based swap routing** to CypherGoat.

This is **not a redesign**, **not a pivot**, and **not a recommendation** that users should prefer THORChain over other routes. It is simply an additional execution option that is now available in cases where it makes sense.

Given the mixed (and often justified) opinions around THORChain in the Monero community, this post is intentionally factual and practical.

## **What Was Added**

CypherGoat now supports routing swaps through THORChain **alongside** existing centralized and off-chain partners.

From a technical standpoint, this means:
- Some swaps can now be executed **entirely on-chain**
- Certain pairs no longer require an intermediary exchange
- Users can choose a route that settles natively on the destination chain

Nothing else about CypherGoat changes.

## **How THORChain Swaps Actually Work**

A THORChain swap is not a “swap” in the traditional exchange sense.

Instead:
1. You send funds to a protocol-controlled vault  
2. Validators observe and confirm the inbound transaction  
3. Liquidity pools are used to convert assets  
4. Funds are released to the destination address  

There is **no custody account**, **no order book**, and **no KYC** — but there *is* on-chain observability and protocol-level logic deciding execution.

This is fundamentally different from how most CypherGoat routes work today.

## **Where This Makes Sense**

THORChain routing can be useful when:
- Centralized liquidity is temporarily unavailable
- Users explicitly want an on-chain execution path
- A native BTC → LTC or BTC → ETH route is preferred
- Off-chain providers are rate-limited or congested

In those cases, THORChain can function as a **fallback liquidity source**, not a replacement.

## **Where It Doesn’t**

THORChain is **not** always the best option.

Known trade-offs include:
- Higher slippage during volatility
- Longer execution times
- Fully traceable transaction flows
- Different privacy properties than off-chain swaps

For many users, especially privacy-maximalists, centralized or hybrid routes may still be preferable.


## **Transparency by Design**

All THORChain swaps are:
- Fully on-chain
- Publicly verifiable
- Trackable via transaction hashes and memos

CypherGoat does not abstract this away or attempt to make it appear “private”.

If you use this route, you know exactly what you’re using.

## **Why We Added It Anyway**

CypherGoat’s role is not to judge which infrastructure is ideologically pure.

Our role is to:
- Aggregate execution paths
- Expose trade-offs honestly
- Let users decide what they trust

Ignoring THORChain entirely would mean **withholding a real execution option**, even if that option isn’t universally liked.

Providing it transparently is the more honest approach.

## **Current Scope**

At launch, THORChain routing is limited to:
- Specific BTC-based routes
- Assets supported natively by the protocol
- Pairs where execution reliability meets our internal thresholds

We will expand or restrict availability based on real-world performance, not narratives.

## **Closing Thoughts**

THORChain is neither a silver bullet nor a threat to CypherGoat’s philosophy.

It’s a tool sometimes useful, sometimes not.

By integrating it as **one route among many**, we give users more control without pretending that every trade-off disappears.

As always:
**understand the route you choose, keep custody of your keys, and don’t trust abstractions you can’t verify.**
