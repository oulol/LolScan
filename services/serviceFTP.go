package services

import (
	"bufio"
	"net"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type ServiceFTP struct {
	ServiceInterface
	Address string
}

func (s *ServiceFTP) Init(address string) {
	s.Address = address
}

func (s *ServiceFTP) GetAddress() string {
	return s.Address
}

func (s *ServiceFTP) CanIdentify() bool {
	conn, err := net.DialTimeout("tcp", s.Address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	if strings.HasPrefix(line, "220") {
		return true
	}

	return false
}

func (s *ServiceFTP) GetName() string {
	return "FTP server"
}

func (s *ServiceFTP) GetType() ServiceType {
	return ServiceTypeFTP
}

func (s *ServiceFTP) TryLogin(login string, password string) LoginStatus {
	c, err := ftp.Dial(s.Address, ftp.DialWithTimeout(timeout))
	if err != nil {
		return LoginFailed
	}
	defer c.Quit()

	err = c.Login(login, password)
	if err != nil {
		if strings.Contains(err.Error(), "530") {
			return LoginFailed
		}
		return LoginFailed
	}

	return LoginSuccess
}

func (s *ServiceFTP) StoreSnapshots(path string) error {
	return nil
}
