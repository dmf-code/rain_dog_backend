package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"sync"
	"time"
)

type MysqlPool struct {
}

var instance *MysqlPool
var once sync.Once

var db *gorm.DB
var errorDb error

func GetInstance() *MysqlPool {
	once.Do(func() {
		instance = &MysqlPool{}
	})
	return instance
}

var Dsn string

func IsInit() (status bool) {
	if db != nil {
		return true
	}

	return false
}

func (m *MysqlPool) InitDataPool() (status bool) {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	dbName := os.Getenv("DBNAME")
	Dsn = fmt.Sprintf("%s:%s@(%s:%s)/", user, password, ip, port)
	str := fmt.Sprintf("%s%s?charset=utf8&parseTime=True&loc=Local", Dsn, dbName)
	fmt.Println(Dsn)
	fmt.Println(str)
	db, errorDb = gorm.Open("mysql", str)

	if errorDb != nil {
		log.Fatal(errorDb)
		return false
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(120)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	// 数据库中查询 show global variables like "%timeout%"
	// 查看interactive_timeout 和 wait_timeout 的值
	// 设置成和这两个值一致
	db.DB().SetConnMaxLifetime(time.Second * 120)

	// 不要默认创建数据表添加s后缀
	db.SingularTable(true)

	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

func (m *MysqlPool) SetDbName(dbName string) (_db *gorm.DB) {
	_ = os.Setenv("DBNAME", dbName)
	m.InitDataPool()

	return db
}

func (m *MysqlPool) GetMysqlDB() (_db *gorm.DB) {
	return db
}
