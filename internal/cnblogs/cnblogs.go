package cnblogs

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

var targetUrl string = `https://www.cnblogs.com/`
var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n

	s := gocron.NewScheduler()
	s.Every(1).Day().At("19:30").Do(GetDailyCNBlogs, n)
	s.Start()
}

func GetDailyCNBlogs(n notify.Notify) {
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

	contents := make([]string, 0, 16)
	doc.Find(".post-item-title").Each(func(_ int, s *goquery.Selection) {
		url := s.AttrOr("href", "")
		title := s.Text()
		contents = append(contents, fmt.Sprintf("%s[%s]", title, url))
	})

	arg := notify.NotifyArg{
		Title: "博客园一览",
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
		chromedp.WaitVisible(`#post_list`),
		chromedp.OuterHTML(`#post_list`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
