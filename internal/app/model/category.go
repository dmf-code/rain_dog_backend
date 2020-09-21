package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"rain/library/helper"
)

type Category struct {
	BaseModel
	Name string `json:"name" gorm:"column:name;comment: '分类名';"`
	Num  int    `json:"num"  gorm:"column:num;comment:'分类使用次数';"`
}

func (m *Category) Index(ctx *gin.Context) {
	db := helper.Db()
	var fields []Category
	if err := db.Table("category").Find(&fields).Error; err != nil {
		helper.Fail(ctx, "查询失败")
		return
	}

	helper.Success(ctx, fields)
}

func (m *Category) Show(ctx *gin.Context) {
	db := helper.Db()
	var field Category
	if err := db.Table("category").Where("id = ?", ctx.Param("id")).First(&field).Error; err != nil {
		helper.Fail(ctx, "查询失败")
		return
	}

	helper.Success(ctx, field)
}

func (m *Category) Store(ctx *gin.Context) {
	db := helper.Db()
	var field Category
	err := ctx.Bind(&field)
	field.Num = 0
	fmt.Println(field)
	if err != nil {
		helper.Fail(ctx, "绑定数据失败")
		return
	}
	if err = db.Table("category").Create(&field).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "success")
}

func (m *Category) Update(ctx *gin.Context) {
	db := helper.Db()
	var filed Category
	requestJson := helper.GetRequestJson(ctx)
	filed.ID = helper.Str2Uint(ctx.Param("id"))
	if err := db.Table("category").Model(&filed).Updates(requestJson).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "更新成功")
}

func (m *Category) Destroy(ctx *gin.Context) {
	db := helper.Db()
	var field Category
	field.ID = helper.Str2Uint(ctx.Param("id"))
	if err := db.Table("category").Delete(&field).Error; err != nil {
		helper.Fail(ctx, err.Error())
		return
	}

	helper.Success(ctx, "删除成功")
}