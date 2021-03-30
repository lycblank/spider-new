package flybook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type FlyBookArg struct {
	Title    string
	Contents []string
}

type FlyBook struct {
	ch       chan FlyBookArg
	webhook  string
}

func NewFlyBook(webhook string) *FlyBook {
	fb := &FlyBook{
		ch:       make(chan FlyBookArg, 512),
		webhook:  webhook,
	}
	go fb.process()
	return fb
}



func (fb *FlyBook) Write(p []byte) (n int, err error) {
	datas := bytes.Split(p, []byte("\n"))
	if len(datas) > 0 {
		contents := make([]string, 0, len(datas) - 1)
		for i ,cnt := 1, len(datas);i<cnt;i++{
			contents = append(contents, string(datas[i]))
		}
		fb.Alert(string(datas[0]), contents...)
	}

	return len(p), nil
}


func (fb *FlyBook) Alert(title string, contents ...string) {
	if len(contents) <= 0 {
		contents = append(contents, title)
	}
	select {
	case fb.ch <- FlyBookArg{Title: title, Contents: contents}:
	default:
	}
}

func (fb *FlyBook) process() {
	for arg := range fb.ch {
		fb.sendMsg(arg)
	}
}

func (fb *FlyBook) sendMsg(arg FlyBookArg) {
	resp, _, errs := gorequest.New().Post(fb.webhook).
		Set(`Content-Type`, `application/json`).
		SendString(string(fb.buildBody(arg))).End()
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

func (fb *FlyBook) buildBody(arg FlyBookArg) []byte {
	title := arg.Title
	lines := make([][]FlyBookMsgCell, 0, len(arg.Contents))
	for i, cnt := 0, len(arg.Contents); i < cnt; i++ {
		lines = append(lines, []FlyBookMsgCell{{Tag: "text", Text: arg.Contents[i]}})
	}

	body := map[string]interface{}{}
	body["msg_type"] = "post"
	body["content"] = map[string]interface{}{
		"post": map[string]interface{}{
			"zh_cn": map[string]interface{}{
				"title":   title,
				"content": lines,
			},
		},
	}

	datas, _ := json.Marshal(body)
	return datas
}

