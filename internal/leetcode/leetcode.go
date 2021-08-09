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

	var title, href, solution, correctate, level string
	doc.Find("tbody tr").EachWithBreak(func(i int, s *goquery.Selection)bool{
		if i == 1 {
			s.Find("td").EachWithBreak(func(k int, s *goquery.Selection) bool {
				if k == 1 {
					t := s.Find("a")
					title = t.Text()
					href = fmt.Sprintf("%s%s", "https://leetcode-cn.com", t.AttrOr("href", ""))
				}
				if k == 2 {
					t := s.Find("a")
					solution = t.Text()
				}
				if k == 3 {
					t := s.Find("span")
					correctate = t.Text()
				}
				if k == 4 {
					t := s.Find("span")
					level = t.Text()
					return false
				}
				return true
			})
			return false
		}
		return true
	})
	strs := strings.Fields(title)
	arg := notify.NotifyArg{
		Title:"每日一题",
		Contents: []string{
			"序号: " + strs[0][:len(strs[0])-1],
			"标题: " + strs[1],
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
		chromedp.WaitVisible(`#__next > div > div > div.grid.grid-cols-4.md\:grid-cols-3.lg\:grid-cols-4.gap-4.lg\:gap-6 > div.col-span-4.md\:col-span-2.lg\:col-span-3 > div.jsx-3812067982 > div.ant-table-wrapper.question-table.-mx-4.md\:mx-0 > div > div > div > div > div > table`),
		chromedp.OuterHTML(`#__next > div > div > div.grid.grid-cols-4.md\:grid-cols-3.lg\:grid-cols-4.gap-4.lg\:gap-6 > div.col-span-4.md\:col-span-2.lg\:col-span-3 > div.jsx-3812067982 > div.ant-table-wrapper.question-table.-mx-4.md\:mx-0 > div > div > div > div > div > table`, &str),
	)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
