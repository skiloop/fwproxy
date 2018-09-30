fwproxy
====================================================
a proxy in go mapping other proxies to local network interface

### Usage

map proxy http://user:Passwd@proxy.com:8012 to local address 127.0.0.1:50000

```bash
fwproxy 127.0.0.1:50000:http://user:Passwd@proxy.com:8012
```