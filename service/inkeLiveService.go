package service

import (
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type InkeLiveService struct{}

func (InkeLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://webapi.busi.inke.cn/web/live_share_pc?uid=" + key
	res, err := requests.Get(roomUrl)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data := json.Get([]byte(res.Text()), "data")
	streamAddr := data.Get("live_addr", 0, "rtmp_stream_addr").ToString()
	return streamAddr, nil
}

func init() {
	RegisterService("inke", new(InkeLiveService))
}
