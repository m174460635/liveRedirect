package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hr3lxphr6j/requests"
	"github.com/tidwall/gjson"

	"github.com/luckycat0426/bililive-client/pkg/utils"
)

var (
	ErrRoomNotExist     = errors.New("room not exists")
	ErrRoomUrlIncorrect = errors.New("room url incorrect")
	ErrInternalError    = errors.New("internal error")
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
)

var CommonUserAgent = requests.UserAgent(userAgent)
var timeout = requests.Timeout(time.Second * 30)

func decode_live_url_info(srcAntiCode string) *map[string]string {
	c := strings.Split(srcAntiCode, "&")
	var cc = []string{}

	n := make(map[string]string)
	for _, x := range c {
		if len(x) > 0 {
			cc = append(cc, x)
			ss := strings.Split(x, "=")
			n[ss[0]] = ss[1]
		}
	}
	c = cc
	fm, _ := url.QueryUnescape(n["fm"])
	uu, _ := base64.StdEncoding.DecodeString(fm)
	u := string(uu)

	liveUrlInfo := make(map[string]string)
	liveUrlInfo["hash_prefix"] = strings.Split(u, "_")[0]
	liveUrlInfo["uuid"] = getOrDefault(&n, "uuid", "")
	liveUrlInfo["ctype"] = getOrDefault(&n, "ctype", "")
	liveUrlInfo["txyp"] = getOrDefault(&n, "txyp", "")
	liveUrlInfo["fs"] = getOrDefault(&n, "fs", "")
	liveUrlInfo["t"] = getOrDefault(&n, "t", "")

	return &liveUrlInfo
}
func getOrDefault(m *map[string]string, key string, defa string) string {
	if value, ok := (*m)[key]; ok {
		return value
	} else {
		return defa
	}

}
func GetHuyaStreamUrls(uurl string) (us []string, err error) {
	resp, err := requests.Get(uurl, CommonUserAgent, timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrRoomNotExist
	}
	body, err := resp.Text()
	if err != nil {
		return nil, err
	}

	// Decode stream part.
	streamInfo := utils.Match1(`"stream": "(.*?)"`, body)
	if streamInfo == "" {
		return nil, ErrInternalError
	}
	streamByte, err := base64.StdEncoding.DecodeString(streamInfo)
	if err != nil {
		return nil, err
	}
	streamStr := utils.UnescapeHTMLEntity(string(streamByte))
	streamInfoList := gjson.Get(streamStr, "data.0.gameStreamInfoList")

	liveUrlInfos := make(map[string]map[string]string)
	streamInfoList.ForEach(func(_, streamInfo gjson.Result) bool {
		liveUrlInfo := make(map[string]string)
		sCdnType := streamInfo.Get("sCdnType").String()
		liveUrlInfo["stream_name"] = streamInfo.Get("sStreamName").String()
		liveUrlInfo["base_url"] = streamInfo.Get("sHlsUrl").String()
		liveUrlInfo["hls_url"] = streamInfo.Get("sHlsUrl").String() + "/" + streamInfo.Get("sStreamName").String() + "." + streamInfo.Get("sHlsUrlSuffix").String()
		liveUrlInfo["sCdnType"] = sCdnType
		sHlsAntiCode := streamInfo.Get("sHlsAntiCode").String()
		info := decode_live_url_info(sHlsAntiCode)
		for k, v := range *info {
			liveUrlInfo[k] = v
		}
		liveUrlInfos[sCdnType] = liveUrlInfo

		return true
	})

	i := time.Now().UnixNano() / 1000000
	i2 := time.Now().UnixNano() / 1000000000
	seqid := strconv.FormatInt(i+1463993859134, 10)
	i2 = 1653882456
	wsTime := toHex(int(i2))

	var ss []string
	for _, liveUrlInfo := range liveUrlInfos {
		hash0 := GetMD5Hash(seqid + "|" + liveUrlInfo["ctype"] + "|" + liveUrlInfo["t"])
		hash1 := GetMD5Hash(liveUrlInfo["hash_prefix"] + "_" + "1463993859134" + "_" + liveUrlInfo["stream_name"] + "_" + hash0 + "_" + wsTime)
		ratio := ""
		if strings.Contains(liveUrlInfo["ctype"], "mobile") {
			url := fmt.Sprintf("%s?wsSecret=%s&wsTime=%s&uuid=%s&uid=%s&seqid=%s&ratio=%s&txyp=%s&fs=%s&ctype=%s&ver=1&t=%s",
				liveUrlInfo["hls_url"], hash1, wsTime, liveUrlInfo["uuid"], "1463993859134", seqid, ratio, liveUrlInfo["txyp"],
				liveUrlInfo["fs"], liveUrlInfo["ctype"], liveUrlInfo["t"])

			ss = append(ss, url)
		} else {
			url := fmt.Sprintf("%s?wsSecret=%s&wsTime=%s&seqid=%s&ctype=%s&ver=1&txyp=%s&fs=%s&ratio=%s&u=%s&t=%s&sv=2107230339",
				liveUrlInfo["hls_url"], hash1, wsTime, seqid, liveUrlInfo["ctype"], liveUrlInfo["txyp"], liveUrlInfo["fs"], ratio, "1463993859134", liveUrlInfo["t"])
			ss = append(ss, url)

		}

	}

	return ss, nil
}
func toHex(ten int) string {
	m := 0
	hex := make([]int, 0)
	for {
		m = ten % 16
		ten = ten / 16
		if ten == 0 {
			hex = append(hex, m)
			break
		}
		hex = append(hex, m)
	}
	hexStr := []string{}
	for i := len(hex) - 1; i >= 0; i-- {
		if hex[i] >= 10 {
			hexStr = append(hexStr, fmt.Sprintf("%c", 'A'+hex[i]-10))
		} else {
			hexStr = append(hexStr, fmt.Sprintf("%d", hex[i]))
		}
	}
	return strings.ToLower(strings.Join(hexStr, ""))
}
