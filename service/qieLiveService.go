package service

import (
	"errors"
	"github.com/asmcos/requests"
	"regexp"
)

type QieLiveService struct{}

func (QieLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://m.live.qq.com/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	statusMatchRes := regexp.MustCompile(`"show_status":"(\d)"`).FindStringSubmatch(resText)
	if statusMatchRes == nil || statusMatchRes[1] != "1" {
		return "", errors.New("直播间不存在或未开播")
	}
	hlsUrlMatchRes := regexp.MustCompile(`"hls_url":"(.*)","use_p2p"`).FindStringSubmatch(resText)
	if hlsUrlMatchRes == nil {
		return "", nil
	}
	return hlsUrlMatchRes[1], nil
}

func init() {
	RegisterService("qie", new(QieLiveService))
}
