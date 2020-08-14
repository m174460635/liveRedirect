package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"strconv"
	"time"
)

type ShowselfLiveService struct{}

func (ShowselfLiveService) GetPlayUrl(key string) (string, error) {
	uidUrl := "https://service.showself.com/v2/custuser/visitor"
	res, err := requests.Get(uidUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	uidData := json.Get([]byte(res.Text()), "data")
	uid := uidData.Get("uid").ToString()
	accessToken := uidData.Get("sessionid").ToString()
	sessionid := accessToken
	data := make(map[string]string)
	data["accessToken"] = accessToken
	data["tku"] = uid
	data["_st1"] = strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	params := requests.Params{}
	copyMap(data, params)
	payload := requests.Datas{
		"groupid":   "999",
		"roomid":    key,
		"sessionid": sessionid,
		"sessionId": sessionid,
	}
	copyMap(payload, data)
	encodedData := urlEncodMap(data) + "sh0wselfh5"
	ajaxData1 := GetMD5Hash(encodedData)
	payload["_ajaxData1"] = ajaxData1

	roomUrl := "https://service.showself.com/v2/rooms/" + key + "/members?" + urlEncodMap(params)
	payloadJson, err := json.MarshalToString(payload)
	if err != nil {
		return "", nil
	}
	res, err = requests.PostJson(roomUrl, payloadJson)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	statuscode := json.Get([]byte(resText), "status", "statuscode").ToString()
	if statuscode != "0" {
		return "", errors.New("房间不存在")
	}
	dataJson := json.Get([]byte(resText), "data")
	liveStatus := dataJson.Get("roomInfo", "live_status").ToString()
	if liveStatus != "1" {
		return "", errors.New("未开播")
	}
	realUrl := dataJson.Get("roomInfo", "anchor", 0, "media_url").ToString()
	return realUrl, nil
}

func urlEncodMap(data map[string]string) string {
	sourceData := url.Values{}
	for k, v := range data {
		sourceData.Add(k, v)
	}
	return sourceData.Encode()
}

func copyMap(sourceMap, targetMap map[string]string) {
	for k, v := range sourceMap {
		targetMap[k] = v
	}
}

func init() {
	RegisterService("showself", new(ShowselfLiveService))
}
