package services

import (
	"fmt"
	"time"
)

type ServiceType int

const (
	ServiceTypeCamera ServiceType = iota
	ServiceTypeWeb
	ServiceTypeSSH
	ServiceTypeFTP
	ServiceTypeUnknown
)

var ServiceNames = [...]string{
	ServiceTypeCamera:  "Camera",
	ServiceTypeWeb:     "Web",
	ServiceTypeSSH:     "SSH",
	ServiceTypeFTP:     "FTP",
	ServiceTypeUnknown: "Unknown",
}

func (s ServiceType) String() string {
	if s < 0 || int(s) >= len(ServiceNames) {
		return fmt.Sprintf("ServiceType(%d)", s)
	}
	return ServiceNames[s]
}

type LoginStatus int

const (
	LoginSuccess LoginStatus = iota
	LoginNotRequired
	LoginBlocked
	LoginFailed
)

var timeout = 1000 * time.Millisecond

func SetTimeout(val time.Duration) {
	timeout = val
}
