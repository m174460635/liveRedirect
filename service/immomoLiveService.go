package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type ImmomoLiveService struct{}

func (ImmomoLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://web.immomo.com/webmomo/api/scene/profile/roominfos"
	data := requests.Datas{
		"stid": key,
		"src":  "url",
	}
	req := requests.Requests()
	req.Get("https://web.immomo.com")
	res, err := req.Post(roomUrl, data)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	resText := res.Text()
	ec := json.Get([]byte(resText), "ec").ToInt()
	if ec != 200 {
		println(ec)
		return "", errors.New("请求参数错误")
	}
	resData := json.Get([]byte(resText), "data")
	if resData.Get("live") == nil {
		return "", errors.New("未开播")
	}
	realUrl := resData.Get("url").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("immomo", new(ImmomoLiveService))
}
