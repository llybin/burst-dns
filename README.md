# BurstCoin DNS

Run under root, because software needs access to resolv.conf and listen 53 port.

Example for test: http://devtrue.burst/

# How to add domain

Open your wallet on any node, e.g.: https://wallet.burst.devtrue.net

![Alias section](https://i.imgur.com/OJ6W27f.png)

![Register alias](https://i.imgur.com/ZKDz1EX.png)

Supported types: a, aaaa

Exmaple: `[{"type": "a", "data": "IP1", "ttl": 3600}, {"type": "a", "data": "IP2", "ttl": 7200}]`
