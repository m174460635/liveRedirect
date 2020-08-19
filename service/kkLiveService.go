package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type KKLiveService struct{}

func (KKLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://sapi.kktv1.com/meShow/entrance?parameter="
	paramter := `{'FuncTag': 10005043, 'userId': '` + key + `', 'platform': 1, 'a': 1, 'c': 100101}`
	res, err := requests.Get(roomUrl + paramter)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	resText := res.Text()
	tagCode := json.Get([]byte(resText), "TagCode").ToString()
	if tagCode != "00000000" {
		return "", errors.New("直播间不存在")
	}
	if json.Get([]byte(resText), "liveType").ToInt() != 1 {
		return "", errors.New("未开播")
	}
	roomId := json.Get([]byte(resText), "roomId").ToString()
	paramter = `{'FuncTag': 60001002, 'roomId': ` + roomId + `, 'platform': 1, 'a': 1, 'c': 100101}`
	res, err = requests.Get(roomUrl + paramter)
	if err != nil {
		return "", err
	}
	realUrl := json.Get([]byte(res.Text()), "liveStream").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("kk", new(KKLiveService))
}
