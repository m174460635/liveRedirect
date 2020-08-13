package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type YizhiboLiveService struct{}

func (YizhiboLiveService) GetPlayUrl(key string) (string, error) {
	err := getStatus(key)
	if err != nil {
		return "", err
	}
	m3u8url := "http://al01.alcdn.hls.xiaoka.tv/live/" + key + ".m3u8"
	return m3u8url, nil
}

func getStatus(key string) error {
	statusUrl := "https://m.yizhibo.com/www/live/get_live_video?scid=" + key
	res, err := requests.Get(statusUrl)
	if err != nil {
		return err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	status := json.Get([]byte(res.Text()), "data", "info", "status").ToInt()
	if status != 10 {
		return errors.New("未开播")
	}
	return nil
}

func init() {
	RegisterService("yizhibo", new(YizhiboLiveService))
}
