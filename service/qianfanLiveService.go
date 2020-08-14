package service

import (
	"github.com/asmcos/requests"
	"regexp"
)

type QianFanLiverService struct{}

func (QianFanLiverService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://qf.56.com/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	urlMatchRes := regexp.MustCompile(`flvUrl:'(.*)\?wsSecret`).FindStringSubmatch(res.Text())
	if urlMatchRes == nil {
		return "", nil
	}
	return urlMatchRes[1], nil
}

func init() {
	RegisterService("56qf", new(QianFanLiverService))
}
