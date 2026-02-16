package dal

import (
	"time"

	"github.com/shopspring/decimal"
)

type Kline struct {
	open         decimal.Decimal
	close        decimal.Decimal
	high         decimal.Decimal
	low          decimal.Decimal
	zdf          decimal.Decimal
	zde          decimal.Decimal
	trade_amount decimal.Decimal
	timestamp    time.Time
}
