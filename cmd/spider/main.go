package main

import (
	"context"
	"github.com/lycblank/spider-new/internal/conf"
	"github.com/lycblank/spider-new/internal/github"
	"github.com/lycblank/spider-new/internal/gonews"
	"github.com/lycblank/spider-new/internal/iciba"
	"github.com/lycblank/spider-new/pkg/chanify"
	"github.com/lycblank/spider-new/pkg/flybook"
	"github.com/lycblank/spider-new/pkg/notify"
	"time"
)

func main() {
	ic := &NotifyContainer{}
	ic.AddNotify(flybook.NewFlyBook(conf.GetConfig().FlyBook.Webhook))
	ic.AddNotify(chanify.NewChanify(conf.GetConfig().Chanify.Webhook))
	iciba.Init(ic)
	time.Sleep(10*time.Second)
	github.Init(ic)
	time.Sleep(10*time.Second)
	gonews.Init(ic)
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


