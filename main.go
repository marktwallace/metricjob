package main

import (
	"log"
	"net/http"
	"os"

	"github.com/marktwallace/metricjob/internal/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func printEnv(name string) {
	log.Println(name, "==", os.Getenv(name))
}

func main() {
	printEnv("APP_VERSION")
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting HTTP listener for /metrics at", ":8745")
	go http.ListenAndServe(":8745", nil)

	app := app.NewApp()
	app.Start()
}
