package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type KugouLiveService struct {
}

func (s *KugouLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://fx1.service.kugou.com/video/mo/live/pull/h5/v3/streamaddr?roomId=" + key + "&platform=18&version=1000&streamType=3-6&liveType=1&ch=fx&ua=fx-mobile-h5&kugouId=0&layout=1"
	resp, err := requests.Get(roomUrl, requests.Header{
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36 ",
	})
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	pageResult := resp.Text()
	fmt.Println(pageResult)
	data := json.Get([]byte(pageResult), "data")
	if data == nil {
		return "", errors.New("房间不存在或者未开播")
	}
	hor := data.Get("horizontal")
	if hor == nil || hor.Size() == 0 {
		return "", errors.New("房间不存在或者未开播")
	}
	d := hor.Get(0)
	if d == nil {
		return "", errors.New("房间不存在或者未开播")
	}
	url := d.Get("httpshls").Get(0).ToString()
	return url, nil
}
func init() {
	RegisterService("kugou", new(KugouLiveService))
}
