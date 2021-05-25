package main

import (
	"context"
	"github.com/lycblank/spider-new/internal/conf"
	"github.com/lycblank/spider-new/internal/github"
	"github.com/lycblank/spider-new/internal/gonews"
	"github.com/lycblank/spider-new/internal/iciba"
	"github.com/lycblank/spider-new/internal/leetcode"
	"github.com/lycblank/spider-new/internal/toutiao"
	"github.com/lycblank/spider-new/pkg/chanify"
	"github.com/lycblank/spider-new/pkg/flybook"
	"github.com/lycblank/spider-new/pkg/notify"
	"github.com/lycblank/spider-new/pkg/pushplus"
	"time"
)

func main() {
	config := conf.GetConfig()
	ic := &NotifyContainer{}
	ic.AddNotify(flybook.NewFlyBook(config.FlyBook.Webhook))
	ic.AddNotify(chanify.NewChanify(config.Chanify.Webhook))
	ic.AddNotify(pushplus.NewPushPlus(config.PushPlus.Webhook, config.PushPlus.Group, config.PushPlus.Token))
	iciba.Init(ic)
	time.Sleep(10*time.Second)
	github.Init(ic)
	time.Sleep(10*time.Second)
	gonews.Init(ic)
	time.Sleep(10*time.Second)
	leetcode.Init(ic)
	time.Sleep(10*time.Second)
	toutiao.Init(ic)
	select {}
}

type NotifyContainer struct {
	ns []notify.Notify
}

func (ic *NotifyContainer) AddNotify(n notify.Notify) {
	ic.ns = append(ic.ns, n)
}

func (ic *NotifyContainer) Send(ctx context.Context, arg notify.NotifyArg) error {
	for _, n := range ic.ns {
		n.Send(ctx, arg)
	}
	return nil
}


