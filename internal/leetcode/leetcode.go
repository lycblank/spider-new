package leetcode

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

var targetUrl string = `https://leetcode-cn.com/problemset/all/`
var defaultNotify notify.Notify
func Init(n notify.Notify) {
	defaultNotify = n
	s := gocron.NewScheduler()
	s.Every(1).Day().At("08:00").Do(GetPerDayProblem, n)
	s.Start()
}

func GetPerDayProblem(n notify.Notify) {
	htmlContent := getHtmlContent()
	if htmlContent == "" {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println(err)
		return
	}

	var title, href, seq, solution, correctate, level string
	doc.Find(".question-title").EachWithBreak(func(_ int, s *goquery.Selection)bool{
		s = s.Find("a")
		title = s.Text()
		href = fmt.Sprintf("%s%s", "https://leetcode-cn.com", s.AttrOr("href", ""))
		return false
	})

	doc.Find("tbody tr td").EachWithBreak(func(i int, s *goquery.Selection)bool{
		if i == 1 {
			seq = s.Text()
		}
		if i == 3 {
			solution = s.Text()
		}
		if i == 4 {
			correctate = s.Text()
		}
		if i == 5 {
			level = s.Text()
			return false
		}
		return true
	})

	arg := notify.NotifyArg{
		Title:"每日一题",
		Contents: []string{
			"序号: " + seq,
			"标题: " + title,
			"链接: " + href,
			"题解: " + solution,
			"通过率:" + correctate,
			"难度: " + level,
		},
	}
	n.Send(context.Background(), arg)
}

func getHtmlContent() string {
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	var str string
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetUrl),
		chromedp.WaitVisible(`#question-app > div > div:nth-child(2) > div.question-list-base > div.table-responsive.question-list-table > table`),
		chromedp.OuterHTML(`#question-app > div > div:nth-child(2) > div.question-list-base > div.table-responsive.question-list-table > table`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
