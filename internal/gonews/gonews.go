package gonews

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/jasonlvhit/gocron"
	"github.com/lycblank/spider-new/pkg/notify"
	"log"
	"strings"
)

var baseAddr = `https://gocn.vip`
var listAddr = `/topics/cate/18?page=1&grade=hot`

var defaultNotify notify.Notify

func Init(n notify.Notify) {
	defaultNotify = n
	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:30").Do(Fetch, n)
	s.Start()
	Fetch(n)
}

func Fetch(n notify.Notify) {
	body := getHtmlFetch()
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		fmt.Printf("new document err. err:%+v\n", err)
		return
	}

	first := doc.Find("a")
	topicUrl, ok := first.Attr("href")
	if !ok {
		fmt.Printf("attr href not exits dddd %s\n", first.Text())
		return
	}
	fmt.Println(topicUrl)
	contentUrl := fmt.Sprintf("%s%s", baseAddr, topicUrl)
	FetchContent(contentUrl, n)
}

func FetchContent(contentUrl string, n notify.Notify) {
	body := getHtmlFetchContent(contentUrl)
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		fmt.Printf("new document err. err:%+v\n", err)
		return
	}

	// 获取标题
	titleSel := doc.Find("span")
	title := strings.Split(titleSel.Text(), "···")[0]
	contents := make([]string, 0, 8)
	doc.Find("ol li").Each(func(i int, s *goquery.Selection) {
		content := strings.TrimSpace(strings.Split(s.Text(), "https")[0])
		href := s.Find("a").AttrOr("href", "")
		contents = append(contents, fmt.Sprintf("[%s](%s)", content, href))
	})
	fmt.Println("content", contents)
	n.Send(context.Background(), notify.NotifyArg{
		Title:    title,
		Contents: contents,
	})
}

func getHtmlFetch() string {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("%s%s", baseAddr, listAddr)),
		chromedp.WaitVisible(`#root > div > section > div > main > div > div.ant-row.ant-row-center.ant-row-top > div:nth-child(1) > div > div > div.ant-spin-nested-loading > div > div > div > div.ant-spin-nested-loading > div > ul > div:nth-child(1) > div > li > div > div.ant-list-item-meta-content`),
		chromedp.OuterHTML(`#root > div > section > div > main > div > div.ant-row.ant-row-center.ant-row-top > div:nth-child(1) > div > div > div.ant-spin-nested-loading > div > div > div > div.ant-spin-nested-loading > div > ul > div:nth-child(1) > div > li > div > div.ant-list-item-meta-content`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}

func getHtmlFetchContent(url string) string {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#root > div > section > div > main > div > div > div > div:nth-child(1) > div > div.ant-card-body > div:nth-child(3) > div > div > ol`),
		chromedp.OuterHTML(`#root > div > section > div > main > div > div > div > div:nth-child(1) > div`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
