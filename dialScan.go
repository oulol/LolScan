package main

import (
	"net"
	"sync"
	"sync/atomic"
)

func dialScan() int64 {
	var opens int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	jobs := make(chan string, discoveryThreads*2)

	for i := 0; i < discoveryThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range jobs {
				conn, err := net.DialTimeout("tcp", target, timeout)
				if err == nil {
					conn.Close()

					mu.Lock()
					logPortOpen(target)
					opens++
					if identify {
						go postOpen(target)
					} else {
						results.WriteString(target + "\n")
					}
					mu.Unlock()
				}
				atomic.AddInt32(&done, 1)
			}
		}()
	}

	for _, ip := range ips {
		strIp := ip.String()
		for _, port := range ports {
			jobs <- strIp + ":" + port
		}
	}
	close(jobs)

	wg.Wait()
	return opens
}
