package main

import (
	"github.com/lycblank/spider-new/internal/conf"
	"github.com/lycblank/spider-new/internal/iciba"
	"github.com/lycblank/spider-new/pkg/flybook"
)

func main() {
	fb := flybook.NewFlyBook(conf.GetConfig().FlyBook.Webhook)
	iciba.Init(fb)
	select {}
}



