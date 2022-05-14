package gopherdaily

import (
    "bytes"
    "context"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "github.com/chromedp/chromedp"
    "github.com/jasonlvhit/gocron"
    "github.com/lycblank/spider-new/pkg/notify"
    "log"
)

var addr = `https://gopher-daily.com/`

var defaultNotify notify.Notify

func Init(n notify.Notify) {
    defaultNotify = n
    s := gocron.NewScheduler()
    s.Every(1).Day().At("08:30").Do(Fetch, n)
    s.Start()
}

func Fetch(n notify.Notify) {
    body := getHtmlFetch()
    fmt.Println("body = ", body)
    doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
    if err != nil {
        fmt.Printf("new document err. err:%+v\n", err)
        return
    }

    contents := make([]string, 0, 16)
    title := doc.Find("h2").Text()
    doc.Find("ol li").Each(func(i int, s *goquery.Selection) {
        content := s.Text()
        href := s.Find("a").AttrOr("href", "")
        contents = append(contents, fmt.Sprintf("[%s](%s)", content, href))
    })

    n.Send(context.Background(), notify.NotifyArg{
        Title:    title,
        Contents: contents,
    })
    fmt.Println(title)
    fmt.Println(contents)
}

func getHtmlFetch() string {
    ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
    defer cancel()

    var str string
    err := chromedp.Run(ctx,
        chromedp.Navigate(addr),
        chromedp.WaitVisible(`body > div > div > div.col.pt-3.pt-sm-2.pt-lg-3.pt-xl-5 > div > div > div.issue-content > ol`),
        chromedp.OuterHTML(`body > div > div > div.col.pt-3.pt-sm-2.pt-lg-3.pt-xl-5 > div > div > div.issue-content`, &str),
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
        chromedp.OuterHTML(`#root > div > section > div > main > div > div > div > div:nth-child(1) > div > div.ant-card-body > div:nth-child(3) > div > div`, &str),
    )
    if err != nil {
        fmt.Println(err)
    }
    return str
}
