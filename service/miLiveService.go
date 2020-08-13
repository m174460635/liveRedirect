package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

//小米直播：https://live.wali.com/fe
type MiLiveService struct{}

func (MiLiveService) GetPlayUrl(key string) (string, error) {
	zuid := strings.Split(key, "_")[0]
	roomUrl := "https://s.zb.mi.com/get_liveinfo?lid=" + key + "&zuid=" + zuid
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := json.Get([]byte(res.Text()), "data")
	if data.Get("status").ToInt() != 1 {
		return "", errors.New("直播间不存在或未开播")
	}
	return strings.Replace(data.Get("video", "flv").ToString(), "http", "https", -1), nil

}

func init() {
	RegisterService("mi", new(MiLiveService))
}
