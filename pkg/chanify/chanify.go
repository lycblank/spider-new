package chanify

import (
	"context"
	"fmt"
	"github.com/lycblank/spider-new/pkg/notify"
	"github.com/parnurzeal/gorequest"
	"strings"
)


type Chanify struct {
	ch       chan notify.NotifyArg
	webhook  string
}

func NewChanify(webhook string) *Chanify {
	c := &Chanify{
		ch:       make(chan notify.NotifyArg, 512),
		webhook:  webhook,
	}
	go c.process()
	return c
}

func (c *Chanify) Send(ctx context.Context, arg notify.NotifyArg) error {
	select {
	case c.ch <- arg:
	default:
	}
	return nil
}

func (c *Chanify) Alert(title string, contents ...string) {
	if len(contents) <= 0 {
		contents = append(contents, title)
	}
	select {
	case c.ch <- notify.NotifyArg{Title: title, Contents: contents}:
	default:
	}
}

func (c *Chanify) process() {
	for arg := range c.ch {
		c.sendMsg(arg)
	}
}

func (c *Chanify) sendMsg(arg notify.NotifyArg) {
	resp, _, errs := gorequest.New().Post(c.webhook + "?sound=1&title=" + arg.Title).
		Set(`Content-Type`, `text/plain`).
		SendString(strings.Join(arg.Contents, "\n")).End()
	if len(errs) > 0 {
		fmt.Printf("send alert failed. errs:%+v\n", errs)
		return
	}
	resp.Body.Close()
	return
}