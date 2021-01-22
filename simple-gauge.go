package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ec2_available_ip = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "aws",
			Name:      "ec2_subnet_available_ip_total",
			Help:      "Available IP number for the subnet",
		})
)

func main() {
	rand.Seed(time.Now().Unix())

	prometheus.MustRegister(ec2_available_ip)

	go func() {
		for {
			ec2_available_ip.Set(rand.Float64()*15 - 5)
			time.Sleep(time.Second)
		}
	}()
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(":8080", nil))
}