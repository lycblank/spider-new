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
    GetPerDayProblem(n)
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

    idx := 0
    var title, href, solution, correctate, level string
    doc.Find("div").EachWithBreak(func(i int, s *goquery.Selection) bool {
        if s.AttrOr("role", "") != "cell" {
            return true
        }
        idx++
        if idx == 1 {
            return true
        }

        if idx == 2 {
            s := s.Find("a")
            link := fmt.Sprintf("https://leetcode-cn.com/%s", s.AttrOr("href", ""))
            href = fmt.Sprintf("[%s](%s)", link, link)
            title = s.Text()
        }
        if idx == 3 {
            solution = s.Find("a").Text()
        }
        if idx == 4 {
            correctate = s.Find("span").Text()
        }
        if idx == 5 {
            level = s.Find("span").Text()
        }
        return true
    })
    strs := strings.Fields(title)
    arg := notify.NotifyArg{
        Title: "每日一题",
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
        chromedp.WaitVisible(`#__next > div > div > div.grid.grid-cols-4.gap-4.md\:grid-cols-3.lg\:grid-cols-4.lg\:gap-6 > div.col-span-4.z-base.md\:col-span-2.lg\:col-span-3 > div:nth-child(7) > div.-mx-4.md\:mx-0 > div > div > div:nth-child(2) > div:nth-child(1)`),
        chromedp.OuterHTML(`#__next > div > div > div.grid.grid-cols-4.gap-4.md\:grid-cols-3.lg\:grid-cols-4.lg\:gap-6 > div.col-span-4.z-base.md\:col-span-2.lg\:col-span-3 > div:nth-child(7) > div.-mx-4.md\:mx-0 > div > div > div:nth-child(2) > div:nth-child(1)`, &str),
    )
    if err != nil {
        fmt.Println(err)
    }
    return str
}
