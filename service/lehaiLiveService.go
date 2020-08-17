package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strconv"
	"time"
)

type LehaiLiveService struct{}

func (LehaiLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://service.lehaitv.com/v2/room/" + key + "/enter"
	urlEncode := url.Values{}
	tt13 := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	accessToken := "s7FUbTJ%2BjILrR7kicJUg8qr025ZVjd07DAnUQd8c7g%2Fo4OH9pdSX6w%3D%3D"
	params := requests.Params{
		"_st1":        tt13,
		"accessToken": accessToken,
		"tku":         "3000006",
	}
	for k, v := range params {
		urlEncode.Add(k, v)
	}
	data := urlEncode.Encode() + "1eha12h5"
	_ajaxData1 := GetMD5Hash(data)
	params["_ajaxData1"] = _ajaxData1
	unquotedToken, err := url.QueryUnescape(accessToken)
	if err != nil {
		return "", err
	}
	params["accessToken"] = unquotedToken
	res, err := requests.Get(roomUrl, params)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	statusCode := json.Get([]byte(resText), "status", "statuscode").ToString()
	if statusCode != "0" {
		return "", errors.New("房间不存在 或 权限检查错误")
	}
	resData := json.Get([]byte(resText), "data")
	liveStatus := resData.Get( "live_status").ToString()
	if liveStatus != "1" {
		return "", errors.New("未开播")
	}
	realUrl := resData.Get("anchor", 0, "media_url").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("lehai", new(LehaiLiveService))
}
