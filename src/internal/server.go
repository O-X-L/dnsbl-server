package internal

import (
	"fmt"
	"net/netip"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const RES_RATE = "127.255.255.255"
const BAD_REQ = "-"
const REGEX_DOMAIN = "^[a-z0-9\\-\\.]{1,253}\\.[a-z0-9\\-]{2,63}(\\.?)$"

var STARTUP_TIME = time.Now().Format("2006010215")

func checkIP(q dns.Question, c *DNSBLRunningConfig) (string, string) {
	req := strings.Replace(q.Name, c.BaseIP, "", 1)
	parts := strings.Split(req, ".")
	slices.Reverse(parts)
	var ip netip.Addr
	var err error
	if len(parts) == 32 {
		ip6Parts := []string{}
		s := 0
		for i := 0; i < 8; i++ {
			s = i * 4
			ip6Parts = append(ip6Parts, strings.Join(parts[s:s+4], ""))
		}
		ip6 := strings.Join(ip6Parts, ":")
		ip, err = netip.ParseAddr(ip6)
		if err != nil {
			return BAD_REQ, req
		}
	} else if len(parts) == 4 {
		ip4 := strings.Join(parts, ".")
		ip, err = netip.ParseAddr(ip4)
		if err != nil {
			return BAD_REQ, req
		}
	} else {
		return BAD_REQ, req
	}

	res, found := c.BL.IPs[ip]
	if found {
		return res, ip.String()
	}

	for n, r := range c.BL.Nets {
		if n.Contains(ip) {
			return r, ip.String()
		}
	}

	return "", ip.String()
}

func checkDomain(q dns.Question, c *DNSBLRunningConfig) (string, string) {
	domain := strings.Replace(q.Name, c.BaseDomain, "", 1)
	validDomain, _ := regexp.MatchString(REGEX_DOMAIN, domain)
	if !validDomain {
		return BAD_REQ, domain
	}

	res, found := c.BL.Domains[domain]
	if !found {
		return "", domain
	}
	return res, domain
}

func soaResponse(d string, c *DNSBLRunningConfig) string {
	return fmt.Sprintf("%s 3600 IN SOA %s %s %s 10800 3600 604800 3600", d, c.NS[0], c.AdminMail, STARTUP_TIME)
}

func getBaseDomain(t int, c *DNSBLRunningConfig) string {
	if t == LOOKUP_IP {
		return c.BaseIP[1:]
	} else {
		return c.BaseDomain[1:]
	}
}

func parseQuery(m *dns.Msg, w dns.ResponseWriter, c *DNSBLRunningConfig, t int) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			var res string
			cli := ""
			query := ""
			if c.Log {
				cli = strings.Split(w.RemoteAddr().String(), ":")[0]
			}

			if t == LOOKUP_IP {
				res, query = checkIP(q, c)
			} else {
				res, query = checkDomain(q, c)
			}

			if res == BAD_REQ {
				logRequest(query, 400, cli, c, t, "")
				m.Rcode = dns.RcodeRefused

			} else if res != "" {
				logRequest(query, 200, cli, c, t, res)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, res))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				m.Rcode = dns.RcodeSuccess

			} else {
				logRequest(query, 404, cli, c, t, "")
				m.Rcode = dns.RcodeNameError
			}

			baseDomain := getBaseDomain(t, c)
			rr, err := dns.NewRR(soaResponse(baseDomain, c))

			if err == nil {
				m.Ns = append(m.Ns, rr)
			}

		case dns.TypeSOA:
			baseDomain := getBaseDomain(t, c)
			rr, err := dns.NewRR(soaResponse(baseDomain, c))

			if err == nil {
				if q.Name == baseDomain {
					m.Answer = append(m.Answer, rr)
				} else {
					m.Ns = append(m.Ns, rr)
				}
				m.Rcode = dns.RcodeSuccess

			} else {
				fmt.Println("SOA ERROR:", q.Name, baseDomain, err)
				m.Rcode = dns.RcodeServerFailure
			}

		case dns.TypeNS:
			baseDomain := getBaseDomain(t, c)

			for _, host := range c.NS {
				rr, err := dns.NewRR(fmt.Sprintf("%s 3600 IN NS %s", baseDomain, host))
				if err == nil {
					if q.Name == baseDomain {
						m.Answer = append(m.Answer, rr)
					} else {
						m.Ns = append(m.Ns, rr)
					}
					m.Rcode = dns.RcodeSuccess

				} else {
					fmt.Println("NS ERROR:", q.Name, baseDomain, err)
					m.Rcode = dns.RcodeServerFailure
				}
			}
		}
	}
}

func HandleDnsRequest(w dns.ResponseWriter, r *dns.Msg, c *DNSBLRunningConfig, l int) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, w, c, l)
	}

	w.WriteMsg(m)
}
