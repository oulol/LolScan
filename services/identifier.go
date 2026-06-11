package services

var Registry []ServiceInterface = []ServiceInterface{
	&ServiceWeb{},
}

func Identify(address string) ServiceInterface {
	for _, dev := range Registry {
		dev.Init(address)
		if dev.CanIdentify() {
			return dev
		}
	}

	return nil
}
