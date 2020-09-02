package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func parseArgs() (bool, int, string) {
	debugPtr := flag.Bool("debug", false, "Debug. Disable search fail2ban-client")
	verPtr := flag.Bool("version", false, "Show version.")
	portPtr := flag.Int("port", 9635, "Listen port number.")
	enginePtr := flag.String("engine", "ip2c",
		fmt.Sprintf(
			"%s\n"+
				"%s\n"+
				"%-12s %s\n"+
				"%-12s %s\n"+
				"%-12s %s\n",
			"GeoIP engine.",
			"Supported values:",
			"ip2c", " - get info from ip2c.org. Unlimited, just be reasonable. Currently ip2c.org can sustain a maximum of about 30 million per day.",
			"geoiplookup", " - use geoiplookup console util (fastest method. Linux only)",
			"freegeoip", " - get info from https://freegeoip.app",
		))
	flag.Parse()

	if *verPtr == true {
		fmt.Println("Fail2ban exporter. v0.20.831\nCopyright © 2020 Dennis Y. Parygin (dyp2000@mail.ru)")
		os.Exit(0)
	}

	switch *enginePtr {
	case "ip2c", "geoiplookup", "freegeoip":
		//
	default:
		log.Fatal("Engine [", *enginePtr, "] not supported.")
		os.Exit(1)
	}

	return *debugPtr, *portPtr, *enginePtr
}

func main() {
	debug, port, engine := parseArgs()

	// Find fail2ban-client executable
	if !debug {
		f2b, err := exec.LookPath("fail2ban-client")
		if err != nil {
			log.Fatal(err)
		}
		// Main metrics
		go mainMetrics(&f2b)
		// GeoIP metrica
		go geoIPLocator(&f2b, engine)
	}

	apiSrv := newServer(port)
	apiSrv.start()
}
