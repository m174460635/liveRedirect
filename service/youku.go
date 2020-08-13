package service

/***
# 获取@优酷轮播台@的真实流媒体地址。
# 优酷轮播台是优酷直播live.youku.com下的一个子栏目，轮播一些经典电影电视剧，个人感觉要比其他直播平台影视区的画质要好，
# 而且没有平台水印和主播自己贴的乱七八糟的字幕遮挡。
# liveId 是如下形式直播间链接:
# “https://vku.youku.com/live/ilproom?spm=a2hcb.20025885.m_16249_c_59932.d_11&id=8019610&scm=20140670.rcmd.16249.live_8019610”中的8019610字段。
*/
import (
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

type YoukuService struct {
}

func (s *YoukuService) GetPlayUrl(key string) (string, error) {
	req := requests.Requests()
	tt := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	data := "{\"liveId\":\"" + key + "\",\"app\":\"Pc\"}"
	roomUrl := "https://acs.youku.com/h5/mtop.youku.live.com.livefullinfo/1.0/?appKey=24679788"
	resp, err := req.Get(roomUrl)
	if err != nil {
		return "", err
	}
	token := ""
	coo := ""
	for _, c := range resp.Cookies() {
		if c.Name == "_m_h5_tk" {
			token = c.Value[0:32]
		}
		coo = coo + c.Name + "=" + c.Value + ";"
	}
	sign := GetMD5Hash(token + "&" + tt + "&" + "24679788" + "&" + data)

	params := requests.Params{
		"t":    tt,
		"sign": sign,
		"data": data,
	}
	res, err := req.Get(roomUrl, params, coo)
	if err != nil {
		return "", err
	}
	pageResult := res.Text()
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	streamname := json.Get([]byte(pageResult), "data", "data", "stream", 0, "streamName").ToString()
	realUrl := "http://lvo-live.youku.com/vod2live/" + streamname + "_mp4hd2v3.m3u8?&expire=21600&psid=1&ups_ts=" + strconv.FormatInt(time.Now().Unix(), 10) + "&vkey="

	return realUrl, nil
}
func init() {
	RegisterService("youku", new(YoukuService))

}
