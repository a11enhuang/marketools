package marketdata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"com.reopenai/marketool/dal"
	"github.com/bytedance/sonic"
)

func QueryKline(ctx context.Context) {
	entries := dal.SelectAll()
	if len(entries) == 0 {
		log.Println("没有找到任何的历史行情数据")
		return
	}
	for i := range entries {
		entry := entries[i]
		requestUrl := fmt.Sprintf("https://proxy.finance.qq.com/ifzqgtimg/appstock/app/newfqkline/get?_var=kline_dayqfq&param=%s,day,,,320,qfq&r=0.5411693584834572", entry.Code)
		log.Println("正在查询历史K线数据.requetUrl=", requestUrl)
		status, body, err := httpClient.Get(ctx, nil, requestUrl)
		if err != nil {
			log.Printf("拉历史k线数据出错.err=%s", err.Error())
			break
		}
		if status != 200 {
			log.Printf("拉历史k线数据API响应非200.status=%d", status)
			break
		}

		type ApiResult struct {
			Data map[string]struct {
				Qfqday []any `json: "qfqday"`
			} `json: "data"`
		}
		result := strings.TrimPrefix(string(body), "kline_dayqfq=")

		var apiResult ApiResult
		if err := sonic.UnmarshalString(result, &apiResult); err != nil {
			log.Println("解析k线历史数据出错.err=", err.Error())
			continue
		}
		data := apiResult.Data[entry.Code]
		fmt.Print(data.Qfqday...)
	}
}
