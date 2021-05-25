package pushplus

import (
	"context"
	"fmt"
	"github.com/lycblank/spider-new/pkg/notify"
	"github.com/parnurzeal/gorequest"
)


type PushPlus struct {
	ch       chan notify.NotifyArg
	webhook  string
	topic string
	token string
}

func NewPushPlus(webhook string, topic string, token string) *PushPlus {
	c := &PushPlus{
		ch:       make(chan notify.NotifyArg, 512),
		webhook:  webhook,
		topic: topic,
		token: token,
	}
	go c.process()
	return c
}

func (c *PushPlus) Send(ctx context.Context, arg notify.NotifyArg) error {
	select {
	case c.ch <- arg:
	default:
	}
	return nil
}

func (c *PushPlus) Alert(title string, contents ...string) {
	if len(contents) <= 0 {
		contents = append(contents, title)
	}
	select {
	case c.ch <- notify.NotifyArg{Title: title, Contents: contents}:
	default:
	}
}

func (c *PushPlus) process() {
	for arg := range c.ch {
		c.sendMsg(arg)
	}
}

func (c *PushPlus) sendMsg(arg notify.NotifyArg) {
	content := map[string]interface{}{
		"token": c.token,
		"title": arg.Title,
		"template": "json",
		"topic": c.topic,
	}
	resp, _, errs := gorequest.New().Post(c.webhook).
		Set(`Content-Type`, `application/json`).
		Send(content).End()
	if len(errs) > 0 {
		fmt.Printf("send alert failed. errs:%+v\n", errs)
		return
	}
	resp.Body.Close()
	return
}