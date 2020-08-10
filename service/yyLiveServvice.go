package service

import (
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"regexp"
)

type YYLiveService struct {
}

func (s *YYLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "http://interface.yy.com/hls/new/get/" + key + "/" + key + "/1200?source=wapyy&callback=jsonp3"
	resp, err := requests.Get(roomUrl, requests.Header{"referer": "http://wap.yy.com/mobileweb/" + key,
		"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36 ",
	})
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	pageResult := resp.Text()
	re := regexp.MustCompile(`\(([\W\w]*)\)`)
	res := re.FindStringSubmatch(pageResult)

	if len(res) > 0 {
		s := res[1]
		if len(s) > 0 {
			d := make(map[string]interface{})
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err := json.Unmarshal([]byte(s), &d)
			if err != nil {
				return "", err
			}
			url := fmt.Sprintf("%v", d["hls"])
			return url, nil
		}
	}

	return "", nil
}
func init() {
	RegisterService("yy", new(YYLiveService))
}