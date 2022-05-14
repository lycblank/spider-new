package dingding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lycblank/spider-new/pkg/notify"
	"github.com/parnurzeal/gorequest"
)

type dingDing struct {
	ch      chan notify.NotifyArg
	webhook string
	secret  string
}

func NewDingDing(webhook string, secret string) *dingDing {
	dd := &dingDing{
		ch:      make(chan notify.NotifyArg, 512),
		webhook: webhook,
		secret:  secret,
	}
	go dd.process()
	return dd
}

func (dd *dingDing) Send(ctx context.Context, arg notify.NotifyArg) error {
	select {
	case dd.ch <- arg:
	default:
	}
	return nil
}

func (dd *dingDing) Alert(title string, contents ...string) {
	if len(contents) <= 0 {
		contents = append(contents, title)
	}
	select {
	case dd.ch <- notify.NotifyArg{Title: title, Contents: contents}:
	default:
	}
}

func (dd *dingDing) process() {
	for arg := range dd.ch {
		dd.sendMsg(arg)
	}
}

func (dd *dingDing) sendMsg(arg notify.NotifyArg) {
	//timestamp := fmt.Sprintf("%d", int64(time.Now().UnixNano()/int64(time.Millisecond)))
	//fmt.Println("xxxx", dd.secret)
	//raw := timestamp + "\n" + dd.secret
	//mac := sha256.New()
	//mac.Write([]byte(raw))
	//sign := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	request := gorequest.New().Post(dd.webhook)
	request.Debug = true
	resp, _, errs := request.Set(`Content-Type`, `application/json`).
		SendString(string(dd.buildBody(arg))).End()
	if len(errs) > 0 {
		fmt.Printf("send alert failed. errs:%+v\n", errs)
		return
	}
	resp.Body.Close()
	return
}

type FlyBookMsgCell struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

func (dd *dingDing) buildBody(arg notify.NotifyArg) []byte {
	var buf bytes.Buffer
	if arg.Title != "" {
		buf.WriteString(fmt.Sprintf("## %s", arg.Title))
	}
	for i, cnt := 0, len(arg.Contents); i < cnt; i++ {
		if buf.Len() > 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("%d. %s", i+1, arg.Contents[i]))
	}

	body := map[string]interface{}{}
	body["msgtype"] = "markdown"
	body["markdown"] = map[string]interface{}{
		"title": arg.Title,
		"text":  buf.String(),
	}

	datas, _ := json.Marshal(body)
	return datas
}
