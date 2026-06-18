package main

import (
	"LolScan/services"
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

var scanner *IpScanner
var credentials []string
var ports []string
var types []services.ServiceType
var brute bool
var identify bool

var discoveryThreads int

var start time.Time
var done int32 = 0

var timeout time.Duration

var Version = "development"

func main() {
	defer func() {
		if r := recover(); r != nil {
			logErr("Panic in main: " + fmt.Sprint(r))
			log("Stack:\n" + string(debug.Stack()))
		}
	}()
	initConsole()
	printLogo()

	println("Running LolScan \"" + Version + "\" on " + runtime.GOOS + " " + runtime.GOARCH)

	ipsFlag := flag.String("ips", "ips.txt", "A file that contains target ips")
	credentialsFlag := flag.String("creds", "credentials.txt", "A file that contains credentials to try (login:password)")
	portsFlag := flag.String("ports", "37777,8000,8001,8080,8081", "Comma separated list of ports to scan")
	discoveryThreadsFlag := flag.Int("threads", 32, "The amount of threads to search ports")
	noBruteforceFlag := flag.Bool("nobrute", false, "Disables bruteforce if present.")
	noIdentifyFlag := flag.Bool("noidentify", false, "Disables service identification if present.")
	typesFlag := flag.String("types", "all", "Scans for only specified types (web,camera,ssh,ftp). Set to all for every type.")
	timeoutFlag := flag.Int("timeout", 700, "Timeout in ms")

	flag.Parse()
	initDirectory()

	ipsFile := *ipsFlag
	ipsRaw, err := os.Open(ipsFile)
	if err != nil {
		logErr("Failed to parse IPs file: " + err.Error())
		return
	}
	log("Parsing IPs from file " + ipsFile)

	scanner = NewIpScanner(bufio.NewScanner(ipsRaw))
	ipsAmount, err := countIPsInFile(ipsFile)
	if err != nil {
		logErr("Error counting IPs: " + err.Error())
		return
	}

	log("Scanning " + fmt.Sprint(ipsAmount) + " IPs")

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

	timeoutInt := *timeoutFlag
	timeout = time.Duration(timeoutInt) * time.Millisecond
	services.SetTimeout(timeout)

	for _, p := range strings.Split(*typesFlag, ",") {
		p = strings.TrimSpace(p)
		for val, str := range services.ServiceNames {
			if strings.EqualFold(str, p) {
				types = append(types, services.ServiceType(val))
			}
		}
	}

	log("Scan started")
	start = time.Now()

	total := ipsAmount * len(ports)
	stopProgress := make(chan struct{})
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				doneVal := atomic.LoadInt32(&done)
				percent := float64(doneVal) / float64(total) * 100

				printBar(doneVal, total, percent, time.Since(start).Round(time.Second).String())
			case <-stopProgress:
				return
			}
		}
	}()

	opens := dialScan()
	close(stopProgress)

	end := time.Now()
	diff := end.Sub(start)
	log("Finished scanning " + fmt.Sprint(ipsAmount) + " ips in " + diff.String() + ". Found " + fmt.Sprint(opens) + " open targets.")

	if brute {
		log("Waiting for processing threads to stop...")
		bruteGroup.Wait()

		end = time.Now()
		diff = end.Sub(start)
		log("Processing finished. Total time: " + diff.String())
	}
}
