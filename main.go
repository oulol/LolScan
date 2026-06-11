package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ips []string
var credentials []string
var ports []string
var brute bool
var identify bool

var discoveryThreads int

var start time.Time
var done int32 = 0

var Version = "development"

func main() {
	initConsole()
	printLogo()

	println("Running LolScan \"" + Version + "\" on " + runtime.GOOS + " " + runtime.GOARCH)

	ipsFlag := flag.String("ips", "ips.txt", "A file that contains target ips")
	credentialsFlag := flag.String("creds", "credentials.txt", "A file that contains credentials to try (login:password)")
	portsFlag := flag.String("ports", "37777,8000,8001,8080,8081", "Comma separated list of ports to scan")
	discoveryThreadsFlag := flag.Int("threads", 32, "The amount of threads to search ports")
	noBruteforceFlag := flag.Bool("nobrute", false, "Disables bruteforce if present.")
	noIdentifyFlag := flag.Bool("noidentify", false, "Disables service identification if present.")

	flag.Parse()
	initDirectory()

	ipsFile := *ipsFlag
	ipsRaw, err := os.ReadFile(ipsFile)
	if err != nil {
		error("Failed to parse IPs file: " + err.Error())
		return
	}
	ipsStr := string(ipsRaw)
	log("Parsing IPs from file " + ipsFile)

	for _, line := range strings.Split(strings.ReplaceAll(ipsStr, "\r", ""), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "/") {
			_, ipNet, err := net.ParseCIDR(line)
			if err != nil {
				warn("Failed to parse subnet " + line + ":" + err.Error())
				continue
			}

			ones, _ := ipNet.Mask.Size()
			if ones < 16 {
				warn("Subnet " + line + " is too huge. Skipped to prevent memory overflow.")
				continue
			}

			for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
				tempIP := make(net.IP, len(ip))
				copy(tempIP, ip)
				ips = append(ips, tempIP.String())
			}
			continue
		}

		ip := net.ParseIP(line)
		if ip == nil {
			warn("Failed to parse IP: " + line)
			continue
		}
		ips = append(ips, ip.String())
	}
	log("Loaded " + fmt.Sprint(len(ips)) + " IPs")

	credentialsFile := *credentialsFlag
	credsRaw, err := os.ReadFile(credentialsFile)
	if err != nil {
		warn("Failed to parse credentials file: " + err.Error() + " Using default ones.")
		credentials = []string{
			"admin:admin",
			"admin:admin123",
			"admin:admin12345",
			"root:root",
			"root:toor",
			"root:admin",
		}
	} else {
		credsStr := string(credsRaw)
		for _, cLine := range strings.Split(strings.ReplaceAll(credsStr, "\r", ""), "\n") {
			cLine = strings.TrimSpace(cLine)
			if cLine != "" {
				credentials = append(credentials, cLine)
			}
		}
	}
	log("Loaded " + fmt.Sprint(len(credentials)) + " credentials")

	for _, p := range strings.Split(*portsFlag, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			ports = append(ports, p)
		}
	}
	log("Searching on " + fmt.Sprint(len(ports)) + " ports")

	discoveryThreads = *discoveryThreadsFlag
	log("Using " + fmt.Sprint(discoveryThreads) + " threads to find open ports")

	brute = !*noBruteforceFlag
	identify = !*noIdentifyFlag

	log("Scan started")
	start = time.Now()

	total := len(ips) * len(ports)
	stopProgress := make(chan struct{})
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				doneVal := atomic.LoadInt32(&done)
				percent := float64(doneVal) / float64(total) * 100

				printBar(doneVal, total, percent, time.Now().Sub(start).Round(time.Second).String())
			case <-stopProgress:
				return
			}
		}
	}()

	var openTargets []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	sem := make(chan struct{}, discoveryThreads)

	for _, ip := range ips {
		for _, port := range ports {
			addr := ip + ":" + port

			wg.Add(1)
			sem <- struct{}{}

			go func(target string) {
				defer wg.Done()
				defer func() { <-sem }()

				conn, err := net.DialTimeout("tcp", target, 700*time.Millisecond)
				if err == nil {
					conn.Close()

					mu.Lock()
					logPortOpen(target)
					openTargets = append(openTargets, target)
					if identify {
						go postOpen(target)
					} else {
						results.WriteString(target + "\n")
					}
					mu.Unlock()
				}

				atomic.AddInt32(&done, 1)
			}(addr)
		}
	}

	wg.Wait()
	close(stopProgress)
	close(sem)

	end := time.Now()
	diff := end.Sub(start)
	log("Finished scanning " + fmt.Sprint(len(ips)) + " ips in " + diff.String() + ". Found " + fmt.Sprint(len(openTargets)) + " open targets.")

	if brute {
		log("Waiting for bruteforcing threads to stop...")
		bruteGroup.Wait()

		end = time.Now()
		diff = end.Sub(start)
		log("Bruteforce finished. Total time: " + diff.String())
	}
}
