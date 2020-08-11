package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
)

type AcFunLiveService struct {
}

func (s *AcFunLiveService) GetPlayUrl(key string) (string, error) {
	roomUrl := "https://id.app.acfun.cn/rest/app/visitor/login"
	headers := requests.Header{
		"content-type": "application/x-www-form-urlencoded",
		"cookie":       "_did=H5_",
		"referer":      "https://m.acfun.cn/",
	}
	datas := requests.Datas{
		"sid": "acfun.api.visitor",
	}
	res, err := requests.Post(roomUrl, datas, headers)
	if err != nil {
		return "", err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	userId := json.Get([]byte(res.Text()), "userId").ToString()
	visitorSt := json.Get([]byte(res.Text()), "acfun.api.visitor_st").ToString()

	playUrl := "https://api.kuaishouzt.com/rest/zt/live/web/startPlay"
	params := requests.Params{
		"subBiz":               "mainApp",
		"kpn":                  "ACFUN_APP",
		"kpf":                  "PC_WEB",
		"userId":               userId,
		"did":                  "H5_",
		"acfun.api.visitor_st": visitorSt,
	}
	datas = requests.Datas{
		"authorId":       key,
		"pullStreamType": "FLV",
	}
	res, err = requests.Post(playUrl, params, datas, headers)
	if err != nil {
		return "", err
	}
	resText := res.Text()
	result := json.Get([]byte(resText), "result").ToInt()
	if result != 1 {
		return "", errors.New("直播已关闭")
	}
	data := json.Get([]byte(resText), "data")
	liveAdaptiveManifest := json.Get([]byte(data.Get("videoPlayRes").ToString()), "liveAdaptiveManifest", 0)
	realUrl := liveAdaptiveManifest.Get("adaptationSet", "representation", liveAdaptiveManifest.Get("adaptationSet", "representation").Size()-1, "url").ToString()
	return realUrl, nil
}

func init() {
	RegisterService("acfun", new(AcFunLiveService))
}
