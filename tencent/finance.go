package tencent

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"com.reopenai/marketool/dal"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
)

func Run() {
	syncData(time.Now())

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		now := time.Now()
		if isWeekday(now) {
			triggerSyncData(now)
		}
		<-ticker.C
	}
}

func isWeekday(t time.Time) bool {
	return t.Weekday() < time.Monday || t.Weekday() > time.Friday
}

func triggerSyncData(t time.Time) {
	h, m := t.Hour(), t.Minute()
	if h < 9 || h > 15 {
		return
	}
	if h == 9 && m < 15 {
		return
	}
	if h == 15 && m > 0 {
		return
	}
	if (h > 11 && m > 30) && h < 13 {
		return
	}
	syncData(t)
}

var httpClient *client.Client

func init() {
	clientCfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, _ := client.NewClient(
		client.WithTLSConfig(clientCfg),
		client.WithKeepAlive(true),
		client.WithDialer(standard.NewDialer()),
	)
	httpClient = c
}

func syncData(now time.Time) {
	version := now.Format("20060102")
	pageNum := 1
	for {
		requestUrl := fmt.Sprintf("https://proxy.finance.qq.com/cgi/cgi-bin/rank/hs/getBoardRankList?_appver=11.17.0&board_code=aStock&sort_type=price&direct=down&count=200&offset=%d", (pageNum-1)*200)
		log.Printf("正在查询数据.pageNun=%d,requestUrl=%s \n", pageNum, requestUrl)
		status, body, err := httpClient.Get(context.Background(), nil, requestUrl)
		if err != nil {
			log.Printf("拉数据出错.err=%s", err.Error())
			break
		}
		if status != 200 {
			log.Printf("拉数据API响应非200.status=%d", status)
			break
		}
		var result apiResult
		err = sonic.Unmarshal(body, &result)
		if err != nil {
			log.Println("反序列化结果出错.err=", err.Error())
			break
		}
		rankList := result.Data.RankList
		if len(rankList) == 0 {
			log.Println("没有拉到数据了.pageNum=", pageNum)
			break
		}

		for _, item := range rankList {
			item.Version = version
			item.Upsert()
		}

		if len(rankList) < 200 {
			log.Println("已经是最后一页了.pageNum=", pageNum)
			break
		}
		pageNum = pageNum + 1
		time.Sleep(5000)
	}

}

type apiResult struct {
	Code   int      `json:"code"`
	Msg    string   `json:"msg"`
	Data   RankData `json:"data"`
	Offset int      `json:"offset"`
}

type RankData struct {
	RankList []*dal.StockPrice `json:"rank_list"`
}
