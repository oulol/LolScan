package main

import (
	"bufio"
	"net"
	"strings"
)

type IpScanner struct {
	scanner    *bufio.Scanner
	currentIP  net.IP
	currentNet *net.IPNet
}

func NewIpScanner(scanner *bufio.Scanner) *IpScanner {
	return &IpScanner{scanner: scanner}
}

func (s *IpScanner) Next() net.IP {
	for {
		if s.currentNet != nil {
			ipCopy := make(net.IP, len(s.currentIP))
			copy(ipCopy, s.currentIP)

			incIP(s.currentIP)

			if !s.currentNet.Contains(s.currentIP) {
				s.currentNet = nil
				s.currentIP = nil
			}
			return ipCopy
		}

		if !s.scanner.Scan() {
			return nil
		}
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "/") {
			_, ipNet, err := net.ParseCIDR(line)
			if err != nil {
				continue
			}
			s.currentNet = ipNet
			s.currentIP = ipNet.IP.Mask(ipNet.Mask)
			continue
		}

		ip := net.ParseIP(line)
		if ip == nil {
			continue
		}
		return ip
	}
}

func (s *IpScanner) Err() error {
	return s.scanner.Err()
}
