package internal

import (
	"fmt"
	"net/netip"
	"os"

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
	Domains []DNSBLEntries `yaml:"domains"`
	IPs     []DNSBLEntries `yaml:"ips"`
	Nets    []DNSBLEntries `yaml:"nets"`
}

type DNSBLConfigFlat struct {
	Domains map[string]string
	IPs     map[netip.Addr]string
	Nets    map[netip.Prefix]string
}

type DNSBLRunningConfig struct {
	BL         DNSBLConfigFlat
	BaseIP     string
	BaseDomain string
	Log        bool
	LogTime    bool
	LogJSON    bool
}

func (config *DNSBLRunningConfig) LookupIP(w dns.ResponseWriter, r *dns.Msg) {
	HandleDnsRequest(w, r, config, LOOKUP_IP)
}

func (config *DNSBLRunningConfig) LookupDomain(w dns.ResponseWriter, r *dns.Msg) {
	HandleDnsRequest(w, r, config, LOOKUP_DOMAIN)
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

func FlattenConfig(c *DNSBLConfigFile, f *DNSBLConfigFlat) {
	f.Domains = map[string]string{}
	f.IPs = map[netip.Addr]string{}
	f.Nets = map[netip.Prefix]string{}

	for _, l := range c.Domains {
		for _, e := range l.Content {
			f.Domains[e] = l.Response
		}
	}

	for _, l := range c.IPs {
		for _, e := range l.Content {
			ip, err := netip.ParseAddr(e)
			if err != nil {
				fmt.Printf("Ignoring IP in invalid format: %v (%v)\n", e, ip)
				continue
			}
			f.IPs[ip] = l.Response
		}
	}

	for _, l := range c.Nets {
		for _, e := range l.Content {
			n, err := netip.ParsePrefix(e)
			if err != nil {
				fmt.Printf("Ignoring network in invalid format: %v (%v)\n", e, n)
				continue
			}
			f.Nets[n] = l.Response
		}
	}

}
