package tencent

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"com.reopenai/marketool/dal"
	"com.reopenai/marketool/dingding"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/robfig/cron/v3"
)

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

func Run() {
	syncData(time.Now())
	buyStocks()
	c := cron.New(cron.WithSeconds())

	task := func() {
		syncData(time.Now())
	}

	c.AddFunc("15,30,45 9 * * 1-5", task)
	c.AddFunc("0,15,30,45 10 * * 1-5", task)
	c.AddFunc("0,15,30 11 * * 1-5", task)

	c.AddFunc("0,15,30,45 13 * * 1-5", task)
	c.AddFunc("0,15,30,45 14 * * 1-5", task)
	c.AddFunc("0 15 15 * * 1-5", task)

	c.AddFunc("30 9 * * 1-5", buyStocks)
	c.AddFunc("30 11 * * 1-5", buyStocks)
	c.AddFunc("30 14 * * 1-5", buyStocks)

	c.Start()
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
	}
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
		time.Sleep(5 * time.Second)
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
