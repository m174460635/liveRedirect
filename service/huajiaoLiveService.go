package service

import (
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type HuajiaoLiveService struct {
}

func (s *HuajiaoLiveService) GetPlayUrl(key string) (string, error) {
	timeStr := fmt.Sprintf("%.7f", float64(time.Now().UnixNano())/1000000000)
	roomUrl := "https://h.huajiao.com/api/getFeedInfo?sid=" + timeStr + "&liveid=" + key
	resp, err := requests.Get(roomUrl, requests.Header{"referer": "https://h.huajiao.com/l/feedlist",
		"User-Agent": "Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
	})
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	resultContent := resp.Text()
	if len(resultContent) <= 0 {
		return "", nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData := make(map[string]interface{})
	jsonerr := json.Unmarshal([]byte(resultContent), &jsonData)
	realUrl := json.Get([]byte(resultContent), "data", "live", "main").ToString()
	if len(realUrl) <= 0 {
		return "", jsonerr
	}

	return realUrl, nil
}
func init() {
	RegisterService("huajiao", new(HuajiaoLiveService))
}
