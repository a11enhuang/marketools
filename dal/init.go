package dal

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sqlDB *gorm.DB

func init() {
	viper.SetDefault("POSTGRES_USERNAME", "postgres")
	viper.SetDefault("POSTGRES_DATABASE", "stock")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	host := viper.GetString("POSTGRES_HOST")
	username := viper.GetString("POSTGRES_USERNAME")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbname := viper.GetString("POSTGRES_DATABASE")
	port := viper.GetString("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, username, password, dbname, port)
	log.Println("正在连接数据库: ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB = db

	instance, err := db.DB()
	if err != nil {
		panic(err)
	}
	instance.SetMaxIdleConns(10)
	instance.SetMaxOpenConns(100)
	instance.SetConnMaxLifetime(time.Hour)

	initTable()
}

func initTable() {
	tableStockPrice := `
	CREATE TABLE IF NOT EXISTS stock_prices (
		id BIGSERIAL PRIMARY KEY,
		code VARCHAR(64) NOT NULL,
		name VARCHAR(128) NOT NULL,
		hsl DECIMAL(18,6) NOT NULL DEFAULT 0,
		lb DECIMAL(18,6) NOT NULL DEFAULT 0,
		ltsz DECIMAL(18,6) NOT NULL DEFAULT 0,
		pe_ttm DECIMAL(18,6) NOT NULL DEFAULT 0,
		pn DECIMAL(18,6) NOT NULL DEFAULT 0,
		speed DECIMAL(18,6) NOT NULL DEFAULT 0,
		turnover DECIMAL(18,6) NOT NULL DEFAULT 0,
		volume DECIMAL(18,6) NOT NULL DEFAULT 0,
		zd DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_d10 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_d20 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_d5 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_d60 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_w52 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zdf_y DECIMAL(18,6) NOT NULL DEFAULT 0,
		zf DECIMAL(18,6) NOT NULL DEFAULT 0,
		zljlr DECIMAL(18,6) NOT NULL DEFAULT 0,
		zllc DECIMAL(18,6) NOT NULL DEFAULT 0,
		zllc_d5 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zllr DECIMAL(18,6) NOT NULL DEFAULT 0,
		zllr_d5 DECIMAL(18,6) NOT NULL DEFAULT 0,
		zsz DECIMAL(18,6) NOT NULL DEFAULT 0,
		zxj DECIMAL(18,6) NOT NULL DEFAULT 0,
		version VARCHAR(64) NOT NULL
	);
	CREATE UNIQUE INDEX IF NOT EXISTS udx_stock_prices ON stock_prices(code, version DESC);
	`
	result := gorm.WithResult()
	err := gorm.G[any](sqlDB, result).Exec(context.Background(), tableStockPrice)
	if err != nil {
		panic(err)
	}

}
