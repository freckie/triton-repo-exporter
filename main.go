package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	var tritonAddr, port string
	flag.StringVar(&tritonAddr, "tritonAddr", "", "address of the triton server (host:port)")
	flag.StringVar(&port, "port", "9100", "listener port")
	flag.Parse()

	if tritonAddr == "" {
		log.Errorf("--tritonAddr not given")
		os.Exit(1)
	}
	log.Infof("tritonAddr=%s", tritonAddr)

	hostname := os.Getenv("NODENAME")
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			log.Errorf("Failed to get hostname: %s", err)
			os.Exit(1)
		}
	}
	log.Infof("hostname=%s", hostname)

	prometheus.MustRegister(&RepositoryCollector{
		hostname:   hostname,
		tritonAddr: tritonAddr,
	})
	log.Info("Registered a new RepositoryCollector successfully")

	http.Handle("/metrics", promhttp.Handler())

	log.Infof("Trying to listen at 0.0.0.0:%s", port)
	if err = http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
		os.Exit(1)
	}
}
