package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fail2banBannedCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fail2ban_banned_current",
			Help: "Number of currently banned IP addresses.",
		},
		[]string{"jail"},
	)

	fail2banBannedTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fail2ban_banned_total",
			Help: "Number of total banned IP addresses.",
		},
		[]string{"jail"},
	)

	fail2banFailedCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fail2ban_failed_current",
			Help: "Number of current failed connections.",
		},
		[]string{"jail"},
	)

	fail2banFailedTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fail2ban_failed_total",
			Help: "Number of total failed total.",
		},
		[]string{"jail"},
	)

	fail2banHackersLocations = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fail2ban_hackers_locations",
			Help: "Location of hackers on world map.",
		},
		[]string{"jail", "location"},
	)
)

func getFloat64Val(src []byte, reStr string) float64 {
	re := regexp.MustCompile(reStr)
	match := re.FindStringSubmatch(string(src))
	val, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		val = -1.0
	}
	return val
}

// Get Jails list
func getJails(cmd *string) []string {
	out, err := exec.Command(*cmd, "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`(Jail list:\s+)(.+)`)
	match := re.FindStringSubmatch(string(out))
	var jails []string = strings.Fields(match[2])
	return jails
}

func parseArgs() (int, string) {
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
			"freegeoip", " - get info from https://freegeoip.app. Not released.",
		))
	flag.Parse()

	if *verPtr == true {
		fmt.Println("Fail2ban exporter. Ver.0.20.827")
		os.Exit(0)
	}

	switch *enginePtr {
	case "ip2c", "geoiplookup":
		//
	default:
		log.Fatal("Engine [", *enginePtr, "] temporary not supported.")
		os.Exit(1)
	}

	return *portPtr, *enginePtr
}

func main() {
	port, engine := parseArgs()

	// Find fail2ban-client executable
	f2b, err := exec.LookPath("fail2ban-client")
	if err != nil {
		log.Fatal(err)
	}

	// Main metrics
	go mainMetrics(&f2b)

	// GeoIP metrica
	go geoIPLocator(&f2b, engine)

	// -------- HTTP SERVER -------- //
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html>\n"+
			"<head>\n"+
			"<title>Fail2ban exporter</title>\n"+
			"<meta http-equiv=\"Cache-Control\" content=\"no-cache, must-revalidate, max-age=0\">\n"+
			"</head>\n"+
			"<body>\n"+
			"<h1>Fail2ban exporter</h1>\n"+
			"<p><a href='/metrics'>Metrics</a></p>\n"+
			"</body>\n"+
			"</html>\n")
	})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
