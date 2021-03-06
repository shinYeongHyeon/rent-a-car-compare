package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shinYeongHyeon/rent-a-car-compare/costCompare/dolHaruPang"
	"github.com/shinYeongHyeon/rent-a-car-compare/costCompare/jejuPass"
	"github.com/shinYeongHyeon/rent-a-car-compare/costCompare/zzimCar"
	"github.com/shinYeongHyeon/rent-a-car-compare/health"
	"time"
)

func main() {
	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string { // 커스텀 로그 (아파치에서 출력하는 형식)
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("template/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.GET("/ping", health.Health)
	r.GET("/jejuPass", jejuPass.JejuPass)
	r.GET("/zzimCar", zzimCar.ZzimCar)
	r.GET("/dolHaruPang", dolHaruPang.DolHaruPang)

	r.Run(":3030")
}


