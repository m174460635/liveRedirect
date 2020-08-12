package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	"regexp"
)

type YuanboboLiveService struct {
}

func (s *YuanboboLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://zhibo.yuanbobo.com/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	streamIdMatchRes := regexp.MustCompile(`stream_id:\s+'(\d+)`).FindStringSubmatch(resText)
	if streamIdMatchRes == nil {
		return "", errors.New("直播间不存在")
	}
	statusMatchRes := regexp.MustCompile(`status:\s+'(\d)`).FindStringSubmatch(resText)
	if statusMatchRes == nil || statusMatchRes[1] != "1" {
		return "", errors.New("未开播")
	}
	realUrl := fmt.Sprintf("http://ks-hlslive.yuanbobo.com/live/%s/index.m3u8", streamIdMatchRes[1])
	return realUrl, nil
}

func init() {
	RegisterService("yuanbobo", new(YuanboboLiveService))
}
