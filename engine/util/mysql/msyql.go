package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func Init() {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWD")
	port := os.Getenv("MYSQL_PORT")
	db := os.Getenv("MYSQL_DB")
	var conStr string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&charset=utf8mb4&collation=utf8mb4_unicode_ci", user, password, host, port, db)
	var err error
	engine, err = xorm.NewEngine("mysql", conStr)
	err = engine.Ping()
	if err != nil {
		log.Fatal("Connect db:", err, conStr)
	}

	engine.ShowSQL(false)
	//tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "tbl_")
	//engine.SetTableMapper(tbMapper)
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(10)
	engine.SetConnMaxLifetime(time.Hour * 7)
}

func InitMysql(host string, user string, password string, port int64, db string) {
	var conStr string = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?&charset=utf8mb4&collation=utf8mb4_unicode_ci", user, password, host, port, db)
	var err error
	engine, err = xorm.NewEngine("mysql", conStr)
	err = engine.Ping()
	if err != nil {
		log.Fatal("Connect db:", err, conStr)
	}

	engine.ShowSQL(false)
	//tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "tbl_")
	//engine.SetTableMapper(tbMapper)
	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(10)
	engine.SetConnMaxLifetime(time.Hour * 7)
}

func GetSession() *xorm.Session {
	return engine.NewSession()
}

func GetEngine() *xorm.Engine {
	return engine
}
