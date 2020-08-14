package service

import (
	"github.com/asmcos/requests"
	"regexp"
	"strconv"
	"time"
)

type PPSLiveService struct{}

func (PPSLiveService) GetPlayUrl(key string) (string, error) {
	anchorUrl := "http://m-x.pps.tv/room/" + key
	res, err := requests.Get(anchorUrl)
	if err != nil {
		return "", err
	}
	anchorIdMatchRes := regexp.MustCompile(`anchor_id":(\d*),"online_uid`).FindStringSubmatch(res.Text())
	if anchorIdMatchRes == nil {
		return "", nil
	}
	tt13 := strconv.FormatInt((time.Now().UnixNano() / 1000000), 10)
	roomUrl := "http://api-live.iqiyi.com/stream/geth5?qd_tm=" + tt13 + "&typeId=1&platform=7&vid=0&qd_vip=0&qd_uid=" + anchorIdMatchRes[1] + "&qd_ip=114.114.114.114&qd_vipres=0&qd_src=h5_xiu&qd_tvid=0&callback="
	headers := requests.Header{
		"Content-Type": "application/x-www-form-urlencoded",
		"Referer":      "http://m-x.pps.tv/",
	}
	res, err = requests.Get(roomUrl, headers)
	if err != nil {
		return "", err
	}
	realUrlMatchRes := regexp.MustCompile(`"hls":"(.*)","rate_list`).FindStringSubmatch(res.Text())
	if realUrlMatchRes == nil {
		return "", nil
	}
	return realUrlMatchRes[1], nil
}

func init() {
	RegisterService("pps", new(PPSLiveService))
}
