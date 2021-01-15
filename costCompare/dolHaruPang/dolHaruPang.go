package dolHaruPang

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/shinYeongHyeon/rent-a-car-compare/lib"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// DolHaruPang
func DolHaruPang(c *gin.Context) {
	startDate := c.Query("start")
	startTime := c.Query("sTime")
	endDate := c.Query("end")
	endTime := c.Query("eTime")

	zzimCarDaysMap := map[time.Weekday]int{}
	zzimCarDaysMap[time.Sunday] = 0
	zzimCarDaysMap[time.Monday] = 1
	zzimCarDaysMap[time.Tuesday] = 2
	zzimCarDaysMap[time.Wednesday] = 3
	zzimCarDaysMap[time.Thursday] = 4
	zzimCarDaysMap[time.Friday] = 5
	zzimCarDaysMap[time.Saturday] = 6

	zzimCarMonthsMap := map[time.Month]int{}
	zzimCarMonthsMap[time.January] = 1
	zzimCarMonthsMap[time.February] = 2
	zzimCarMonthsMap[time.March] = 3
	zzimCarMonthsMap[time.April] = 4
	zzimCarMonthsMap[time.May] = 5
	zzimCarMonthsMap[time.June] = 6
	zzimCarMonthsMap[time.July] = 7
	zzimCarMonthsMap[time.August] = 8
	zzimCarMonthsMap[time.September] = 9
	zzimCarMonthsMap[time.October] = 10
	zzimCarMonthsMap[time.November] = 11
	zzimCarMonthsMap[time.December] = 12

	startYear, _ := strconv.Atoi(startDate[:4])
	startMonth, _ := strconv.Atoi(startDate[4:6])
	startDay, _ := strconv.Atoi(startDate[6:8])
	startMonthTime := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	start := time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	startWeek := getWeekOfMonth(startDay, getSaturdayDateOfFirstWeekOfMonth(zzimCarDaysMap[startMonthTime.Weekday()]))
	startHourSpanChild := getSpanChildNumber(startTime)

	fmt.Println(startHourSpanChild)

	endYear, _ := strconv.Atoi(endDate[:4])
	endMonth, _ := strconv.Atoi(endDate[4:6])
	endDay, _ := strconv.Atoi(endDate[6:8])
	endMonthTime := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	endWeek := getWeekOfMonth(endDay, getSaturdayDateOfFirstWeekOfMonth(zzimCarDaysMap[endMonthTime.Weekday()]))
	endHourSpanChild := getSpanChildNumber(endTime)

	monthDiff := startMonth - zzimCarMonthsMap[time.Now().Month()]
	startAndEndMonthDiff := endMonth - startMonth

	opts := []chromedp.ExecAllocatorOption{
		chromedp.WindowSize(1920, 1080),
	}

	contextVar, cancelFunc := chromedp.NewExecAllocator(context.Background(), opts...)

	contextVar, cancelFunc = chromedp.NewContext(contextVar)

	defer cancelFunc()

	var strVar string
	dolHaruPangResult := map[string]string{}

	contextVar, cancelFunc = context.WithTimeout(contextVar, 60*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://www.dolharupang.com/list/car`),
		chromedp.Click(`input[name="startDate"]`),
	)

	errCheck(err)

	if monthDiff != 0 {
		for i := 0; monthDiff - i > 0; i++ {
			fmt.Println("i: ", i)
			errFor := chromedp.Run(contextVar,
				chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed th.next`,chromedp.ByQuery),
			)

			if errFor != nil {
				fmt.Println(err)
			}
		}
	}

	err3 := chromedp.Run(contextVar,
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed tbody tr:nth-child(` + strconv.Itoa(startWeek) + `) td:nth-child(` + strconv.Itoa(zzimCarDaysMap[start.Weekday()] + 1) + `)`),
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
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-days table.table-condensed tbody tr:nth-child(` + strconv.Itoa(endWeek) + `) td:nth-child(` + strconv.Itoa(zzimCarDaysMap[end.Weekday()] + 1) + `)`),
		chromedp.Click(`body div.datetimepicker[style*="block"] div.datetimepicker-hours table.table-condensed tbody tr td span:nth-child(` + strconv.Itoa(endHourSpanChild) + `)`),
		chromedp.Click(`.input-insur li:nth-child(2) button`),
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

func getSaturdayDateOfFirstWeekOfMonth(startDay int) int {
	return 7 - startDay
}

func getWeekOfMonth(date, saturdayDateOfFirstWeekOfMonth int) int {
	if saturdayDateOfFirstWeekOfMonth >= date {
		return 1
	}

	return 1 + int(math.Ceil(float64(date - saturdayDateOfFirstWeekOfMonth) / 7))
}

func getSpanChildNumber(time string) int {
	hour, _ := strconv.Atoi(time[:2])
	hourSpanChild := (hour - 8) * 2 + 1
	if time[2:4] == "30" {
		hourSpanChild += 1
	}

	return hourSpanChild
}