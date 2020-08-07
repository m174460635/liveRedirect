package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/markbates/pkger"
	"io/ioutil"
)

type LiveService interface {
	GetPlayUrl(key string) (string, error)
}

var serviceMap = make(map[string]LiveService)

func GetServiceMap() map[string]LiveService {
	return serviceMap
}
func RegisterService(path string, s LiveService) {
	serviceMap[path] = s
	fmt.Println("加载：" + path + " 服务")
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func GetFromResource(path string) (string, error) {
	f, err := pkger.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	text := string(content)
	return text, nil
}
