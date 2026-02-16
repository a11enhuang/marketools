package marketdata

import (
	"context"
	"log"
	"strings"
	"time"

	"com.reopenai/marketool/dal"
	"com.reopenai/marketool/schedules"
	"com.reopenai/marketool/service/dingding"
)

func Run(ctx context.Context) {
	schedules.AddWeekdaysTask("09:00", buyStocks)
	schedules.AddWeekdaysTask("09:30", buyStocks)
	schedules.AddWeekdaysTask("10:00", buyStocks)
	schedules.AddWeekdaysTask("11:45", buyStocks)
	schedules.AddWeekdaysTask("13:20", buyStocks)
	schedules.AddWeekdaysTask("14:30", buyStocks)
	schedules.AddWeekdaysTask("14:45", buyStocks)
}

func buyStocks() {
	title := time.Now().Format("2006-01-02")
	entries := dal.SelectBuyStocks()
	if len(entries) > 0 {
		builder := strings.Builder{}
		builder.WriteString("# ")
		builder.WriteString(title)
		builder.WriteString(" 评级买入股票 \n")

		for i := range entries {
			entry := entries[i]
			builder.WriteString(" - ")
			builder.WriteString(entry.Name)
			builder.WriteRune('(')
			builder.WriteString(entry.Code)
			builder.WriteRune(')')
			builder.WriteString("最新价:")
			builder.WriteString(entry.Price.String())
			builder.WriteString("涨跌幅:")
			builder.WriteString(entry.Zdf.String())
			builder.WriteString(",换手率:")
			builder.WriteString(entry.Hsl.String())
			builder.WriteString(",量比:")
			builder.WriteString(entry.Lb.String())
			builder.WriteString(" \n")
		}

		dingding.Send(builder.String())
	} else {
		log.Println("未匹配到任何建议买入的股票")
	}
}
