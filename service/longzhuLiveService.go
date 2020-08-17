package service

import (
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"regexp"
)

type LongzhuLiveService struct{}

func (LongzhuLiveService) GetPlayUrl(key string) (string, error) {
	res, err := requests.Get("http://m.longzhu.com/" + key)
	if err != nil {
		return "", err
	}
	roomIdMatchRes := regexp.MustCompile(`roomId = (\d*);`).FindStringSubmatch(res.Text())
	if roomIdMatchRes == nil {
		return "", nil
	}
	res, err = requests.Get("http://livestream.longzhu.com/live/getlivePlayurl?roomId=" + roomIdMatchRes[1] + "&hostPullType=2&isAdvanced=true&playUrlsType=1")
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	urls := json.Get([]byte(res.Text()), "playLines", 0, "urls")
	realUrl := urls.Get(urls.Size() - 1, "securityUrl").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("longzhu", new(LongzhuLiveService))
}
