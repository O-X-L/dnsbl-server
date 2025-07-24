package internal

import (
	"fmt"
	"log"
	"net/netip"
	"regexp"
	"slices"
	"strings"

	"github.com/miekg/dns"
)

const RES_RATE = "127.255.255.255"
const BAD_REQ = "-"
const REGEX_DOMAIN = "^[a-z0-9\\-\\.]{1,253}\\.[a-z0-9\\-]{2,63}(\\.?)$"

func logRequest(q string, res string, cli string, c *DNSBLRunningConfig, t int) {
	if !c.Log {
		return
	}
	ts := "IP"
	if t == LOOKUP_DOMAIN {
		ts = "Domain"
	}
	l := fmt.Sprintf("[%s] => %s: %s <= %s", cli, ts, q, res)
	if c.LogTime {
		log.Println(l)
	} else {
		fmt.Println(l)
	}
}

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
				logRequest(query, "400", cli, c, t)

			} else if res != "" {
				logRequest(query, "200", cli, c, t)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, res))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}

			} else {
				logRequest(query, "404", cli, c, t)
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
