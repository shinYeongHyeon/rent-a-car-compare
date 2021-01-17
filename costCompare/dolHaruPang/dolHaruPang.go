package dolHaruPang

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/shinYeongHyeon/go-times"
	"github.com/shinYeongHyeon/rent-a-car-compare/lib"
	"os"
	"strconv"
	"log"
	"strings"
	"time"
)

// DolHaruPang
func DolHaruPang(c *gin.Context) {
	startDate := c.Query("start")
	startTime := c.Query("sTime")
	endDate := c.Query("end")
	endTime := c.Query("eTime")

	startYear, _ := strconv.Atoi(startDate[:4])
	startMonth, _ := strconv.Atoi(startDate[4:6])
	startDay, _ := strconv.Atoi(startDate[6:8])
	start := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	startWeek := times.GetNthWeekOfMonth(start)
	startHourSpanChild := getSpanChildNumber(startTime)

	endYear, _ := strconv.Atoi(endDate[:4])
	endMonth, _ := strconv.Atoi(endDate[4:6])
	endDay, _ := strconv.Atoi(endDate[6:8])
	end := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	endWeek := times.GetNthWeekOfMonth(end)
	endHourSpanChild := getSpanChildNumber(endTime)

	monthDiff := startMonth - times.MonthMap[time.Now().Month()]
	startAndEndMonthDiff := endMonth - startMonth

	contextVar, cancelFunc := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancelFunc()

	var strVar string
	dolHaruPangResult := map[string]string{}

	contextVar, cancelFunc = context.WithTimeout(contextVar, 60*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://www.dolharupang.com/list/car`),
		chromedp.RemoveAttribute(`div.mobile-menu.visible-xs`, "class"),
		chromedp.Click(`input[name="startDate"]`),
		chromedp.Click(`input[name="startDate"]`),
	)

	errCheck(err)

	if monthDiff != 0 {
		for i := 0; monthDiff - i > 0; i++ {
			errFor := chromedp.Run(contextVar,
				chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed th.next`,chromedp.ByQuery),
			)

			if errFor != nil {
				fmt.Println(err)
			}
		}
	}

	err3 := chromedp.Run(contextVar,
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed tbody tr:nth-child(` + strconv.Itoa(startWeek) + `) td:nth-child(` + strconv.Itoa(times.DaysMap[start.Weekday()] + 1) + `)`),
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-hours table.table-condensed tbody tr td span:nth-child(` + strconv.Itoa(startHourSpanChild) + `)`),
		chromedp.Click(`input[name="endDate"]`),
	)

	errCheck(err3)

	if startAndEndMonthDiff != 0 {
		for i := 0; startAndEndMonthDiff - i > 0; i++ {
			fmt.Println("i: ", i)
			errFor := chromedp.Run(contextVar,
				chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed th.next`,chromedp.ByQuery),
			)

			if errFor != nil {
				fmt.Println(err)
			}
		}
	}

	err4 := chromedp.Run(contextVar,
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed tbody tr:nth-child(` + strconv.Itoa(endWeek) + `) td:nth-child(` + strconv.Itoa(times.DaysMap[end.Weekday()] + 1) + `)`),
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-hours table.table-condensed tbody tr td span:nth-child(` + strconv.Itoa(endHourSpanChild) + `)`),
		chromedp.Sleep(150 * time.Millisecond),
		chromedp.Click(`.input-insur li:nth-child(2) button`),
		chromedp.Sleep(150 * time.Millisecond),
		chromedp.Click("#searchSchedule div.block-btn-search button.btn"),
		chromedp.Sleep(8 * time.Second),
	)

	errCheck(err4)

	for i := 0; i < 100; i++ {
		errLoop := chromedp.Run(contextVar,
			chromedp.ScrollIntoView(".bottom-corp-info"),
			chromedp.Sleep(150 * time.Millisecond),
		)

		errCheck(errLoop)
	}

	debugErr := chromedp.Run(contextVar,
		chromedp.InnerHTML("#ajaxListCar", &strVar),
	)

	errCheck(debugErr)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(strVar))
	if err != nil {
		fmt.Println("Document Err : ", err)
	}

	doc.Find("div.block-car.original.ok").Each(func (i int, s *goquery.Selection) {
		model, _ := s.Attr("data-product-nm")
		cost := s.Find("div.inline.section-right").Find("div.part-bottom").Find("div.inline.price-total.ani-popping-down").Find("span").Text()

		dolHaruPangResult[model] = cost
	})

	dolHaruPangFileName := startDate + startTime + "~" + endDate + endTime + "_dolHaruPang.csv"
	isSuccess := lib.ExtractCsv("dolHaruPang", dolHaruPangFileName, dolHaruPangResult)

	if isSuccess {
		c.FileAttachment(dolHaruPangFileName, dolHaruPangFileName)

		defer os.Remove(dolHaruPangFileName)
	} else {
		c.String(500, "Fail")
	}
}

func errCheck(err error) {
	if err != nil {
		fmt.Println("err : ", err)
	}
}

func getSpanChildNumber(time string) int {
	hour, _ := strconv.Atoi(time[:2])
	hourSpanChild := (hour - 8) * 2 + 1
	if time[2:4] == "30" {
		hourSpanChild += 1
	}

	return hourSpanChild
}