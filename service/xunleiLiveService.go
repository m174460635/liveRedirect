package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strconv"
	"time"
)

type XunleiLiveService struct{}

func (s *XunleiLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://biz-live-ssl.xunlei.com//caller"
	headers := requests.Header{
		"cookie": "appid=1002",
	}
	tt13 := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	data := url.Values{}
	data.Add("_t", tt13)
	data.Add("a", "play")
	data.Add("c", "room")
	data.Add("hid", "h5-e70560ea31cc17099395c15595bdcaa1")
	data.Add("uuid", key)
	params := requests.Params{
		"_t":   tt13,
		"a":    "play",
		"c":    "room",
		"hid":  "h5-e70560ea31cc17099395c15595bdcaa1",
		"uuid": key,
	}
	f := "&*%$7987321GKwq"
	p := GetMD5Hash("1002" + data.Encode() + f)
	params["sign"] = p
	res, err := requests.Get(roomUrl, params, headers)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	if json.Get([]byte(res.Text()), "result").ToInt() != 0 {
		return "", errors.New("直播间不存在")
	}
	resData := json.Get([]byte(res.Text()), "data")
	if resData.Get("play_status").ToInt() != 1 {
		return "", errors.New("未开播")
	}
	realUrl := resData.Get("data", "stream_pull_https").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("xunlei", new(XunleiLiveService))
}
