---
title: How to use the Monero CLI wallet
description: In this guide I will be showing you how to use the monero cli wallet
date: 2025-04-05T16:19:36+03:00
image: "/blog/images/cli.jpg"
draft: false
author: "4rkal"
---


### Requirements
1. A computer (Windows, Mac or Linux)
2. Some basic command line comfort

In this guide I will be showing you how to use it on Linux mostly, since that is where most people use it.

## Download and install
You can download the cli from the official website [getmonero.org/downloads](https://web.getmonero.org/downloads/)
Choose the one for your operating system. 

After it has downloaded extract it

On Linux simply run: 

`tar -xvjf monero-linux*.tar.bz2`

Now that it has been extracted simply cd into it

`cd monero*`


## Using

### 1. Creating a New Wallet
The first thing we will have to do is create a new wallet. Simply run 

`./monero-wallet-cli`

You should now see something like this:

![](/blog/images/cli1.png)


Enter a wallet file name. I will enter main for example. After that confirm the creation of that wallet. 

After that you will be prompted to enter a new password for the wallet. Use something secure, but remember it. 

![](/blog/images/cli2.png)

Confirm the password. After that you should select the language of your seed. Just choose English for simplicity.

Now you will see something like this:

![](/blog/images/cli3.png)

It is extremely important that you BACKUP this seedphrase without it you will loose access to your funds. 

After that select whether or not you want to mine (most probably no).


After that you can enter welcome (type welcome and hit enter) in order to get some more info

![](/blog/images/cli4.png)


### Basic commands
Initially, unless you are running a node on your computer you will have to use a [remote node](https://www.getmonero.org/resources/moneropedia/remote-node.html). To do that you will have to run this command

```
./monero-wallet-cli --daemon-address NODEADDRESS
```

You can find a full list of monero nodes at [monero.fail](https://monero.fail)

For example
```
./monero-wallet-cli --daemon-address node.sethforprivacy.com:18089
```

After running the command you will see something like this
<img src="/blog/images/cli5.png" alt="Image">

It will show you your full balance etc.


We will now be going over the most essential commands
#### Balance
In order to get your balance all you have to do is run

`balance`

You should see something like this:

```
[wallet 43LSEb]: balance
Currently selected account: [0] Primary account
Tag: (No tag assigned)
Balance: 0.000000000000, unlocked balance: 0.000000000000
```

#### Address
In order to display the primary receiving address run

`address`

You should see something like this:

```
[wallet 43LSEb]: address
0  43LSEbXqZ2hjhrogZsSZ9TKnDSKLxueyY2kUyUwJfHh2L5bJjX3KLeWdD4A6XEsJ8QNqspfmDk5qVdj5QmDPyAPvEmbuZ7W  Primary address 
```

In order to generate a new subaddress you can run

`address new`

Output:

```
[wallet 43LSEb]: address new
1  85TT84Ma6yTXzrcz5yfFpBWYgf4guCzdHE1ip4jwomxV6KV5BB2qqyyQtfPXn2YWThC4YrRVUDrgxj3niZrKtQa377J6odA  (Untitled address) 
```

#### Sending

In order to send monero we use the `transfer` command

The syntax is 

`transfer <address> <amount>`

For example
`transfer 888tNkZrPN6JsEgekjMnABU4TBzc2Dt29EPAvkRxbANsAnjyPbb3iQ1YBRk1UXcdRsiKc9dhwMVgN5S9cQUiyoogDavup3H 1`

This will transfer 1 XMR to the address of the [CCS](https://ccs.getmonero.org/).


#### Refresh
In order to see new incoming transactions you can run

`refresh`

Output:

```
[wallet 43LSEb]: refresh
Starting refresh...
Refresh done, blocks received: 0                                
Currently selected account: [0] Primary account
Tag: (No tag assigned)
Balance: 0.000000000000, unlocked balance: 0.000000000000
```

#### Incoming transactions

In order to see incoming transactions we use the `incoming_transfers` command

To see incoming transactions that are available to spend you may use `incoming_transfers available`

#### Save & Exit

When exiting the wallet first enter `save` in order to not loose any recent transactions or changes

Then type `exit` to safely exit

## Conclusion
And that's about it. The monero wallet cli is an extremely well developed solution suitable for advanced users. While it is not for everyone if you enjoy using the command line, it's amazing!

This article was written by [4rkal](https://4rkal.com) part of the [CypherGoat.com](https://cyphergoat.com) open source instant exchange aggregator team!

