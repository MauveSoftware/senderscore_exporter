package main

import (
	"fmt"
	"net"
	"strconv"
)

func reverseIP(ip net.IP) string {
	if ip.To4() != nil {
		return reverseIPv4(ip.To4())
	}

	return reverseIPv6(ip)
}

func reverseIPv4(ip net.IP) string {
	s := ""

	for i := 3; i >= 0; i-- {
		s += strconv.Itoa(int(ip[i])) + "."
	}

	return s
}

func reverseIPv6(ip net.IP) string {
	s := ""

	a := make([]uint8, 16)
	for i, v := range ip {
		a[i] = v
	}
	fmt.Println(a)

	for i := 0; i < 8; i++ {
		idx := 14 - (2 * i)
		s += revserseBlock(ip[idx : idx+2])
	}

	return s
}

func revserseBlock(b []uint8) string {
	hex := fmt.Sprintf("%02x%02x", b[0], b[1])

	s := ""
	for i := 3; i >= 0; i-- {
		s += string(hex[i]) + "."
	}

	return s
}
