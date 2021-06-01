package main

import (
	"context"
	"github.com/lycblank/spider-new/internal/cnblogs"
	"github.com/lycblank/spider-new/internal/conf"
	"github.com/lycblank/spider-new/internal/github"
	"github.com/lycblank/spider-new/internal/gonews"
	"github.com/lycblank/spider-new/internal/iciba"
	"github.com/lycblank/spider-new/internal/leetcode"
	"github.com/lycblank/spider-new/internal/shequ"
	"github.com/lycblank/spider-new/internal/toutiao"
	"github.com/lycblank/spider-new/pkg/chanify"
	"github.com/lycblank/spider-new/pkg/flybook"
	"github.com/lycblank/spider-new/pkg/notify"
	"github.com/lycblank/spider-new/pkg/pushplus"
)

func main() {
	config := conf.GetConfig()
	ic := &NotifyContainer{}
	ic.AddNotify(flybook.NewFlyBook(config.FlyBook.Webhook))
	ic.AddNotify(chanify.NewChanify(config.Chanify.Webhook))
	ic.AddNotify(pushplus.NewPushPlus(config.PushPlus.Webhook, config.PushPlus.Group, config.PushPlus.Token))
	iciba.Init(ic)
	github.Init(ic)
	gonews.Init(ic)
	leetcode.Init(ic)
	toutiao.Init(ic)
	shequ.Init(ic)
	cnblogs.Init(ic)
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


