// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under MIT license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"io/ioutil"
	"net"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config represents the config file
type Config struct {
	// Addresses is the list of IP addresses monitored
	Addresses []net.IP `yaml:"addresses"`
}

func loadConfigFromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open config file")
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse config file")
	}

	return cfg, nil
}
