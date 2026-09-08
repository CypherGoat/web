---
title: "Kraken Puts Monthly Caps on Monero Deposits and Withdrawals"
description: "Kraken is rolling out rolling 30-day transaction limits on XMR from August 15th, gated behind a source-of-funds review for anyone who wants more. Here's what it means and how to route around it."
image: "/blog/images/kraken.jpg"
date: 2026-07-16
author: 4rkal
---

Kraken has started sending notices to XMR traders announcing new monthly limits on Monero deposits and withdrawals. The changes take effect **August 15th** and apply on a rolling 30-day basis. The exchange frames them as part of its "risk-based compliance controls" while it continues to support XMR trading and transfers.

In plain terms: Monero on Kraken now comes with a ceiling, and lifting that ceiling means opening your books.

## The new limits

| Account type | Deposit limit (30-day) | Withdrawal limit (30-day) |
|---|---|---|
| Individual | $3,500 | $1,500 |
| Corporate | $30,000 | $10,000 |

Deposits and withdrawals are calculated independently, so hitting one cap doesn't touch the other. The limits reset on a sliding 30-day window rather than resetting on the 1st of each month — the same rolling mechanism Kraken already uses for its general funding limits.

For a supposedly "supported" asset, an individual withdrawal cap of $1,500 a month is a striking number. At current prices that's a fraction of a single XMR position. Anyone using Kraken as an actual on/off-ramp for Monero rather than a place to park a small balance will hit the wall almost immediately.

## The real story is the whitelisting review

The limits themselves are the headline, but the mechanism for exceeding them is where the compliance creep lives. Clients who want higher limits can apply for an exemption through a **whitelisting review process**. As part of that review, Kraken says it may request:

- The purpose of your XMR activity
- Source of funds
- Expected transaction volumes
- Business or trading use case, where applicable

This is source-of-funds gating dressed as a courtesy. The default posture becomes: your Monero activity is capped unless you can justify it to a compliance desk. That's a meaningful shift from "verified account, trade freely" to "verified account, but privacy-coin flow is throttled until you explain yourself." Kraken recommends contacting support and starting the review *before* the effective date if you expect to exceed the caps.

## This isn't happening in a vacuum

The limits are the latest step in a long, steady retreat from Monero across regulated venues:

- **EEA (Oct 2024):** Kraken delisted XMR entirely, halting trading and deposits, with withdrawals cut off at the end of 2024 and residual balances force-converted to BTC.
- **Ireland & Belgium (2024):** Delisted ahead of the broader EEA action, following Binance and OKX out the door.
- **Canada (Sept 2025):** XMR delisted for Canadian clients, moved to withdrawals-only.
- **UAE (2026):** XMR included in a scheduled delisting cycle alongside ZEC, DASH and several stablecoins.

Read together, the pattern is clear. Where regulation forces Kraken's hand, it delists outright. Where it still lists XMR, it is now narrowing the pipe and adding friction. Limits are what a delisting looks like before it's a delisting.

## What it means if you use Kraken for XMR

If your Monero activity is likely to exceed the caps, you have three options, in rough order of how much they cost you:

1. **Submit to the whitelisting review.** You keep using Kraken, but you hand over source-of-funds documentation and a stated purpose for your XMR activity. You're trading privacy for a higher number.
2. **Stay under the cap and treat Kraken as a small-balance venue.** Workable if you only ever move modest amounts, but you're one policy update away from another squeeze.
3. **Stop routing Monero through a KYC exchange at all.** This is the option that actually addresses the root problem rather than the symptom.

## The takeaway

The point of Monero is to not need permission. Every additional layer of caps, reviews and source-of-funds questions on a custodial exchange is a reminder that "supported on a major exchange" and "usable as private money" are drifting further apart.

If you want to move XMR without a monthly allowance and a compliance interview attached, the tooling exists and doesn't require an account: **atomic swaps**, non-custodial swap aggregators, and DEX-style routing let you go in and out of Monero without depositing to a third party that can cap, freeze, or question the flow. Non-custodial routes ([CypherGoat](https://cyphergoat.com) among them) don't hold your funds and don't gate withdrawals behind a review desk — there's nothing to whitelist because there's no custodian in the middle.

Kraken's caps aren't the end of Monero access. They're just one more nudge toward the answer the ecosystem has been giving all along: hold your own keys, and route around anyone who won't let you use them.

---

*Effective date and figures per Kraken's client notice (August 15th). Deposit and withdrawal caps apply where Kraken still supports XMR trading; the asset has already been delisted in the EEA, Canada, the UAE and elsewhere. Not financial advice.*