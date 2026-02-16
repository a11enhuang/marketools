package marketdata

import (
	"context"
	"fmt"
	"log"
	"time"

	"com.reopenai/marketool/schedules"
	"github.com/bytedance/sonic"
)

func startSyncLastPrice() {
	syncDataTask := func() {
		syncLastPrice(time.Now())
	}

	go syncDataTask()

	schedules.AddWeekdaysTask("09:20", syncDataTask)
	schedules.AddWeekdaysTask("09:50", syncDataTask)
	schedules.AddWeekdaysTask("11:00", syncDataTask)
	schedules.AddWeekdaysTask("11:30", syncDataTask)
	schedules.AddWeekdaysTask("13:10", syncDataTask)
	schedules.AddWeekdaysTask("13:30", syncDataTask)
	schedules.AddWeekdaysTask("14:00", syncDataTask)
	schedules.AddWeekdaysTask("14:20", syncDataTask)
	schedules.AddWeekdaysTask("14:35", syncDataTask)
	schedules.AddWeekdaysTask("14:40", syncDataTask)
	schedules.AddWeekdaysTask("15:10", syncDataTask)
}

func syncLastPrice(now time.Time) {
	version := now.Format("20060102")
	pageNum := 1
	for {
		requestUrl := fmt.Sprintf("https://proxy.finance.qq.com/cgi/cgi-bin/rank/hs/getBoardRankList?_appver=11.17.0&board_code=aStock&sort_type=price&direct=down&count=200&offset=%d", (pageNum-1)*200)
		log.Printf("[LastPrice]正在查询数据.pageNun=%d,requestUrl=%s \n", pageNum, requestUrl)
		status, body, err := httpClient.Get(context.Background(), nil, requestUrl)
		if err != nil {
			log.Printf("[LastPrice]拉数据出错.err=%s", err.Error())
			break
		}
		if status != 200 {
			log.Printf("[LastPrice]拉数据API响应非200.status=%d", status)
			break
		}
		var result apiResult
		err = sonic.Unmarshal(body, &result)
		if err != nil {
			log.Println("[LastPrice]反序列化结果出错.err=", err.Error())
			break
		}
		rankList := result.Data.RankList
		if len(rankList) == 0 {
			log.Println("[LastPrice]没有拉到数据了.pageNum=", pageNum)
			break
		}

		for _, item := range rankList {
			item.Version = version
			item.Upsert()
		}

		if len(rankList) < 200 {
			log.Println("[LastPrice]已经是最后一页了.pageNum=", pageNum)
			break
		}
		pageNum = pageNum + 1
		time.Sleep(2 * time.Second)
	}

}
