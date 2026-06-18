package services

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

type ServiceWebDahuaCamera struct {
	ServiceInterface
	Address string
}

func (s *ServiceWebDahuaCamera) Init(address string) {
	s.Address = address
}

func (s *ServiceWebDahuaCamera) GetAddress() string {
	return s.Address
}

func (s *ServiceWebDahuaCamera) MD5Of(path string) (string, error) {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get("http://" + s.Address + path)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	hashBytes := md5.Sum(body)

	hashString := hex.EncodeToString(hashBytes[:])

	return hashString, nil
}

func (s *ServiceWebDahuaCamera) CanIdentify() bool {
	faviconHashes := []string{"bd9e17c46bbbc18af2a2bd718dddad0e", "605f51b413980667766a9aff2e53b9ed", "b39f249362a2e4ab62be4ddbc9125f53"}

	faviconMD5, err := s.MD5Of("/favicon.ico")
	if err != nil {
		return false
	}

	imgMD5, err := s.MD5Of("/image/lgbg.jpg")
	if err != nil {
		return false
	}
	return slices.Contains(faviconHashes, faviconMD5) || imgMD5 == "4ff53be6165e430af41d782e00207fda"
}

func (s *ServiceWebDahuaCamera) TryLogin(user string, password string) LoginStatus {
	ip := strings.Split(s.Address, ":")[0]
	headers := map[string]string{
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36",
		"Host":             ip,
		"Origin":           "http://" + ip,
		"Referer":          "http://" + ip,
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding":  "gzip, deflate",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"Connection":       "close",
		"X-Requested-With": "XMLHttpRequest",
	}

	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"method": "global.login",
		"params": map[string]interface{}{
			"userName":      user,
			"password":      password,
			"clientType":    "Web3.0",
			"loginType":     "Direct",
			"authorityType": "Default",
			"passwordType":  "Plain",
		},
		"id":      1,
		"session": 0,
	})
	println(string(jsonData))
	if err != nil {
		return LoginFailed
	}

	req, _ := http.NewRequest("POST", "http://"+s.Address+"/RPC2_Login", bytes.NewBuffer(jsonData))

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return LoginFailed
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoginFailed
	}
	stringBody := strings.ReplaceAll(string(body), " ", "")
	if strings.Contains(stringBody, "\"result\":true") && resp.StatusCode == 200 {
		return LoginSuccess
	}
	println(stringBody)
	return LoginFailed
}

func (s *ServiceWebDahuaCamera) GetName() string {
	return "Dahua Camera (WEB)"
}

func (s *ServiceWebDahuaCamera) GetType() ServiceType {
	return ServiceTypeCamera
}

func (s *ServiceWebDahuaCamera) StoreSnapshots(path string) error {
	return nil
}
