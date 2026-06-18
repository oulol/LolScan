package services

import (
	"bytes"
	"net"

	"golang.org/x/crypto/ssh"
)

type ServiceSSH struct {
	ServiceInterface
	Address string
}

func (s *ServiceSSH) Init(address string) {
	s.Address = address
}

func (s *ServiceSSH) GetAddress() string {
	return s.Address
}

func (s *ServiceSSH) CanIdentify() bool {
	var err error
	conn, err := net.Dial("tcp", s.Address)
	if err != nil {
		return false
	}
	defer conn.Close()

	buf := make([]byte, 3)
	_, err = conn.Read(buf)
	if err != nil {
		return false
	}

	if bytes.Equal(bytes.ToLower(buf), []byte("ssh")) {
		return true
	}

	return false
}

func (s *ServiceSSH) GetName() string {
	return "SSH server"
}

func (s *ServiceSSH) GetType() ServiceType {
	return ServiceTypeSSH
}

func (s *ServiceSSH) TryLogin(login string, password string) LoginStatus {
	config := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	client, err := ssh.Dial("tcp", s.Address, config)
	if err != nil {
		return LoginFailed
	}

	defer client.Close()
	return LoginSuccess
}

func (s *ServiceSSH) StoreSnapshots(path string) error {
	return nil
}
