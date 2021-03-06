package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"rain/library/config"
	"rain/library/database"
	resp "rain/library/response"
	"strconv"
)

func Env(str string) (res string) {
	res = os.Getenv(str)
	return
}


// 丢弃BindJSON这种臃肿的获取值模式，采用灵活的MAP
func GetRequestJson(ctx *gin.Context) (requestMap map[string]interface{}) {
	requestData, err := ctx.GetRawData()
	if err != nil {
		resp.Error(ctx, 400,"参数获取失败")
		return
	}
	err = json.Unmarshal(requestData, &requestMap)
	if err != nil {
		resp.Error(ctx, 400,"参数获取失败")
	}

	fmt.Println(requestMap)
	return
}

func Db() (con *gorm.DB) {
	master := database.NewMySQL(config.RainDog).Master()
	return master.Write()
}

func Paginate(r *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.Query("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.Query("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
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

func Float64ToInt(f float64) (res int) {
	tmp := strconv.FormatFloat(f, 'f', -1, 64)
	var err error
	res, err = strconv.Atoi(tmp)
	if err != nil {
		fmt.Println(err)
	}
	return
}

