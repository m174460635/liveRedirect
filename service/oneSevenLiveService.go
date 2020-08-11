package service

import (
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type OneSevenLiveService struct {

}

func (s *OneSevenLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://api-dsa.17app.co/api/v1/lives/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	realUrlDefault := json.Get([]byte(res.Text()), "rtmpUrls", 0, "url").ToString()
	realUrlModify := strings.Replace(realUrlDefault, "global-pull-rtmp.17app.co", "china-pull-rtmp-17.tigafocus.com", -1)
	// 这里有点奇怪，为什么python返回数组了
	return realUrlModify, nil
}

func init() {
	RegisterService("17", new(OneSevenLiveService))
}
