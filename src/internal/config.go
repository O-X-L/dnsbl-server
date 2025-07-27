package internal

import (
	"fmt"
	"net/netip"
	"os"
	"regexp"
	"strings"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

const LOOKUP_IP = 1
const LOOKUP_DOMAIN = 2

type DNSBLEntries struct {
	Response string   `yaml:"response"`
	Content  []string `yaml:"content"`
}

type DNSBLConfigFile struct {
	Domain    string         `yaml:"domain"`
	AdminMail string         `yaml:"admin_mail"`
	NS        []string       `yaml:"nameservers"`
	Domains   []DNSBLEntries `yaml:"domains"`
	IPs       []DNSBLEntries `yaml:"ips"`
	Nets      []DNSBLEntries `yaml:"nets"`
}

type DNSBLConfigFlat struct {
	Domains map[string]string
	IPs     map[netip.Addr]string
	Nets    map[netip.Prefix]string
}

type DNSBLRunningConfig struct {
	BL         DNSBLConfigFlat
	Root       string
	BaseIP     string
	BaseDomain string
	NS         []string
	AdminMail  string
	Log        bool
	LogTime    bool
	LogJSON    bool
}

func (config *DNSBLRunningConfig) LookupIP(w dns.ResponseWriter, r *dns.Msg) {
	HandleDnsBLRequest(w, r, config, LOOKUP_IP)
}

func (config *DNSBLRunningConfig) LookupDomain(w dns.ResponseWriter, r *dns.Msg) {
	HandleDnsBLRequest(w, r, config, LOOKUP_DOMAIN)
}

func (config *DNSBLRunningConfig) LookupRoot(w dns.ResponseWriter, r *dns.Msg) {
	HandleDNSRootDomain(w, r, config)
}

func LoadConfig(config_file string, d *DNSBLConfigFile) error {
	file, err := os.ReadFile(config_file)
	if err != nil {
		return fmt.Errorf("config file does not exist %v: %v", config_file, err)
	}
	err = yaml.Unmarshal(file, d)
	if err != nil {
		return fmt.Errorf("config file could not be parsed %v: %v", config_file, err)
	}
	return nil
}

func ValidateFlattenConfig(c *DNSBLConfigFile, r *DNSBLRunningConfig) {
	if c.Domain == "" || c.AdminMail == "" || len(c.NS) == 0 {
		fmt.Println("ERROR: Domain, Admin-Email-Address and Nameserver(s) are required!")
		os.Exit(1)
	}

	validDomain, _ := regexp.MatchString(REGEX_DOMAIN, c.Domain)
	if !validDomain {
		fmt.Println("ERROR: Invalid domain provided! Example: 'dnsbl.risk.oxl.app'")
		os.Exit(1)
	}
	if !strings.HasSuffix(c.Domain, ".") {
		c.Domain += "."
	}

	r.Root = c.Domain
	r.BaseIP = fmt.Sprintf(".ip.%v", c.Domain)
	r.BaseDomain = fmt.Sprintf(".d.%v", c.Domain)

	for _, hostname := range c.NS {
		if strings.HasSuffix(hostname, ".") {
			hostname = strings.TrimSuffix(hostname, ".")
		}
		validDomain, _ := regexp.MatchString(REGEX_DOMAIN, hostname)
		if !validDomain {
			fmt.Printf("ERROR: Invalid nameserver hostname provided: %s\n", hostname)
			os.Exit(1)
		}
	}
	r.NS = c.NS
	r.AdminMail = strings.ReplaceAll(c.AdminMail, "@", ".")
	if strings.HasSuffix(r.AdminMail, ".") {
		r.AdminMail = strings.TrimSuffix(r.AdminMail, ".")
	}
	if !validDomain || !strings.Contains(c.AdminMail, "@") {
		fmt.Printf("ERROR: Invalid Admin-Email-Address provided: %s\n", c.AdminMail)
		os.Exit(1)
	}

	r.BL.Domains = map[string]string{}
	r.BL.IPs = map[netip.Addr]string{}
	r.BL.Nets = map[netip.Prefix]string{}

	for _, l := range c.Domains {
		for _, e := range l.Content {
			r.BL.Domains[strings.ToLower(e)] = l.Response
		}
	}

	for _, l := range c.IPs {
		for _, e := range l.Content {
			ip, err := netip.ParseAddr(e)
			if err != nil {
				fmt.Printf("Ignoring IP in invalid format: %v (%v)\n", e, ip)
				continue
			}
			r.BL.IPs[ip] = l.Response
		}
	}

	for _, l := range c.Nets {
		for _, e := range l.Content {
			n, err := netip.ParsePrefix(e)
			if err != nil {
				fmt.Printf("Ignoring network in invalid format: %v (%v)\n", e, n)
				continue
			}
			r.BL.Nets[n] = l.Response
		}
	}
}
