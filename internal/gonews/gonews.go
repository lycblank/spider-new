package gonews

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"github.com/lycblank/spider-new/pkg/notify"
	"net/http"
)

var baseAddr = `https://gocn.vip`
var listAddr = `/topics/node18`

var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n
	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:30").Do(Fetch, n)
	s.Start()
}

func Fetch(n notify.Notify) {
	listUrl := fmt.Sprintf("%s%s", baseAddr, listAddr)
	fmt.Println(listUrl)
	httpClient := &http.Client{}
	resp, err := httpClient.Get(listUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("new document failed. err:%+v\n", err)
		return
	}

	first := doc.Find(".row .topics .topic .title a")

	topicUrl, ok := first.Attr("href")
	if !ok {
		fmt.Printf("attr href not exits dddd %s\n", first.Text())
		return
	}
	contentUrl := fmt.Sprintf("%s%s", baseAddr, topicUrl)
	FetchContent(contentUrl, n)
}

func FetchContent(contentUrl string, n notify.Notify) {
	http := &http.Client{}
	resp, err := http.Get(contentUrl)
	if err != nil {
		fmt.Printf("get %s error err:%+v\n", contentUrl, err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("new document err. err:%+v\n", err)
		return
	}

	// 获取标题
	titleSel := doc.Find(".row .topic-detail .title")
	title := titleSel.Text()
	contents := make([]string, 0, 8)
	doc.Find(".card-body ol li").Each(func(i int, s *goquery.Selection){
		contents = append(contents, s.Text())
	})

	n.Send(context.Background(), notify.NotifyArg{
		Title: title,
		Contents: contents,
	})
}