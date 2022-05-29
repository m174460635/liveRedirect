package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"liveRedirect/jsengine"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DouyuLiveService struct {
}

func getTt() [2]string {
	tt10 := strconv.FormatInt(time.Now().Unix(), 10)
	tt13 := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	return [2]string{tt10, tt13}
}

func (s *DouyuLiveService) get_js(key, did, realKeyRespText string, tt [2]string) (string, error) {
	jsReg := regexp.MustCompile(`(function ub98484234.*)\s(var.*)`)
	jsRegMatchResult := jsReg.FindStringSubmatch(realKeyRespText)
	if jsRegMatchResult == nil {
		return "", errors.New("未找到直播地址")
	}
	replaceReg := regexp.MustCompile(`eval.*;}`)
	jsContent := jsRegMatchResult[2] + replaceReg.ReplaceAllString(jsRegMatchResult[1], "strc;}")
	jsRes, err := jsengine.RunJSFunc(jsContent, "ub98484234")
	if err != nil {
		return "", err
	}
	vMatchResult := regexp.MustCompile(`v=(\d+)`).FindStringSubmatch(jsRes)
	if vMatchResult == nil {
		return "", errors.New("未找到直播地址")
	}
	rb := GetMD5Hash(key + did + tt[0] + vMatchResult[1])
	func_sign := regexp.MustCompile(`return rt;}\);?`).ReplaceAllString(jsRes, "return rt;}")
	func_sign = strings.ReplaceAll(func_sign, `(function (`, `function sign(`)
	func_sign = strings.ReplaceAll(func_sign, `CryptoJS.MD5(cb).toString()`, `"`+rb+`"`)
	signRes, err := jsengine.RunJSFunc(func_sign, "sign", key, did, tt[0])
	if err != nil {
		return "", err
	}
	signRes = signRes + "&ver=219032101&rid=" + key + "&rate=-1"
	rateRes, err := requests.Post("https://m.douyu.com/api/room/ratestream?" + signRes)
	if err != nil {
		return "", err
	}
	rateResText := rateRes.Text()
	jsKeyMatchRes := regexp.MustCompile(`(\d{1,7}[0-9a-zA-Z]+)_?\d{0,4}(.m3u8|/playlist)`).FindStringSubmatch(rateResText)
	if jsKeyMatchRes == nil {
		return "", nil
	}
	return jsKeyMatchRes[1], nil
}

func (s *DouyuLiveService) GetPlayUrl(key string) (string, error) {
	var did = "10000000000000000000000000001501"
	tt := getTt()
	realKey, realKeyRespText, err := getRealKey(key)
	if err != nil {
		return "", err
	}
	roomUrl := "https://playweb.douyucdn.cn/lapi/live/hlsH5Preview/" + realKey
	auth := GetMD5Hash(realKey + tt[1])
	data := requests.Datas{
		"rid": realKey,
		"did": did,
	}
	headers := requests.Header{
		"rid":  realKey,
		"time": tt[1],
		"auth": auth,
	}
	roomRes, err := requests.Post(roomUrl, data, headers)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	preRoomText := roomRes.Text()
	errorCode := json.Get([]byte(preRoomText), "error").ToInt()
	var target string
	if errorCode == 102 {
		return "", errors.New("房间不存在")
	} else if errorCode == 104 {
		return "", errors.New("房间未开播")
	} else if errorCode == 0 {
		dataJson := json.Get([]byte(preRoomText), "data")
		rtmp_live := dataJson.Get("rtmp_live").ToString()
		preRoomReg := regexp.MustCompile(`(\d{1,7}[0-9a-zA-Z]+)_?\d{0,4}(/playlist|.m3u8)`)
		preRoomMatchResult := preRoomReg.FindStringSubmatch(rtmp_live)
		if preRoomMatchResult == nil {
			return "", errors.New("未找到房间URL")
		}
		target = preRoomMatchResult[1]
	} else {
		key, err := s.get_js(realKey, did, realKeyRespText, tt)
		if err != nil {
			return "", err
		}
		target = key
	}
	if len(target) <= 0 {
		return "", nil
	}
	return "http://vplay3a.douyucdn.cn/live/" + target + ".flv?uuid=", nil
}

func getRealKey(key string) (string, string, error) {
	res, err := requests.Get("https://m.douyu.com/" + key)
	if err != nil {
		return "", "", err
	}

	reg := regexp.MustCompile(`rid":(\d{1,7}),"vipId`)
	realKeyRespText := res.Text()
	matchedResult := reg.FindStringSubmatch(realKeyRespText)
	if matchedResult == nil {
		return "", "", errors.New("房间号错误")
	}
	realKey := matchedResult[1]
	return realKey, realKeyRespText, nil
}

func init() {
	RegisterService("douyu", new(DouyuLiveService))
}
