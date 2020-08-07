package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"time"
)

type LMLiveService struct {
}

func g(d, key string) string {
	return GetMD5Hash(d + "#programId=" + key + "#Ogvbm2ZiKE")
}

func (s *LMLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://www.51lm.tv/live/room/info/basic"
	lminfo := url.Values{}

	lminfo.Add("h", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	lminfo.Add("i", "-246397986")
	lminfo.Add("o", "iphone")
	lminfo.Add("s", "G_c17a64eff3f144a1a48d9f0 ̰2e8d981c2")
	lminfo.Add("t", "H")
	lminfo.Add("v", "4.20.43")
	lminfo.Add("w", "a710244508d3cc14f50d24e9fecc496a")
	encodedLmInfo := lminfo.Encode()
	u := g(encodedLmInfo, key)
	jsonStr := "{\"programId\":\"" + key + "\"}"
	headerInfo := "G=" + u + "&" + encodedLmInfo
	res, err := requests.PostJson(roomUrl, jsonStr, requests.Header{"lminfo": headerInfo})
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	resText := res.Text()
	if len(resText) <= 0 {
		return "", nil
	}
	code := json.Get([]byte(resText), "code").ToInt()
	if code == -1 {
		return "", errors.New("输入信息错误")
	} else if code == 1201 {
		return "", errors.New("直播间不存在")
	} else if code != 200 {
		return "", nil
	}
	status := json.Get([]byte(resText), "data", "isLiving").ToString()
	if status != "True" {
		return "", errors.New("未开播")
	}
	return json.Get([]byte(resText), "data", "playUrl").ToString(), nil
}
func init() {
	RegisterService("51lm", new(LMLiveService))
}
