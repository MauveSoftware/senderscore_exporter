// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under MIT license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const version string = "0.2.1"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9665", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	configFile    = flag.String("config.path", "config.yml", "Path to config file")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("senderscore_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Copyright: 2020, Mauve Mailorder Software GmbH & Co. KG, Licensed under MIT license")
	fmt.Println("Metric exporter for senderscore.org scores")
}

func startServer() {
	logrus.Infof("Starting Senderscore exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Senderscore Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>Senderscore Exporter by Mauve Mailorder Software</h1>
			<h2>Metrics</h2>
			<p><a href="/metrics">here</a></p>
			<h2>More information</h2>
			<p><a href="https://github.com/MauveSoftware/senderscore_exporter">github.com/MauveSoftware/senderscore_exporter</a></p>
			</body>
			</html>`))
	})

	cfg, err := loadConfigFromFile(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	prometheus.MustRegister(newCollector(cfg))
	http.Handle("/metrics", promhttp.Handler())

	logrus.Infof("Listening for %s on %s", *metricsPath, *listenAddress)
	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}
