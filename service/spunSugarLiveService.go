package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type SpunSugarLiveService struct {
}

func (s *SpunSugarLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://www.2cq.com/proxy/room/room/info?roomId=" + key + "&appId=1004"
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	resStr := res.Text()
	if len(resStr) <= 0 {
		return "", errors.New("查询失败")
	}
	if json.Get([]byte(resStr), "status").ToInt() != 1 {
		return "", errors.New("房间不存在")
	}
	if json.Get([]byte(resStr), "result", "liveState").ToInt() != 1 {
		return "", errors.New("未开播")
	}
	realUrl := json.Get([]byte(resStr), "result", "pullUrl").ToString()
	if len(realUrl) <= 0 {
		return "", nil
	}
	return realUrl, nil
}
func init() {
	RegisterService("2cp", new(SpunSugarLiveService))
}
