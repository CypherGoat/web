---
title: "One Interface, JavaScript Optional"
description: CypherGoat now runs as a single adaptive interface that works fully without JavaScript while offering optional enhancements when JavaScript is enabled.
date: 2026-01-18T12:00:00+02:00
image: "/blog/images/nojs.jpg"
draft: false
author: "4rkal"
---

CypherGoat has always treated JavaScript as **optional**, not mandatory.

Until now, this was implemented through two separate interfaces:
- `cyphergoat.com` for the JavaScript-enhanced experience  
- `nojs.cyphergoat.com` for a fully JavaScript-free interface  

While functional, this split created unnecessary fragmentation.

That architecture has now changed.


## **A Single Adaptive Interface**

CypherGoat now operates as **one unified interface** under a single domain.

The behavior is straightforward:

- If **JavaScript is enabled**, CypherGoat progressively enables optional features that improve usability.
- If **JavaScript is disabled**, CypherGoat automatically serves a fully functional, JavaScript-free interface.

No redirects.  
No broken flows.  
No reduced functionality.


## **What JavaScript Adds (When Enabled)**

When JavaScript is available, CypherGoat may offer:

- A **more usable searchable coin selection dropdown**, making it easier to browse and select assets  
- A **shorter swap flow**, reducing the number of pages required to create a swap  

These features exist purely for convenience.

They are **not required** to create or complete a swap.


## **What Happens With JavaScript Disabled**

With JavaScript turned off:

- All core swap functionality remains available  
- Pages are rendered server-side  
- Standard HTML forms are used  
- No client-side logic is required  
- No scripts fail silently or degrade security  

The JavaScript-free path is not a fallback or compatibility mode.  
It is a **first-class interface**.


## **Why This Change Was Made**

Many CypherGoat users operate under strict threat models:

- JavaScript disabled by default  
- Aggressive script whitelisting  
- Hardened browser profiles  
- Tor Browser or similar environments  

Requiring a separate domain for these users introduced friction that didn’t need to exist.

A privacy-respecting service should adapt to the user — not the other way around.

## **What This Means in Practice**

- One canonical domain  
- One interface  
- No forced trust assumptions  
- No loss of functionality when JavaScript is disabled  

Existing bookmarks continue to work.  
`nojs.cyphergoat.com` remains valid.  
`cyphergoat.com` now adapts automatically.


## **Open Source by Design**

The JavaScript-optional architecture is not based on trust alone.

CypherGoat’s frontend and core interface logic are open source, allowing anyone to inspect how JavaScript is used, what it does, and just as importantly what it does *not* do.

This makes it possible to verify that:
- JavaScript is optional  
- Core swap logic does not depend on client-side execution  
- Disabling JavaScript does not introduce hidden behavior  

The goal is not transparency as a slogan, but verifiability in practice.

## **The Philosophy Behind the Change**

This update is not about aesthetics or performance metrics.

It is about aligning the interface with CypherGoat’s core principles:

- User control over execution environment  
- Minimal and explicit trust assumptions  
- Compatibility with adversarial threat models  
- No artificial barriers for privacy-conscious users  

CypherGoat should work the way *you* choose to use it.

This change makes that explicit.
