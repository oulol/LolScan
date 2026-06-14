package main

import (
	"LolScan/services"
	"slices"
	"strings"
	"sync"
)

var bruteGroup sync.WaitGroup

func postOpen(address string) {
	defer bruteGroup.Done()
	bruteGroup.Add(1)
	device := services.Identify(address)
	if device == nil {
		return
	}

	if !slices.Contains(types, device.GetType()) {
		return
	}

	log("Target " + address + " identified (" + device.GetName() + ")")

	if brute {
		for _, cred := range credentials {
			splat := strings.Split(cred, ":")
			login, password := splat[0], splat[1]
			status := device.TryLogin(login, password)
			if status == services.LoginNotRequired {
				logConnected(address)
				addResult(device, "")
				break
			} else if status == services.LoginSuccess {
				logCredentialsFound(address, cred)
				addResult(device, cred)
				break
			} else if status == services.LoginBlocked {
				warn("Target " + address + " is now locked.")
				break
			}
		}
	} else {
		logConnected(address)
		addResult(device, "")
	}
}
