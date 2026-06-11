package services

import "fmt"

type ServiceType int

const (
	ServiceTypeCamera ServiceType = iota
	ServiceTypeWeb
	ServiceTypeUnknown
)

var serviceNames = [...]string{
	ServiceTypeCamera:  "Camera",
	ServiceTypeWeb:     "Web",
	ServiceTypeUnknown: "Unknown",
}

func (s ServiceType) String() string {
	if s < 0 || int(s) >= len(serviceNames) {
		return fmt.Sprintf("ServiceType(%d)", s)
	}
	return serviceNames[s]
}

type LoginStatus int

const (
	LoginSuccess LoginStatus = iota
	LoginNotRequired
	LoginBlocked
	LoginFailed
)
