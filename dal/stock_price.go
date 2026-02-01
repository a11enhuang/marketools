package dal

import (
	"context"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockPrice struct {
	Id        uint64          `json:"id" gorm:"primaryKey;column:id"`
	Code      string          `json:"code" gorm:"column:code"`
	Name      string          `json:"name" gorm:"column:name"`         //名字
	Hsl       decimal.Decimal `json:"hsl" gorm:"column:"`              //换手率
	Lb        decimal.Decimal `json:"lb" gorm:"column:"`               //量比
	Ltsz      decimal.Decimal `json:"ltsz" gorm:"column:"`             //流通市值
	PeTtm     decimal.Decimal `json:"pe_ttm" gorm:"column:pe_ttm"`     //市盈率
	Pn        decimal.Decimal `json:"pn" gorm:"column:pn"`             //市净率
	Speed     decimal.Decimal `json:"speed" gorm:"column:speed"`       //
	Turnover  decimal.Decimal `json:"turnover" gorm:"column:turnover"` //成交额
	Volume    decimal.Decimal `json:"volume" gorm:"column:volume"`     //成交量
	Zd        decimal.Decimal `json:"zd" gorm:"column:zd"`             //涨跌额
	Zdf       decimal.Decimal `json:"zdf" gorm:"column:zdf"`           //涨跌幅
	ZdfD10    decimal.Decimal `json:"zdf_d10" gorm:"column:zdf_d10"`   //10日涨跌幅
	ZdfD20    decimal.Decimal `json:"zdf_d20" gorm:"column:zdf_d20"`   //20日涨跌幅
	ZdfD5     decimal.Decimal `json:"zdf_d5" gorm:"column:zdf_d5"`     //5日涨跌幅
	ZdfD60    decimal.Decimal `json:"zdf_d60" gorm:"column:zdf_d60"`   //60日涨跌幅
	ZdfW52    decimal.Decimal `json:"zdf_w52" gorm:"column:zdf_w52"`   //52周涨跌幅
	ZdfY      decimal.Decimal `json:"zdf_y" gorm:"column:zdf_y"`       //年涨跌幅
	Zf        decimal.Decimal `json:"zf" gorm:"column:zf"`             //震幅
	Zljlr     decimal.Decimal `json:"zljlr" gorm:"column:zljlr"`       //主力尽流入
	Zllc      decimal.Decimal `json:"zllc" gorm:"column:zllc"`         //主力流出
	ZllcD5    decimal.Decimal `json:"zllc_d5" gorm:"column:zllc_d5"`   //主力5日流出
	Zllr      decimal.Decimal `json:"zllr" gorm:"column:zllr"`         //主力流入
	ZllrD5    decimal.Decimal `json:"zllr_d5" gorm:"column:zllr_d5"`   //主力五日流入
	Zsz       decimal.Decimal `json:"zsz" gorm:"column:zsz"`           //总市值
	Price     decimal.Decimal `json:"zxj" gorm:"column:zxj"`           //最新价
	Version   string          `gorm:"column:version"`                  // 版本号
	Timestamp time.Time       `gorm:"column:timestamp"`                // 时间
}

func (this *StockPrice) Upsert() {
	sqlDB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}, {Name: "version"}},
		DoNothing: true,
	}).Create(this)
}

func SelectBuyStocks() []StockPrice {
	entries, err := gorm.G[StockPrice](sqlDB).
		Where("lb between 1.5 and 4.5 and zsz between 50 and 300 and pe_ttm < 90 and zdf between 6 and 13 order by  zsz DESC, zdf DESC").
		Find(context.TODO())
	if err != nil {
		log.Println("查询买盘失败.error=", err.Error())
		return nil
	}
	return entries
}
