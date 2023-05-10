package translate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
	"fmt"

	"github.com/HCY2315/chaoyue-golib/log"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
)

type UrlQueue struct {
	From  string `json:"from" form:"from"`
	Q     string `json:"q" form:"q"`
	To    string `json:"to" form:"to"`
	Appid string `json:"appid" form:"appid"`
	Salt  string `json:"salt" form:"salt"`
	Sign  string `json:"sign" form:"sign"`
}

// {
//     "from": "en",
//     "to": "zh",
//     "trans_result": [
//         {
//             "src": "apple",
//             "dst": "苹果"
//         }
//     ]
// }

type TranslateRes struct {
	Form        string `json:"form"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func NewUrlQueue(from, q, to, appid, key string) *UrlQueue {
	salt := strconv.Itoa(time.Now().Second())
	return &UrlQueue{
		From:  from,
		Q:     q,
		To:    to,
		Appid: appid,
		Salt:  salt,
		Sign:  utils.Md5(appid + q + salt + key),
	}
}

func (uq *UrlQueue) Translate() (TranslateRes, error) {
	var buf bytes.Buffer
	buf.WriteString("http://fanyi-api.baidu.com/api/trans/vip/translate?")
	buf.WriteString("from=" + uq.From)
	buf.WriteString("&q=" + uq.Q)
	buf.WriteString("&to=" + uq.To)
	buf.WriteString("&appid=" + uq.Appid)
	buf.WriteString("&salt=" + uq.Salt)
	buf.WriteString("&sign=" + uq.Sign)
	fmt.Println(buf.String())
	var translateRes TranslateRes
	res, err := http.Get(buf.String())
	if err != nil {
		log.Errorf("访问百度API翻译接口出现问题, error：" + err.Error())
		return translateRes, err
	}
	resBody, err := io.ReadAll(res.Body)
	fmt.Println(string(resBody))
	if err != nil {
		log.Errorf("读取返回信息失败，error：", err.Error())
		return translateRes, err
	}
	err = json.Unmarshal(resBody, &translateRes)
	if err != nil {
		log.Errorf("解析 body 内容失败，error：", err.Error())
		return translateRes, err
	}
	return translateRes, nil
}
