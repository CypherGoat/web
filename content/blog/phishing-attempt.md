---
title: "A Phishing Attempt Against CypherGoat and How We Shut It Down in 24 Hours"
description: How CypherGoat identified a phishing clone early, neutralized it instantly and kept users fully protected.
date: 2025-12-04T22:00:00+03:00
image: "/blog/images/phishing.png"
draft: false
author: "4rkal"
---

On **16 November 2025**, we detected something unusual in our traffic logs: a single visit referral coming from a domain that looked suspiciously close to ours **cyhpergoat.com**.

At first glance it appeared harmless, but one letter was out of place. Instead of the standard “ph,” the attacker swapped the characters to **“hp”**, a classic typo-squatting technique designed to trick users into visiting a fake version of a site they trust.

A quick inspection revealed that the entire website was an exact clone of CypherGoat, down to the last pixel. The domain had been registered just one day earlier, clearly with the intention of impersonation. Thanks to our monitoring systems, we caught it before it ever became functional.


## **Immediate Response**

Once we verified the clone, we took action right away:
- We contacted **Cloudflare**, which was providing DNS services for the phishing domain.
- We contacted **Nicenic**, the domain registrar responsible for issuing the domain.
Both were provided with documentation of the impersonation and an explanation of the potential threat to users.

Throughout the entire incident, no user ever interacted with the phishing domain.

## **Breaking the Clone**
While investigating, we realized something important:  
**The phishing site had no backend and no ability to perform swaps. It was simply forwarding all traffic to us**

This meant users were never at risk. But it also meant we could shut down their attempt instantly.


We modified our backend to detect requests originating from that specific host. Instead of serving our usual interface, we responded to any traffic coming from _cyhpergoat.com_ with a simple message:

`Access denied. Nice try, counterfeit goat.`

![](/blog/images/cyhpergoat.jpg)


Their entire operation collapsed on the spot.
No functionality, no proxying, no future exploit.

## **Resolution**
Within a couple of hours of our initial report, Cloudflare and Nicenic confirmed that the domain was under review.

A day later just **48 hours after the fake domain had been registered** it was **fully taken down**!

**The phishing clone never advanced beyond its harmless forwarding stage**.
## **Why It Matters**

Phishing remains one of the most common attack vectors in crypto.

CypherGoat is built on transparency and trust, so we treat impersonation attempts seriously. We encourage all users to:

- Bookmark **[https://cyphergoat.com](https://cyphergoat.com)**
- Double-check URLs before initiating swaps

A few seconds of caution can prevent major loss - but rest assured, CypherGoat is watching for these threats long before they ever reach our users.
## **Closing Thoughts**

The attack was unsophisticated, but it was a reminder that the more CypherGoat grows, the more we’ll attract this kind of behaviour.

We’ll continue monitoring for impersonation attempts and coordinating with infrastructure providers to ensure malicious actors can’t compromise our users.

And if anyone tries to spin up another counterfeit goat?

Let’s just say:  
**we’re ready.**