package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func cleanRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	err = rdb.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis database cleared successfully!")
}

func cleanMySql() {
	var conStr string = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		"root", "123456@password", "localhost", 3306, "xorderbook_dev")
	var err error
	engine, err := xorm.NewEngine("mysql", conStr)
	if err != nil {
		log.Fatal("Connect db:", err, conStr)
	}
	err = engine.Ping()
	if err != nil {
		log.Fatal(err)
	}

	engine.ShowSQL(false)
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(10)
	engine.SetConnMaxLifetime(time.Hour * 7)

	_, err = engine.Exec("truncate table limit_order")
	if err != nil {
		log.Fatal(err)
	}
	_, err = engine.Exec("truncate table match_record")
	if err != nil {
		log.Fatal(err)
	}
	_, err = engine.Exec("truncate table trade")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cleanRedis()
	cleanMySql()
}
