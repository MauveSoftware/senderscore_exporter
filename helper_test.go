package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseIP(t *testing.T) {
	tests := []struct {
		name     string
		cidr     string
		expected string
	}{
		{
			name:     "ipv4",
			cidr:     "185.138.53.66/32",
			expected: "66.53.138.185.",
		},
		{
			name:     "ipv6",
			cidr:     "2a07:a40::/128",
			expected: "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.4.a.0.7.0.a.2.",
		},
	}

	for _, test := range tests {
		ip, _, _ := net.ParseCIDR(test.cidr)
		assert.Equal(t, test.expected, reverseIP(ip))
	}
}
