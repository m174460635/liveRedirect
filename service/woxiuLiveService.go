package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type WoxiuLiveService struct{}

func (s *WoxiuLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://m.woxiu.com/index.php?action=M/Live&do=LiveInfo&room_id=" + key
	res, err := requests.Get(roomUrl, requests.Header{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) " +
		"Version/11.0 Mobile/15A372 Safari/604.1"})
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	status := json.Get([]byte(res.Text()), "online")
	if status == nil {
		return "", errors.New("未开播")
	}
	url := json.Get([]byte(res.Text()), "live_stream").ToString()
	return url, nil
}

func init() {
	RegisterService("woxiu", new(WoxiuLiveService))
}
