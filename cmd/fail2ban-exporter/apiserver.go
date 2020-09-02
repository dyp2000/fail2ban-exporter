package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	port int
}

// newServer ...
func newServer(port int) *server {
	return &server{
		port: port,
	}
}

func (s *server) rootHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, r.URL.Path+"\n")
	fmt.Fprintf(w, r.Method+"\n")

	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Fprint(w, "<html>\n"+
			"<head>\n"+
			"<title>Fail2ban exporter</title>\n"+
			"<meta http-equiv=\"Cache-Control\" content=\"no-cache, must-revalidate, max-age=0\">\n"+
			"</head>\n"+
			"<body>\n"+
			"<h1>Fail2ban exporter</h1>\n"+
			"<p><a href='/metrics'>Metrics</a></p>\n"+
			"</body>\n"+
			"</html>\n",
		)
	case "POST":
		switch r.URL.Path {
		case "/search":

		case "/query":

		case "/annotation":

		case "/tag-keys":

		case "/tag-values":

		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// Start ...
func (s *server) start() {
	http.HandleFunc("/", s.rootHandler)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil); err != nil {
		log.Fatal(err)
	}

}
