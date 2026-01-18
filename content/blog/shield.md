---
title: "What Is CypherGoat Shield?"
description: CypherGoat Shield is an execution-risk protection mechanism designed to cover rare swap failures using exchange-provided guarantees — without custody or pooled user funds.
date: 2026-01-11T16:00:00+02:00
image: "/blog/images/shield.png"
draft: false
author: "4rkal"
---

Instant swaps are **trust-minimized**, but they are not risk-free.

CypherGoat routes swaps through independent third-party exchanges. While CypherGoat **never takes custody of user funds**, successful execution still depends on external infrastructure and in extremely rare cases, that infrastructure fails or behaves unfairly.

**CypherGoat Shield exists to address that reality directly.**


## **What CypherGoat Shield Is**

CypherGoat Shield is an **execution-risk protection program**.

It is **not insurance**, **not custody**, and **not a price or rate guarantee**.

Shield allows CypherGoat to compensate users when a swap fails due to exchange-side issues that are **outside the user’s control**, including:

- Exchange infrastructure or backend failures  
- Execution errors after funds are received  
- Exchanges holding funds due to KYC or verification issues **without a valid justification**  

If Shield applies, the user does **not** need to negotiate with the exchange.  
CypherGoat intervenes directly.


## **Which Swaps Are Covered**

**Only exchanges that participate in the Shield program are covered.**

Every exchange listed with a Shield icon has provided CypherGoat with a **financial guarantee**. This guarantee:

- Is paid by the exchange  
- Is held by CypherGoat  
- Exists solely to cover eligible Shield claims  

If an exchange has **not** provided a guarantee, swaps routed through that exchange are **not eligible** for Shield coverage.

Shield eligibility — and the maximum covered amount — is **always shown next to the Shield icon before the user starts a swap**.  
The displayed amount is locked in at transaction creation.


## **What Shield Actually Covers**

Shield may apply when **all** of the following conditions are met:

- The user followed all on-screen instructions  
- Funds were sent to the correct address  
- The correct network, asset, and memo (if required) were used  
- The exchange received the funds  
- The exchange failed to complete the transaction **or** withheld funds without a valid reason  

If an exchange cannot justify withholding funds due to:
- High AML risk  
- A legal or law-enforcement order  

CypherGoat may reimburse the user **up to the covered limit** shown at trade start.


## **What Shield Does *Not* Cover**

Shield does **not** apply to:

- High AML-risk transactions  
- Funds originating from mixers, tumblers, or sanctioned addresses  
- Funds frozen due to police requests, court orders, or legal obligations  
- Incorrect destination addresses  
- Wrong networks or chains  
- Missing, incorrect, or malformed memos  
- Market volatility or slippage  
- Blockchain congestion outside execution guarantees  
- Funds stolen from a user’s wallet  
- Claims reported more than **four (4) weeks** after transaction creation  

Shield is about **unfair or failed execution**, not trading risk or user error.


## **How Shield Claims Work**

If a problem occurs:

1. **Contact CypherGoat support** within four (4) weeks of creating the transaction and provide your transaction ID.  
2. CypherGoat will **mediate with the exchange first** to attempt resolution or completion.  
3. If the exchange fails to resolve the issue and cannot provide a valid justification, CypherGoat will reimburse the user **up to the covered amount**.

The process may take several days or longer depending on complexity.


## **How Shield Is Funded**

Shield is **not funded by users**.

It is funded through exchange guarantees paid to CypherGoat  


User funds are:

- Never pooled  
- Never rehypothecated  
- Never touched  

When Shield is triggered, **CypherGoat absorbs the cost** using exchange guarantees.


## **Why Shield Exists**

Most aggregators forward users to exchanges and disengage when something breaks.

CypherGoat takes a different approach:

- Risks are acknowledged  
- Trade-offs are explicit  
- Failures cost **CypherGoat**, not users  

Shield does not pretend instant swaps are riskless.  
It exists to **align incentives**, **reduce downside**, and **enforce accountability** across the execution stack.

For full terms and conditions, see:  
https://cyphergoat.com/cyphergoat-shield
