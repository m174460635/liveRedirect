package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type ZhanqiLiveService struct {
}

func (s *ZhanqiLiveService) GetPlayUrl(key string) (string, error) {
	resp, err := requests.Get("https://m.zhanqi.tv/api/static/v2.1/room/domain/" + key + ".json")
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	res := make(map[string]interface{})
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonError := json.Unmarshal([]byte(resp.Text()), &res)
	if jsonError != nil {
		return "", err
	}
	cookies := resp.Cookies()
	videoid := res["data"].(map[string]interface{})["videoId"]
	status := res["data"].(map[string]interface{})["status"]

	coo := ""
	for _, c := range cookies {
		coo = coo + c.Name + "=" + c.Value + ";"
	}

	if "4" == status {
		resp, err = requests.Get("https://dlhdl-cdn.zhanqi.tv/zqlive/"+(videoid.(string))+".flv?get_url=1", requests.Header{"Cookie": coo})
		if err != nil {
			fmt.Print(err.Error())
			return "", err
		}
		realUrl := resp.Text()
		return realUrl, nil
	} else {
		fmt.Println("未开播")
		return "", errors.New("未开播")
	}

	return "http://www.baidu.com", nil
}
func init() {
	RegisterService("zhanqi", new(ZhanqiLiveService))
}
