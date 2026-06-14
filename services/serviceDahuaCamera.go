package services

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type ServiceDahuaCamera struct {
	ServiceInterface
	Address string
	SN      string
	Random  int64
	Conn    net.Conn
}

func (s *ServiceDahuaCamera) Init(address string) {
	s.Address = address
}

func (s *ServiceDahuaCamera) GetAddress() string {
	return s.Address
}

func (s *ServiceDahuaCamera) CanIdentify() bool {
	var err error
	s.Conn, err = net.Dial("tcp", s.Address)
	if err != nil {
		return false
	}

	s.Conn.Write([]byte{0xA0, 0x05, 0x00, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x02, 0x00, 0x01, 0x00, 0x00, 0xA1, 0xAA})

	buf := make([]byte, 1)
	s.Conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err = s.Conn.Read(buf)
	if err != nil {
		s.Conn.Close()
		return false
	}

	if buf[0] == 0xB0 {
		s.Conn.Read(make([]byte, 31))
		l := string(s.ReadLine())
		if strings.HasPrefix(l, "Realm:Login to ") {
			s.SN, _ = strings.CutPrefix(l, "Realm:Login to ")
		}
		r, _ := strings.CutPrefix(string(s.ReadLine()), "Random:")
		s.Random, _ = strconv.ParseInt(r, 10, 64)
		return true
	}

	if s.Conn != nil {
		s.Conn.Close()
	}
	return false
}

func (s *ServiceDahuaCamera) GetName() string {
	return "Dahua Camera, SN=" + s.SN
}

func (s *ServiceDahuaCamera) GetType() ServiceType {
	return ServiceTypeCamera
}

func (s *ServiceDahuaCamera) TryLogin(login string, password string) LoginStatus {
	return LoginNotRequired
}

func (s *ServiceDahuaCamera) StoreSnapshots(path string) error {
	return nil
}

func (s *ServiceDahuaCamera) ReadLine() []byte {
	var line []byte
	for {
		b := make([]byte, 1)
		_, err := s.Conn.Read(b)
		if err != nil {
			return nil
		}
		if b[0] == '\r' {
			next := make([]byte, 1)
			_, err := s.Conn.Read(next)
			if err != nil {
				return nil
			}
			if next[0] == '\n' {
				return line
			}
			line = append(line, '\r')
			line = append(line, next[0])
		} else {
			line = append(line, b[0])
		}
	}
}
