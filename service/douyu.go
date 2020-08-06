package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DouyuLiveService struct {
}

func getTt() [3]string {
	tt1 := strconv.FormatInt(time.Now().Unix(), 10)
	tt2 := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	today := time.Now().Format("20060102")
	return [3]string{tt1, tt2, today}
}
func mixRoom(rid string) string {
	result1 := "PKing"
	return result1
}
func getPreUrl(rid string, tt string) string {
	requestUrl := "https://playweb.douyucdn.cn/lapi/live/hlsH5Preview/" + rid
	postData := requests.Datas{
		"rid": rid,
		"did": "10000000000000000000000000001501",
	}
	tt = "1596636061448"
	auth := GetMD5Hash(rid + tt)
	header := requests.Header{
		"content-type": "application/x-www-form-urlencoded",
		"rid":          rid,
		"time":         tt,
		"auth":         auth,
	}
	resp, _ := requests.Post(requestUrl, postData, header)
	preUrl := ""
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	pageResult := resp.Text()
	if json.Get([]byte(pageResult), "error").ToInt() == 0 {
		realUrl := json.Get([]byte(pageResult), "data", "rtmp_live").ToString()
		if strings.Contains(realUrl, "mix=1") {
			preUrl = mixRoom(rid)
		} else {
			pattern1 := `(?i)^[0-9a-zA-Z]*`
			re := regexp.MustCompile(pattern1)
			res := re.FindStringSubmatch(realUrl)
			preUrl = res[0]
		}

	}

	return preUrl
}

func (s *DouyuLiveService) GetPlayUrl(key string) (string, error) {
	tt := getTt()

	url := getPreUrl(key, tt[1])
	if url != "" {
		return "http://tx2play1.douyucdn.cn/live/" + url + ".flv?uuid=", nil
	}

	return "", errors.New("直播间不存在或者未开播")
}
