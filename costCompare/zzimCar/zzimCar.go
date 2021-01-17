package zzimCar

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/shinYeongHyeon/go-times"
	"github.com/shinYeongHyeon/rent-a-car-compare/lib"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// ZzimCar cars
func ZzimCar(c *gin.Context) {
	startDate := c.Query("start")
	startTime := c.Query("sTime")
	endDate := c.Query("end")
	endTime := c.Query("eTime")

	startYear, _ := strconv.Atoi(startDate[:4])
	startMonth, _ := strconv.Atoi(startDate[4:6])
	startDay, _ := strconv.Atoi(startDate[6:8])
	start := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	startTimeHourAndMinute := startTime[:2] + ":" + startTime[2:4]
	startWeek := times.GetNthWeekOfMonth(start)

	endYear, _ := strconv.Atoi(endDate[:4])
	endMonth, _ := strconv.Atoi(endDate[4:6])
	endDay, _ := strconv.Atoi(endDate[6:8])
	end := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	endTimeHourAndMinute := endTime[:2] + ":" + endTime[2:4]
	endWeek := times.GetNthWeekOfMonth(end)

	monthDiff := startMonth - times.MonthMap[time.Now().Month()]
	startAndEndMonthDiff := endMonth - startMonth

	contextVar, cancelFunc := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancelFunc()

	var strVar string
	var picbuf []byte
	zzimCarResults := map[string]string{}

	contextVar, cancelFunc = context.WithTimeout(contextVar, 60*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://zzimcar.co.kr/`),
		chromedp.Click(`div.cookiePopup div.popup-area a.popup-close`),
		chromedp.Sleep(2 * time.Second),
	)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}

	if monthDiff != 0 {
		for i := 0; monthDiff - i > 0; i++ {
			errFor := chromedp.Run(contextVar, chromedp.Click("ui-datepicker-next"))

			if errFor != nil {
				fmt.Println(err)
			}
		}
	}

	err2 := chromedp.Run(contextVar,
		chromedp.Click("div.dual-input div.stime div.nice-select"),
		chromedp.Click(`div.dual-input div.stime div.nice-select ul li[data-value="` + startTimeHourAndMinute + `"]`),
		chromedp.Click("div.dual-input div.etime div.nice-select"),
		chromedp.Click(`div.dual-input div.etime div.nice-select ul li[data-value="` + endTimeHourAndMinute + `"]`),
		chromedp.Click(".ui-datepicker-calendar tbody tr:nth-child(" + strconv.Itoa(startWeek) + ") td:nth-child(" + strconv.Itoa(times.DaysMap[start.Weekday()] + 1) + ")"),
		)


	if err2 != nil {
		fmt.Println("error")
		fmt.Println(err2)
	}

	if startAndEndMonthDiff != 0 {
		chromedp.Run(contextVar, chromedp.Click("ui-datepicker-next"))
	}

	err3 := chromedp.Run(contextVar,
		chromedp.Click(".ui-datepicker-calendar tbody tr:nth-child(" + strconv.Itoa(endWeek) + ") td:nth-child(" + strconv.Itoa(times.DaysMap[end.Weekday()] + 1) + ")"),
		chromedp.Click("section.section-main article.search-area div.box-step2 div.select-type div.nice-select.selected"),
		chromedp.Click(`section.section-main article.search-area div.box-step2 div.select-type div.nice-select.selected ul .option[data-value="JEJU"]`),
		chromedp.Click("section.section-main article.search-area div.box-step3 button", chromedp.ByQuery),
		chromedp.Sleep(10 * time.Second),
		chromedp.InnerHTML(`div.wrap-search div.map-root section.section-search article.result-area div.box-wrap div.box-result-list ul`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			picbuf, err = page.CaptureScreenshot().
				WithQuality(90).
				WithClip(&page.Viewport{
					X: contentSize.X,
					Y: contentSize.Y,
					Width: contentSize.Width,
					Height: contentSize.Height,
					Scale: 1,
				}).Do(ctx)

			if err != nil {
				return err
			}

			return nil
		}),
	)

	if err3 != nil {
		fmt.Println("error")
		fmt.Println(err3)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(strVar))

	if err != nil {
		fmt.Println("Document Err : ", err)
	}

	doc.Find("li").Each(func (i int, s *goquery.Selection) {
		box := s.Find("dl.car-info-box dd")
		model := box.Find("div.title-box").Text()
		cost := strings.Split(box.Find("div.price-box p").Text(), " ~")[0]

		zzimCarResults[model] = cost
	})

	ioutil.WriteFile("ds.png", picbuf, 0o644)

	zzimCarFileName := startDate + startTime + "~" + endDate + endTime + "_zzimCar.csv"
	isSuccess := lib.ExtractCsv("zzimCar", zzimCarFileName, zzimCarResults)

	if isSuccess {
		c.FileAttachment(zzimCarFileName, zzimCarFileName)

		defer os.Remove(zzimCarFileName)
	} else {
		c.String(500, "Fail")
	}
}