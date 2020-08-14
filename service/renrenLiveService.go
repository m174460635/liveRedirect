package service

import (
	"errors"
	"github.com/asmcos/requests"
	"regexp"
	"strings"
)

type RenrenLiveService struct{}

func (RenrenLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "http://activity.renren.com/liveroom/" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	liveStateMatchRes := regexp.MustCompile(`"liveState":(\d)`).FindStringSubmatch(resText)
	if liveStateMatchRes == nil {
		return "", errors.New("直播间不存在")
	}
	playUrlMatchRes := regexp.MustCompile(`"playUrl":"([\s\S]*?)"`).FindStringSubmatch(resText)
	if playUrlMatchRes == nil {
		return "", nil
	}
	s := playUrlMatchRes[1]
	if liveStateMatchRes[1] == "1" {
		return s, nil
	}
	if liveStateMatchRes[1] != "0" {
		return "", nil
	}
	accessTokenRes := regexp.MustCompile(`accesskey=(\w+)`).FindStringSubmatch(s)
	expireRes := regexp.MustCompile(`expire=(\d+)`).FindStringSubmatch(s)
	liveRes := regexp.MustCompile(`(/live/\d+)`).FindStringSubmatch(s)
	c := accessTokenRes[1] + expireRes[1] + liveRes[1]
	md5Key := GetMD5Hash(c)
	e := strings.Split(strings.Split(s, "?")[0], "/")[4]
	realUrl := "http://ksy-hls.renren.com/live/" + e + "/index.m3u8?key=" + md5Key
	return realUrl, nil
}

func init() {
	RegisterService("renren", new(RenrenLiveService))
}
