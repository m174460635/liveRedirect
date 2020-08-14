package service

import (
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type NowLiveService struct{}

func (NowLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://now.qq.com/cgi-bin/now/web/room/get_live_room_url?room_id=" + key + "&platform=8"
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	result := json.Get([]byte(res.Text()), "result")
	raw_rtmp_url := result.Get("raw_rtmp_url")
	if raw_rtmp_url != nil {
		return raw_rtmp_url.ToString(), nil
	}
	raw_flv_url := result.Get("raw_flv_url")
	if raw_flv_url != nil {
		return raw_flv_url.ToString(), nil
	}
	return "", nil
}

func init() {
	RegisterService("now", new(NowLiveService))
}
