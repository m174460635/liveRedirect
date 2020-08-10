package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type NineXiuLiveService struct {

}

func (s *NineXiuLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://h5.9xiu.com/room/live/enterRoom?rid=" + key
	headers := requests.Header{
		"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
	}
	res, err := requests.Get(roomUrl, headers)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	resCode := json.Get([]byte(resText), "code").ToInt()
	if resCode != 200 {
		return "", errors.New("直播间不存在")
	}
	data := json.Get([]byte(resText), "data")
	status := data.Get("status").ToInt()
	if status == 0 {
		return "", errors.New("未开播")
	}
	liveUrl := data.Get("live_url").ToString()
	return liveUrl, nil
}

func init() {
	RegisterService("9xiu", new(NineXiuLiveService))
}
