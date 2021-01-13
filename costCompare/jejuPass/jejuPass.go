package jejuPass

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"strings"
	"time"
	"log"
)

// JejuPass cars
func JejuPass(startDate, startTime, endDate, endTime string) map[string]string {
	contextVar, cancelFunc := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancelFunc()

	var strVar string
	jejuPassResult := map[string]string{}

	contextVar, cancelFunc = context.WithTimeout(contextVar, 60*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://www.jejupassrent.com/home/search/list.do`),
		chromedp.Click(`#YMD button`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="` + startDate + `"]`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="` + endDate + `"]`),
		chromedp.Click(`.popover-content button.next.applyCalendarPanel`),
		chromedp.SetValue(`.popover-content > div.pop-select-date > div.panelSearchTimeDiv > div.time > dl > dd > select.i-hour.clickStartTime`, startTime[:2]),
		chromedp.SetValue(`.popover-content > div.pop-select-date > div.panelSearchTimeDiv > div.time > dl > dd > select.i-minute.clickStartMin`, startTime[2:4]),
		chromedp.SetValue(`.popover-content > div.pop-select-date > div.panelSearchTimeDiv > div.time > dl > dd > select.i-hour.clickEndTime`, endTime[:2]),
		chromedp.SetValue(`.popover-content > div.pop-select-date > div.panelSearchTimeDiv > div.time > dl > dd > select.i-minute.clickEndMin`, endTime[2:4]),
		chromedp.Click(`#container .popover-content .btn-wrap button.next`),
		chromedp.Click(`#container > div.content-wrap > div.sidebar > form#leftForm > div.search-car-input > div.btn-wrap > #searchBtn`),
		chromedp.WaitVisible(`#content #drawAreaView`),
		chromedp.InnerHTML(`#content #drawAreaView`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
	)

	if err != nil {
		fmt.Println("Err : ", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(strVar))

	if err != nil {
		fmt.Println("Document Err : ", err)
	}

	doc.Find("li.panel").Each(func (i int, s *goquery.Selection) {
		model := s.Find("div.title > h3").Text()
		cost := s.Find("div.right span.text-price").Text()

		jejuPassResult[model] = cost
	})

	return jejuPassResult
}
