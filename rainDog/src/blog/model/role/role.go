package role

import (
	"blog/model/manage"
	"blog/utils/helper"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	db := helper.Db("rain_dog")
	var fields []manage.Role
	if err := db.Table("role").Find(&fields).Error; err != nil {
		helper.Fail(ctx, "查询失败")
		return
	}

	helper.Success(ctx, fields)
}

func Show(ctx *gin.Context) {
	db := helper.Db("rain_dog")
	var field manage.Role
	if err := db.Table("role").Where("id = ?", ctx.Param("id")).First(&field).Error; err != nil {
		helper.Fail(ctx, "查询失败")
		return
	}

	helper.Success(ctx, field)
}

func Store(ctx *gin.Context) {
	db := helper.Db("rain_dog")
	var field manage.Role
	err := ctx.Bind(&field)
	fmt.Println(field)
	if err != nil {
		helper.Fail(ctx, "绑定数据失败")
		return
	}
	if err = db.Table("role").Create(&field).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "success")
}

func Update(ctx *gin.Context) {
	db := helper.Db("rain_dog")
	var filed manage.Role
	requestJson := helper.GetRequestJson(ctx)
	filed.ID = helper.Str2Uint(ctx.Param("id"))
	if err := db.Table("role").Model(&filed).Updates(requestJson).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "更新成功")
}

func Destroy(ctx *gin.Context) {
	db := helper.Db("rain_dog")
	var field manage.Role
	field.ID = helper.Str2Uint(ctx.Param("id"))
	if err := db.Table("role").Delete(&field).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "删除成功")
}
