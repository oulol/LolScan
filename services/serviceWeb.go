package services

import (
	"crypto/tls"
	"net/http"
)

type ServiceWeb struct {
	ServiceInterface
	Address string
	Proto   string
	Server  string
	Status  string
}

func (s *ServiceWeb) Init(address string) {
	s.Address = address
}

func (s *ServiceWeb) GetAddress() string {
	return s.Address
}

func (s *ServiceWeb) CanIdentify() bool {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	r, err := client.Head("https://" + s.Address)
	if err == nil {
		s.Proto = r.Proto + " + https"
		s.Status = r.Status
		s.Server = r.Header.Get("Server")
		return true
	}

	r, err = client.Get("http://" + s.Address)
	if err == nil {
		s.Proto = r.Proto
		s.Status = r.Status
		s.Server = r.Header.Get("Server")
		return true
	}

	return false
}

func (s *ServiceWeb) GetName() string {
	return s.Proto + " " + s.Server + " server - " + s.Status
}

func (s *ServiceWeb) GetType() ServiceType {
	return ServiceTypeWeb
}

func (s *ServiceWeb) TryLogin(login string, password string) LoginStatus {
	return LoginNotRequired
}

func (s *ServiceWeb) StoreSnapshots(path string) error {
	return nil
}
