package costCompare

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

type JejuPassRequest struct {
	vhctySeCode string
	resveBeginDe string
	resveBeginTime string
	resveEndDe string
	resveEndTime string
	insrncApplcCode string
	monthYn string
	yearLmtYn string
	driverLicenseOneYearUnder string
}

// CostCompare
func CostCompare(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	contextVar, cancelFunc := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancelFunc()

	var strVar string
	costMaps := map[string]string{}
	resultString := ""

	contextVar, cancelFunc = context.WithTimeout(contextVar, 15*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://www.jejupassrent.com/home/search/list.do`),
		chromedp.Click(`#YMD button`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="` + start + `"]`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="` + end + `"]`),
		chromedp.Click(`.popover-content button.next.applyCalendarPanel`),
		// TODO: 시간
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

		resultString += model + " : " + cost + "\n"
		costMaps[model] = cost
	})

	c.String(200, resultString)
}
