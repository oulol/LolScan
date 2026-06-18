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
	sem := make(chan struct{}, discoveryThreads)

	for _, ip := range ips {
		for _, port := range ports {
			addr := ip + ":" + port
			wg.Add(1)
			sem <- struct{}{}

			go func(target string) {
				defer wg.Done()
				defer func() { <-sem }()

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
			}(addr)
		}
	}

	wg.Wait()
	close(sem)
	return opens
}
