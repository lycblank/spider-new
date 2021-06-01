package shequ

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
	"github.com/lycblank/spider-new/pkg/notify"
	"net/http"
	"strings"
	"time"
)

var targetUrl string = `https://studygolang.com/go/godaily`
var animUrl string = `https://studygolang.com`
var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n

	GetDailyShequ(n)
	s := gocron.NewScheduler()
	s.Every(1).Day().At("9:10").Do(GetDailyShequ, n)
	s.Start()
}

func GetDailyShequ(n notify.Notify) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(targetUrl)
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

	nowTime := time.Now().Unix()
	contents := make([]string, 0, 16)
	title := doc.Find("#wrapper > div > div > div.col-md-9.col-sm-6 > div:nth-child(3) > div > div.title > h2").Text()
	doc.Find("#wrapper > div > div > div.col-md-9.col-sm-6 .box_white .topics .topic").Each(func(_ int, s *goquery.Selection)  {
		title := strings.ReplaceAll(strings.ReplaceAll(s.Find(".right-info .title").Text(), "\t", ""), "\n", "")
		targetUrl := animUrl + s.Find(".right-info .title a").AttrOr("href", "")
		html, _ := s.Find("div > dd > div.meta > span").Html()
		sendTime := strings.Split(html, "\"")[3]
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t , _ :=time.ParseInLocation("2006-01-02 15:04:05", sendTime,loc)
		timeU := t.Unix()
		if nowTime < timeU+24*60*60 {
			contents = append(contents, fmt.Sprintf("%s[%s]", title, targetUrl))
		}
	})


	arg := notify.NotifyArg{
		Title: title,
		Contents: contents,
	}
	n.Send(context.Background(), arg)
}
