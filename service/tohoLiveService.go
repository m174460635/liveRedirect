package service

import (
	"errors"
	"github.com/asmcos/requests"
	"regexp"
	"strings"
)

//星光直播：https://www.tuho.tv/28545037
type TohoLiveService struct{}

func (TohoLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://www.tuho.tv/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	flvMatchRes := regexp.MustCompile(`videoPlayFlv":"(https[\s\S]+?flv)`).FindStringSubmatch(resText)
	if flvMatchRes == nil {
		return "", nil
	}
	statusMatchRes := regexp.MustCompile(`isPlaying\s:\s(\w+),`).FindStringSubmatch(resText)
	if statusMatchRes == nil || statusMatchRes[1] != "true" {
		return "", errors.New("未开播")
	}
	realUrl := strings.Replace(flvMatchRes[1], "\\", "", -1)
	return realUrl, nil
}

func init() {
	RegisterService("toho", new(TohoLiveService))
}
