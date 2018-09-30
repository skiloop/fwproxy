[![Go Report Card](https://goreportcard.com/badge/github.com/skiloop/fwproxy)](https://goreportcard.com/report/github.com/skiloop/fwproxy)


fwproxy
====================================================
a proxy in go using [google martian](https://github.com/google/martian) mapping other proxies to local network interface

### Usage

map proxy http://user:Passwd@proxy.com:8012 to local address 127.0.0.1:50000

```bash
fwproxy 127.0.0.1:50000:http://user:Passwd@proxy.com:8012
```