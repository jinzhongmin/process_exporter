package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Export = newExport("process.yml")

func init() {
	prometheus.MustRegister(Export)
}

func main() {
	http.Handle(Export.config.Global.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(Export.config.Global.Port, nil))
}
