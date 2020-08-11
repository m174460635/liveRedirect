package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

//艺气山直播：http://www.173.com/room/category?categoryId=11
type OneSevenThreeLiveService struct {

}

func (s *OneSevenThreeLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "http://www.173.com/room/getVieoUrl"
	params := requests.Params{
		"roomId": key,
		"format": "m3u8",
	}
	if res, err := requests.Post(roomUrl, params);err == nil {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		dataIntf := json.Get([]byte(res.Text()), "data")
		if dataIntf == nil {
			return "", errors.New("直播间不存在")
		}
		status := dataIntf.Get("status").ToInt()
		if status != 2 {
			return "", errors.New("未开播")
		}
		return dataIntf.Get("url").ToString(), nil
	}
	return "", nil
}

func init() {
	RegisterService("173", new(OneSevenThreeLiveService))
}
