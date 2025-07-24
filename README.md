# DNS-BL Microservice

## Config

```yaml
domains:
  - response: 127.0.0.2
    content:
      - 'malicious.risk.oxl.app'

ips:
  - response: 127.0.0.2
    content:
      - '192.0.2.88'
      - 'fe80::9fe:dc1c:42f0:6e60'

nets:
  - response: 127.0.0.2
    content:
      - '192.0.2.128/29'
```


## Usage

Users can query the DNS-BL as configured in your config-file through:
* `ip.<DOMAIN>` => for IP-Lookups
* `d.<DOMAIN>` => for Domain-Lookups

```bash
rath@gate:~ dnsbl-server -help
> Usage of dnsbl-server:
>   -config string
>         Path to the config file (in YAML format) (required)
>   -domain string
>         Domain to serve for (required)
>   -log-json
>         Log in JSON-format (defaut false)
>   -no-log
>         Disable request logging (defaut false)
>   -no-log-time
>         Disable log timestamp (defaut false)
>   -port int
>         Port to listen on (default 5353)

rath@gate:~ dnsbl-server -domain test.at -config ./config.yml -port 10000

2025/07/24 21:46:12 DNS-BL server listening on 10000
 > IP Lookup: ip.test.at.
 > Domain Lookup: d.test.at.
# <time> [<client-IP>] => <IP/DOMAIN>: <request> <= <status> <response>
#   200 = found, 400 = bad request, 404 = not found
2025/07/24 21:46:16 [127.0.0.1] => IP: 192.0.2.88 <= 200 127.0.0.2
2025/07/24 21:46:18 [127.0.0.1] => IP: 192.0.2.130 <= 200 127.0.0.2
2025/07/24 21:46:23 [127.0.0.1] => IP: 1.1.1.1 <= 404
2025/07/24 21:46:53 [127.0.0.1] => IP: fe80::9fe:dc1c:42f0:6e60 <= 200 127.0.0.2
2025/07/24 21:48:08 [127.0.0.1] => Domain: malicious.risk.oxl.app <= 200 127.0.0.2
2025/07/24 21:47:42 [127.0.0.1] => Domain: good.oxl.app <= 404

# examples of bad requests
2025/07/24 21:46:49 [127.0.0.1] => IP: 1 <= 400  # bad IP
2025/07/24 21:46:59 [127.0.0.1] => IP: 0.6.e.6.0.f.2.4.c.1.c.d.e.f.9.0.0.0.0.0.0.0.0.0.0.8.e.f <= 400  # bad IPv6
2025/07/24 21:48:13 [127.0.0.1] => IP: malicious.risk.oxl.app <= 400  # domain on IP-lookup
2025/07/24 21:48:16 [127.0.0.1] => Domain: 1.1.1.1 <= 400  # IP on domain-lookup

```

**Client**:

<details>

```
nslookup 
> set port=10000
> server 127.0.0.1
Default server: 127.0.0.1
Address: 127.0.0.1#10000

# IPv4 MATCH:
> 88.2.0.192.ip.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
Name:   88.2.0.192.ip.test.at
Address: 127.0.0.2


# IPv4 NETWORK MATCH:
> 130.2.0.192.ip.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
Name:   130.2.0.192.ip.test.at
Address: 127.0.0.2


# IPv6 MATCH:
> 0.6.e.6.0.f.2.4.c.1.c.d.e.f.9.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
Name:   0.6.e.6.0.f.2.4.c.1.c.d.e.f.9.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip.test.at
Address: 127.0.0.2


# DOMAIN MATCH
> malicious.risk.oxl.app.d.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
Name:   malicious.risk.oxl.app.d.test.at
Address: 127.0.0.2


# IP NOT LISTED:
> 1.1.1.1.ip.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
*** Can't find 1.1.1.1.ip.test.at: No answer


# DOMAIN NOT LISTED
> good.oxl.app.d.test.at
Server:         127.0.0.1
Address:        127.0.0.1#10000

Non-authoritative answer:
*** Can't find good.oxl.app.d.test.at: No answer
```
</details>
