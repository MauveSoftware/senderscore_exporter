// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under MIT license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var scoreDesc = prometheus.NewDesc("senderscore_score", "senderscore.org score of the IP address", []string{"ip", "ptr"}, nil)

type collector struct {
	cfg      *Config
	resolver *net.Resolver
}

func newCollector(cfg *Config) *collector {
	return &collector{
		cfg:      cfg,
		resolver: resolverFromConfig(cfg),
	}
}

func resolverFromConfig(cfg *Config) *net.Resolver {
	if cfg.DNSServer == "" {
		return net.DefaultResolver
	}

	dnsServer := fmt.Sprintf("%s:53", cfg.DNSServer)

	return &net.Resolver{
		PreferGo: true, // Forces the Go resolver to be used
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(5) * time.Second,
			}
			return d.DialContext(ctx, network, dnsServer)
		},
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
			err := c.collectForIP(context.Background(), a, ch)
			if err != nil {
				logrus.Error(err)
			}

			wg.Done()
		}(addr)
	}

	wg.Wait()
}

func (c *collector) collectForIP(ctx context.Context, ip net.IP, ch chan<- prometheus.Metric) error {
	host := reverseIP(ip) + "score.senderscore.com"

	ips, err := c.resolver.LookupIP(ctx, "ip4", host)
	if err != nil {
		return errors.Wrapf(err, "could not get score for %s", ip)
	}

	if len(ips) == 0 {
		return nil
	}

	// Lookup PTR record
	names, err := c.resolver.LookupAddr(ctx, ip.String())
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
