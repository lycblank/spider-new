package toutiao

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

var targetUrl string = `https://toutiao.io/`
var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n

	GetDailyToutiao(n)

	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:25").Do(GetDailyToutiao, n)
	s.Start()
}

func GetDailyToutiao(n notify.Notify) {
	htmlContent := getHtmlContent()
	if htmlContent == "" {
		return
	}
	fmt.Println(htmlContent)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println(err)
		return
	}

	var title string
	doc.Find(".date").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		s.Find("span").Text()
		title = fmt.Sprintf("%s(%s)", s.Find("span").Text(), s.Find("small").Text())
		return false
	})

	contents := make([]string, 0, 16)
	doc.Find("div.content > h3 > a").Each(func(_ int, s *goquery.Selection){
		url := targetUrl + s.AttrOr("href", "")
		title := s.AttrOr("title", "")
		contents = append(contents, fmt.Sprintf("%s[%s]", title, url))
	})

	arg := notify.NotifyArg{
		Title: title,
		Contents: contents,
	}
	n.Send(context.Background(), arg)
}

func getHtmlContent() string {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible(`#daily > div > div.daily`),
		chromedp.OuterHTML(`#daily > div > div.daily`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
