// 每日一句 英语
package iciba

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/jasonlvhit/gocron"
	"github.com/lycblank/spider-new/pkg/notify"
	"log"
	"strings"
)

var targetUrl string = `http://news.iciba.com/`
var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n

	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:00").Do(GetEnginishAndChinese, n)
	s.Start()
}

func GetEnginishAndChinese(n notify.Notify) {
	htmlContent := getHtmlContent()
	if htmlContent == "" {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println(err)
	}
	var enginish, chinese string
	if s := doc.Find(".swiper-container-place .swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active .item.item-big .item-bottom .english"); s != nil {
		enginish = s.Text()
	}
	if s := doc.Find(".swiper-container-place .swiper-slide.swiper-slide-0.swiper-slide-visible.swiper-slide-active .item.item-big .item-bottom .chinese"); s != nil {
		chinese = s.Text()
	}

	arg := notify.NotifyArg{
		Title:"每日一句",
		Contents: []string{enginish, chinese},
	}
	n.Send(context.Background(), arg)
}

func getHtmlContent() string {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible(`.banner`),
		chromedp.OuterHTML(`.banner`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
