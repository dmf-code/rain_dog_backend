package main

import (
	"app/bootstrap"
	"app/routes"
	"app/utils/mysqlTools"
	"app/utils/permission"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func init()  {
	// 配置日志
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 加载.env配置
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化Mysql连接池
	if !mysqlTools.GetInstance().InitDataPool() {
		log.Println("init database mysqlTools failure...")
		os.Exit(1)
	}

	// 权限初始化
	permission.Init()

	// 迁移数据
	bootstrap.InitTable()

}

func main() {
	secretKey := os.Getenv("SECRET_KEY")
	fmt.Println(secretKey)

	r := routes.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}
