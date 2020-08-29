package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func mainMetrics(cmd *string) {
	for {
		jails := getJails(cmd)
		for _, jail := range jails {
			out, err := exec.Command(*cmd, "status", jail).Output()
			if err != nil {
				log.Fatal(err)
			}
			// fail2ban_banned_current
			fail2banBannedCurrent.WithLabelValues(jail).Set(getFloat64Val(out, `(Currently banned:\s+)(.+)`))
			// fail2ban_banned_total
			fail2banBannedTotal.WithLabelValues(jail).Set(getFloat64Val(out, `(Total banned:\s+)(.+)`))
			// fail2ban_failed_current
			fail2banFailedCurrent.WithLabelValues(jail).Set(getFloat64Val(out, `(Currently failed:\s+)(.+)`))
			// fail2ban_failed_total
			fail2banFailedTotal.WithLabelValues(jail).Set(getFloat64Val(out, `(Total failed:\s+)(.+)`))
		}
		time.Sleep(10 * time.Second)
	}
}

func getIPList(cmd *string, jail *string) []string {
	out, err := exec.Command(*cmd, "status", *jail).Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	return re.FindAllString(string(out), -1)
}

func geoIPLocator(cmd *string, engine string) {
	var (
		sec   time.Duration = 60
		ipGrp map[string]int
	)
	for {
		jails := getJails(cmd)
		for _, jail := range jails {
			ipList := getIPList(cmd, &jail)
			switch engine {
			case "ip2c":
				// Amount of requests per user / per day is unlimited, just be reasonable.
				// Currently ip2c.org can sustain a maximum of about 30 million per day.
				ipGrp = IP2c(cmd, &ipList)
				fmt.Println(ipGrp)
				sec = 600
			case "geoiplookup":
				ipGrp = GeoIPLookup(cmd, &ipList)
			case "freegeoip":
				// Limit 15000 requests per hour
				ipGrp = FreeGeoIP(cmd, &ipList)
				sec = 900
			}
			for key, val := range ipGrp {
				fail2banHackersLocations.WithLabelValues(jail, key).Set(float64(val))
			}
		}
		time.Sleep(sec * time.Second)
	}
}

// IP2c engine
func IP2c(cmd *string, ipList *[]string) map[string]int {
	ipGrp := make(map[string]int)
	for _, ip := range *ipList {
		resp, err := http.Get(fmt.Sprintf("http://ip2c.org/?ip=%s", ip))
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		data := strings.Split(string(body), ";")
		res, err := strconv.Atoi(data[0])
		if err != nil {
			log.Fatalln(err)
		}

		switch res {
		case 0: // WRONG INPUT
		case 2: // UNKNOWN
		default:
			ipGrp[data[1]] = ipGrp[data[1]] + 1
		}
	}
	return ipGrp
}

// GeoIPLookup engine
func GeoIPLookup(cmd *string, ipList *[]string) map[string]int {
	// Find geoiplookup executable
	geoipLookup, err := exec.LookPath("geoiplookup")
	if err != nil {
		log.Fatal(err)
	}
	ipGrp := make(map[string]int)
	for _, ip := range *ipList {
		geoipInfo, err := exec.Command(geoipLookup, ip).Output()
		if err != nil {
			log.Fatal(err)
		}
		re := regexp.MustCompile(`\s[A-Z]{2}`)
		country := strings.Trim(re.FindString(string(geoipInfo)), " ")
		ipGrp[country] = ipGrp[country] + 1
	}
	return ipGrp
}

// FreeGeoIP ...
func FreeGeoIP(cmd *string, ipList *[]string) map[string]int {
	//
	return make(map[string]int)
}
