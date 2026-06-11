package services

type ServiceInterface interface {
	Init(address string)
	CanIdentify() bool
	TryLogin(login string, password string) LoginStatus
	GetName() string
	GetAddress() string
	GetType() ServiceType
	StoreSnapshots(path string) error
}
