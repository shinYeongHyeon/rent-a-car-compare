package costCompare

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq"
	"log"
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
	contextVar, cancelFunc := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancelFunc()

	var strVar string
	var nodes []cdp.NodeID

	contextVar, cancelFunc = context.WithTimeout(contextVar, 15*time.Second)
	defer cancelFunc()

	err := chromedp.Run(contextVar,
		chromedp.Navigate(`https://www.jejupassrent.com/home/search/list.do`),
		chromedp.Click(`#YMD button`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="20210110"]`),
		chromedp.Click(`.popover-content #date table tbody td[data-num="20210111"]`),
		chromedp.Click(`.popover-content button.applyCalendarPanel`),
		// 시간 생략 ... 귀찮아..
		chromedp.Click(`#container .popover-content .btn-wrap button.next`),
		chromedp.NodeIDs(`#container #drawAreaView span.text-price`, &nodes, chromedp.BySearch),
		chromedp.Text(nodes, &strVar, chromedp.ByNodeID),
		//chromedp.Text(`#container #drawAreaView span.text-price`, &strVar, chromedp.ByQueryAll),
		//chromedp.Text(`td[data-num]="20210110"`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
		//chromedp.InnerHTML(`#popover399820`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
		//chromedp.InnerHTML(`#container .popover-content .btn-wrap button.next`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
		//chromedp.InnerHTML(`.popover-content td[data-num]="20210110"`, &strVar, chromedp.NodeVisible, chromedp.ByQuery),
	)

	if err != nil {
		fmt.Println("Err : ", err)
	}
	fmt.Println(strVar)


}

func getCarsAndCost(json *gojsonq.JSONQ) {

}