package costCompare

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/shinYeongHyeon/rent-a-car-compare/costCompare/jejuPass"
	"os"
)

// CostCompare
func CostCompare(c *gin.Context) {
	start := c.Query("start")
	startTime := c.Query("sTime")
	end := c.Query("end")
	endTime := c.Query("eTime")

	results := jejuPass.JejuPass(start, startTime, end, endTime)

	fileName := start + startTime + "~" + end + endTime + ".csv"

	isSuccess := extractCsv(fileName, results)

	if isSuccess {
		c.FileAttachment(fileName, fileName)
		defer os.Remove(fileName)
	} else {
		c.String(500, "Fail")
	}
}

func extractCsv(fileName string, maps map[string]string) bool {
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Jejupass", ""}
	err = writer.Write(headers)

	for key, value := range maps {
		err = writer.Write([]string{key, value})
	}

	return true
}
