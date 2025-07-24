package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"git.oxl.at/dnsbl-server/src/internal"
	"github.com/miekg/dns"
)

func main() {
	var domain string
	var configFile string
	var port int
	var noLog bool
	var noLogTime bool

	flag.StringVar(&domain, "domain", "", "Domain to serve for")
	flag.StringVar(&configFile, "config", "", "Path to the config file (in YAML format)")
	flag.IntVar(&port, "port", 5353, "Port to listen on")
	flag.BoolVar(&noLog, "no-log", false, "Disable request logging")
	flag.BoolVar(&noLogTime, "no-log-time", false, "Disable log timestamp")
	flag.Parse()

	if domain == "" || configFile == "" {
		fmt.Println("Domain and config-file need to be provided!")
		os.Exit(1)
	}

	validDomain, _ := regexp.MatchString(internal.REGEX_DOMAIN, domain)
	if !validDomain {
		fmt.Println("Invalid domain provided! Example: 'dnsbl.risk.oxl.app'")
		os.Exit(1)
	}
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	baseIP := fmt.Sprintf("ip.%v", domain)
	baseDomain := fmt.Sprintf("d.%v", domain)
	config := internal.DNSBLRunningConfig{
		BL:         internal.DNSBLConfigFlat{},
		Log:        !noLog,
		LogTime:    !noLogTime,
		BaseIP:     fmt.Sprintf(".%v", baseIP),
		BaseDomain: fmt.Sprintf(".%v", baseDomain),
	}

	configRaw := internal.DNSBLConfigFile{}
	internal.LoadConfig(configFile, &configRaw)
	internal.FlattenConfig(&configRaw, &config.BL)

	dns.HandleFunc(baseIP, config.LookupIP)
	dns.HandleFunc(baseDomain, config.LookupDomain)

	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("DNS-BL server listening on %d\n > IP Lookup: %v\n > Domain Lookup: %v", port, baseIP, baseDomain)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
