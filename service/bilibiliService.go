package service

import (
	"errors"
	"fmt"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type BilibiliService struct {
}

func u(roomId string, pf string, cookies string) (string, error) {
	req := requests.Requests()

	fUrl := "https://api.live.bilibili.com/xlive/web-room/v1/playUrl/playUrl"
	p := requests.Params{
		"cid":           roomId,
		"qn":            "10000",
		"platform":      pf,
		"https_url_req": "1",
		"ptype":         "16",
	}
	res, err := req.Get(fUrl, p, requests.Header{"Cookie": cookies})
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	durl := json.Get([]byte(res.Text()), "data", "durl")
	if durl == nil || durl.Size() == 0 {
		return "", errors.New("房间不存在")
	}
	s := durl.Get(durl.Size() - 1).Get("url").ToString()
	return s, nil
}
func (s *BilibiliService) GetPlayUrl(key string) (string, error) {
	rUrl := "https://api.live.bilibili.com/room/v1/Room/room_init?id=" + key
	resp, err := requests.Get(rUrl)
	if err != nil {
		fmt.Print(err.Error())
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	pageResult := resp.Text()
	code := json.Get([]byte(pageResult), "code").ToInt()
	if code != 0 {
		return "", errors.New("房间不存在")
	}

	cookies := resp.Cookies()
	coo := ""
	for _, c := range cookies {
		coo = coo + c.Name + "=" + c.Value + ";"
	}

	liveStatus := json.Get([]byte(pageResult), "data", "live_status").ToInt()
	if liveStatus != 1 {
		return "", errors.New("未开播")
	}
	roomId := json.Get([]byte(resp.Text()), "data", "room_id").ToString()

	return u(roomId, "h5", coo)
}
func init() {
	RegisterService("bilibili", new(BilibiliService))

}
