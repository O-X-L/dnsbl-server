package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"git.oxl.at/dnsbl-server/src/internal"
	"github.com/miekg/dns"
)

const VERSION = "1.0.0"

func main() {
	var configFile string
	var port int
	var noLog bool
	var noLogTime bool
	var logJSON bool

	flag.StringVar(&configFile, "config", "", "Path to the config file (in YAML format)")
	flag.IntVar(&port, "port", 5353, "Port to listen on")
	flag.BoolVar(&noLog, "no-log", false, "Disable request logging")
	flag.BoolVar(&noLogTime, "no-log-time", false, "Disable log timestamp")
	flag.BoolVar(&logJSON, "log-json", false, "Log in JSON-format")
	flag.Parse()

	fmt.Printf("DNS-BL Server v%v\n  © OXL IT Service\n  License: GPLv3\n\n", VERSION)

	if configFile == "" {
		fmt.Println("ERROR: Config-file is required!")
		os.Exit(1)
	}

	config := internal.DNSBLRunningConfig{
		BL:      internal.DNSBLConfigFlat{},
		Log:     !noLog,
		LogTime: !noLogTime,
		LogJSON: logJSON,
	}

	configRaw := internal.DNSBLConfigFile{}
	err := internal.LoadConfig(configFile, &configRaw)
	if err != nil {
		fmt.Printf("ERROR: Failed to load config-file - %v\n", err)
		os.Exit(1)
	}
	internal.ValidateFlattenConfig(&configRaw, &config)

	if len(config.BL.Domains) == 0 && len(config.BL.IPs) == 0 && len(config.BL.Nets) == 0 {
		fmt.Printf("ERROR: Empty config-file - %v\n", err)
		os.Exit(1)
	}

	log.Printf("DNS-BL server listening on %d\n", port)
	if len(config.BL.IPs) > 0 || len(config.BL.Nets) > 0 {
		dns.HandleFunc(config.BaseIP[1:], config.LookupIP)
		fmt.Printf(" > IP Lookup: %v\n", config.BaseIP[1:])
	}
	if len(config.BL.Domains) > 0 {
		dns.HandleFunc(config.BaseDomain[1:], config.LookupDomain)
		fmt.Printf(" > Domain Lookup: %v\n", config.BaseDomain[1:])
	}

	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	err = server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
