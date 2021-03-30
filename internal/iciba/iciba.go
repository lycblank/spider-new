// 每日一句 英语
package iciba

import (
	"context"
	"fmt"
	"io"
	"github.com/jasonlvhit/gocron"
	"github.com/chromedp/chromedp"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

var targetUrl string = `http://news.iciba.com/`
var writer io.Writer
func Init(w io.Writer) {
	writer = w

	GetEnginishAndChinese(w)

	s := gocron.NewScheduler()
	s.Every(1).Day().At("09:00").Do(GetEnginishAndChinese, writer)
	s.Start()
}

func GetEnginishAndChinese(w io.Writer) {
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

	content := fmt.Sprintf("%s\n%s\n%s", "每日一句", enginish, chinese)
	io.WriteString(w, content)
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
