---
title: "SafeRoute Explained: Why Splitting Your Swap Can Get You a Better Rate"
description: "SafeRoute is CypherGoat's algorithm that automatically splits large swaps across multiple exchanges when it improves your rate. Here's how it works and why it matters."
date: 2026-06-24T16:00:00+02:00
image: "/blog/images/saferoute.jpg"
draft: false
author: "4rkal"
---

# SafeRoute Explained: Why Splitting Your Swap Can Get You a Better Rate

If you've swapped a large amount on CypherGoat recently, you may have noticed a "SafeRoute" badge appear on your quote. This post explains what it is, why it exists, and what's actually happening behind the scenes when it shows up.

## The problem: big orders move the price

Every exchange has a limited amount of liquidity sitting at the best available rate. Small swaps barely touch it. Large swaps don't.

Once your order size gets big enough, you start eating into worse and worse liquidity on the order book. The rate you're quoted at the top isn't the rate you actually get filled at — the bigger the order, the bigger the gap. This is called **slippage**, and it grows non-linearly: a 10x bigger order doesn't cost 10x more slippage, it often costs much more.

## The fix: split the order

If no single exchange can absorb your full order without slippage, the obvious move is to not give the full order to a single exchange.

SafeRoute checks, on every quote, whether splitting your swap across two CypherGoat partner exchanges produces a better total payout than routing it through one. Each leg only has to absorb half the size, so each leg sits closer to the top of its own order book. Add the two outputs together, and the combined total can beat what a single exchange would have given you.

## How it works for you

You don't configure anything. The flow is:

1. **You enter an amount**, exactly like any normal swap.
2. **SafeRoute scores every exchange pair combination** in the background and checks if a split would materially beat the best single-exchange quote.
3. **If it wins, you get two orders instead of one** — two deposit addresses, two amounts — and both execute at the same time. The combined output lands in the same destination wallet you specified.

If splitting doesn't help, you never see it. SafeRoute is gated on confidence: it only surfaces when the improvement is real, not rounding noise.

## Why this only matters above a certain size

Slippage is negligible on small orders, so there's nothing for SafeRoute to fix below a certain threshold. Based on data from over 2,700 real swaps, SafeRoute's gains become consistently meaningful above roughly **1.2 BTC** in order size. Below that, single-exchange execution is usually already optimal, and SafeRoute correctly stays out of the way.

Above that threshold, the average improvement measured across those swaps was **+0.067–0.083%** in total receive amount — small per swap, but it compounds, and critically: **zero negative outcomes** were recorded. Splitting never made the trade worse; it either helped or did nothing.

## Your funds are still covered

A split swap still only touches CypherGoat's vetted exchange partners — never random third parties. That means **CG Shield** coverage applies to each leg independently. If one partner in a split mishandles or freezes funds, Shield's collateral mechanism kicks in on that leg exactly as it would on a normal single-exchange swap.

The two legs also have no dependency on each other. One slow exchange can't hold up the other.

## The talk

We presented SafeRoute's design and the underlying data at MoneroKon 2026 in Warsaw — "Liquidity Fragmentation In Non Custodial XMR Markets" If you want the full technical breakdown, [watch the talk](https://cyphergoat.com/saferoute).

## TL;DR

- Large swaps suffer slippage on any single exchange.
- SafeRoute automatically checks if splitting across two CypherGoat partners beats a single exchange.
- It only shows up when the gain is real — typically above ~1.2 BTC.
- Both legs are Shield-eligible and execute simultaneously.
- No setup, no account, no extra steps. It's just a better quote when one exists.

Try it yourself — [get a quote on CypherGoat](https://cyphergoat.com/) and watch for the badge.
