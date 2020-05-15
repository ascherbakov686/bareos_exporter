package main

import (
	"github.com/dreyau/bareos_exporter/error"

	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var connectionString string

var (
	exporterPort     = flag.Int("port", 9625, "Bareos exporter port")
	exporterEndpoint = flag.String("endpoint", "/metrics", "Bareos exporter endpoint")
	dbUser        = flag.String("u", "postgres", "Bareos Postgres username")
	dbAuthFile    = flag.String("p", "./auth", "Postgres password file path")
	dbHostname    = flag.String("h", "127.0.0.1", "Postgres hostname")
	dbPort        = flag.String("P", "5432", "Postgres port")
	db            = flag.String("db", "bareos", "Postgres database name")
)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: bareos_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	pass, err := ioutil.ReadFile(*dbAuthFile)
	error.Check(err)
	
        connectionString = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full", *dbUser, strings.TrimSpace(string(pass)), *dbHostname,*db)
	
	collector := bareosCollector()
	prometheus.MustRegister(collector)

	http.Handle(*exporterEndpoint, promhttp.Handler())
	log.Info("Beginning to serve on port ", *exporterPort)

	addr := fmt.Sprintf(":%d", *exporterPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
