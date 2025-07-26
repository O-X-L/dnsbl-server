# Changelog

# 1.0.1

## Features

* Disable IP- or Domain-Listeners if no config for them was provided

## Chore

* Go Version 1.23 => 1.24

----

# 1.0.0

## Features

* **Very Simple config**
* **Support for matches**:
  * IPv4
  * IPv6
  * IP Networks (CIDR format)
  * Domains
* **Handling of all basic DNS-BL Requests**
  * IPv4 query (Example: `192.0.2.88 = 88.2.0.192.ip.<domain>`)
  * IPv6 query (Example: `fe80::9fe:dc1c:42f0:6e60 = 0.6.e.6.0.f.2.4.c.1.c.d.e.f.9.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip.<domain>`)
  * Domain query (Example: `malicious.risk.oxl.app = malicious.risk.oxl.app.d.<domain>`)
* **Handling of bad requests**
* **Logging options**
  * Basic
  * Basic without timestamp (if ran as service)
  * JSON-format
  * JSON-format without timestamp
