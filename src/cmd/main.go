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

const VERSION = "1.0.0"

func main() {
	var domain string
	var configFile string
	var port int
	var noLog bool
	var noLogTime bool
	var logJSON bool

	flag.StringVar(&domain, "domain", "", "Domain to serve for")
	flag.StringVar(&configFile, "config", "", "Path to the config file (in YAML format)")
	flag.IntVar(&port, "port", 5353, "Port to listen on")
	flag.BoolVar(&noLog, "no-log", false, "Disable request logging")
	flag.BoolVar(&noLogTime, "no-log-time", false, "Disable log timestamp")
	flag.BoolVar(&logJSON, "log-json", false, "Log in JSON-format")
	flag.Parse()

	fmt.Printf("DNS-BL Server v%v\n  © OXL IT Service\n  License: GPLv3\n\n", VERSION)

	if domain == "" || configFile == "" {
		fmt.Println("ERROR: Domain and config-file need to be provided!")
		os.Exit(1)
	}

	validDomain, _ := regexp.MatchString(internal.REGEX_DOMAIN, domain)
	if !validDomain {
		fmt.Println("ERROR: Invalid domain provided! Example: 'dnsbl.risk.oxl.app'")
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
		LogJSON:    logJSON,
		BaseIP:     fmt.Sprintf(".%v", baseIP),
		BaseDomain: fmt.Sprintf(".%v", baseDomain),
	}

	configRaw := internal.DNSBLConfigFile{}
	err := internal.LoadConfig(configFile, &configRaw)
	if err != nil {
        fmt.Printf("ERROR: Failed to load config-file - %v\n", err)
		os.Exit(1)
	}
	internal.FlattenConfig(&configRaw, &config.BL)

	dns.HandleFunc(baseIP, config.LookupIP)
	dns.HandleFunc(baseDomain, config.LookupDomain)

	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("DNS-BL server listening on %d\n > IP Lookup: %v\n > Domain Lookup: %v", port, baseIP, baseDomain)
	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
