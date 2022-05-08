package thirdparty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/HCY2315/chaoyue-golib/pkg/errors"
	"github.com/HCY2315/chaoyue-golib/pkg/log"
)

type DingMsg interface {
	Title() string
	Content() string
}
type DingService interface {
	SendPlain(text string) error
	SendMarkdown(title, markdownText string) error
}

type simpleDingService struct {
	keyword    string
	dingbotURl string
}

func (s *simpleDingService) SendMarkdown(title, markdownText string) error {
	msgContent := fmt.Sprintf("[%s]\n%s", s.keyword, markdownText)
	msg := newDingMsgMarkdown(title, msgContent)
	bodyBuf, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "encode %+v", msg)
	}
	return s.sendMsg(bodyBuf)
}

func NewSimpleDingService(keyword, dingbotURl string) DingService {
	return &simpleDingService{
		dingbotURl: dingbotURl,
		keyword:    keyword,
	}
}

type dingMsgMarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type dingMsgMarkdown struct {
	MsgType  string                 `json:"msgtype"`
	Markdown dingMsgMarkdownContent `json:"markdown"`
}

func newDingMsgMarkdown(title, content string) *dingMsgMarkdown {
	return &dingMsgMarkdown{
		Markdown: dingMsgMarkdownContent{
			Title: title,
			Text:  content,
		},
		MsgType: "markdown",
	}
}

type dingMsgText struct {
	Content string `json:"content"`
}

type dingPlainMsg struct {
	MsgType string      `json:"msgtype"`
	Text    dingMsgText `json:"text"`
}

func newPlainDingMsg(content string) dingPlainMsg {
	return dingPlainMsg{
		MsgType: "text",
		Text:    dingMsgText{Content: content},
	}
}

func (s *simpleDingService) SendPlain(text string) error {
	msgContent := fmt.Sprintf("[%s]\n%s", s.keyword, text)
	msg := newPlainDingMsg(msgContent)
	bodyBuf, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "encode %+v", msg)
	}
	return s.sendMsg(bodyBuf)
}

func (s *simpleDingService) sendMsg(bodyBuf []byte) error {
	resp, err := http.Post(s.dingbotURl, "application/json", bytes.NewReader(bodyBuf))
	if err != nil {
		return errors.Wrap(err, "send request to %s with body:%+v", s.dingbotURl, bodyBuf)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		reqBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ding response %s with code %d. url:%s, request body:%+v", string(reqBody),
			resp.StatusCode, s.dingbotURl, string(bodyBuf))
	}
	log.Debugf("[ding] %s sent", string(bodyBuf))
	return nil
}
