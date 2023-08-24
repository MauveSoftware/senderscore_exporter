// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under MIT license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	scoreDesc = prometheus.NewDesc("senderscore_score", "senderscore.org score of the IP address", []string{"ip"}, nil)
)

type collector struct {
	cfg *Config
}

func newCollector(cfg *Config) *collector {
	return &collector{
		cfg: cfg,
	}
}

// Describe implements prometheus.Collector interface
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scoreDesc
}

// Collect implements prometheus.Collector interface
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	wg := &sync.WaitGroup{}
	wg.Add(len(c.cfg.Addresses))

	for _, addr := range c.cfg.Addresses {
		go func(a net.IP) {
			err := collectForIP(a, ch)
			if err != nil {
				logrus.Error(err)
			}

			wg.Done()
		}(addr)
	}

	wg.Wait()
}

func collectForIP(ip net.IP, ch chan<- prometheus.Metric) error {
	host := reverseIP(ip) + "score.senderscore.com"

	ips, err := net.LookupIP(host)
	if err != nil {
			return errors.Wrapf(err, "could not get score for %s", ip)
	}

	if len(ips) == 0 {
			return nil
	}

	// Lookup PTR record
	names, err := net.LookupAddr(ip.String())
	ptr := ""
	if err == nil && len(names) > 0 {
			ptr = names[0]
	}

	resolvedIP := ips[0].To4()
	score := resolvedIP[3]

	// Always add IP as label, add PTR record as additional label if it resolves
	ch <- prometheus.MustNewConstMetric(scoreDesc, prometheus.GaugeValue, float64(score), ip.String(), ptr)
	return nil
}
