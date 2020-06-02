package helper

import (
	"app/utils/mysqlTools"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"strconv"
)

func Success(ctx *gin.Context, data interface{}, args ...interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"status": true, "data": data, "args": args})
	ctx.Abort()
}

func Fail(ctx *gin.Context, data interface{})  {
	ctx.JSON(http.StatusOK, gin.H{"status": false, "data": data})
	ctx.Abort()
}

func Bad(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusBadRequest, gin.H{"response": data})
}

// 丢弃BindJSON这种臃肿的获取值模式，采用灵活的MAP
func GetRequestJson(ctx *gin.Context) (requestMap map[string]interface{}) {
	requestData, err := ctx.GetRawData()
	if err != nil {
		Fail(ctx, "参数获取失败")
		return
	}
	err = json.Unmarshal(requestData, &requestMap)
	if err != nil {
		Fail(ctx, "参数获取失败")
	}

	fmt.Println(requestMap)
	return
}

func Db() (con *gorm.DB) {
	return mysqlTools.GetInstance().GetMysqlDB()
}

func Env(str string) (res string) {
	res = os.Getenv(str)
	return
}

// 判断某一个值是否含在切片之中
func InArray(need interface{}, haystack interface{}) bool {
	switch key := need.(type) {
	case int:
		for _, item := range haystack.([]int) {
			if item == key {
				return true
			}
		}
	case string:
		for _, item := range haystack.([]string) {
			if item == key {
				return true
			}
		}
	case int64:
		for _, item := range haystack.([]int64) {
			if item == key {
				return true
			}
		}
	case float64:
		for _, item := range haystack.([]float64) {
			if item == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func Str2Uint(str string) (b uint) {
	a,_ := strconv.ParseUint(str, 10, 64)
	b = uint(a)
	return
}


func Float64ToInt(f float64) (res int) {
	tmp := strconv.FormatFloat(f, 'f', -1, 64)
	var err error
	res, err = strconv.Atoi(tmp)
	if err != nil {
		fmt.Println(err)
	}
	return
}