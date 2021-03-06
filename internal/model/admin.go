package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"rain/library/format"
	"rain/library/go-str"
	"rain/library/helper"
	"rain/library/response"
	"strings"
	"time"
)

type Row struct {
	Username string
	RoleIds  string
}

type Admin struct {
	BaseModel
	Username string `gorm:"column:username;size:32;not null;unique_index; comment:'用户名';" json:"username" form: "username"`
	Password string `gorm:"column:password;size:128;not null;" json:"password; comment:'密码';" form: "password"`
}

func (Admin) TableName() string {
	return TableName("admin")
}

// 添加之前
func (m *Admin) BeforeCreate(scope *gorm.Scope) error {
	m.CreatedAt = format.JSONTime{Time: time.Now()}
	m.UpdatedAt = format.JSONTime{Time: time.Now()}
	return nil
}

// 更新之前
func (m *Admin) BeforeUpdate(scope *gorm.Scope) error {
	m.UpdatedAt = format.JSONTime{Time: time.Now()}
	return nil
}

type AdminSon struct {
	Admin
	RoleIds  uint64 `gorm:"column:role_id;unique_index:uk_admin_role_admin_id;not null;comment='角色ID'"`
}

func (m *Admin) Index(ctx *gin.Context) {
	db := helper.Db()
	var fields []AdminSon
	if err := db.Table("admin").Select("id,username,created_at,updated_at").Find(&fields).Error; err != nil {
		resp.Error(ctx, 400, "查询失败")
		return
	}
	var rows []Row
	db.Table("admin").
		Select("admin.username, group_concat(admin_role.role_id) as role_ids").
		Joins("left join admin_role on admin_role.admin_id = admin.id").
		Group("admin.username").
		Find(&rows)
	fmt.Println("start")
	fmt.Println(rows)

	for k, v := range fields {
		for _, vv := range rows {
			if v.Username == vv.Username {
				fields[k].RoleIds = uint64(str.ToUint(vv.RoleIds))
			}
		}
	}

	resp.Success(ctx, "ok", fields)
}

func (m *Admin) Show(ctx *gin.Context) {
	db := helper.Db()
	var field Admin
	if err := db.Table("admin").Select("id,username").Where("id = ?", ctx.Param("id")).First(&field).Error; err != nil {
		resp.Error(ctx, 400, "查询失败")
		return
	}

	resp.Success(ctx, "ok", field)
}

func (m *Admin) Store(ctx *gin.Context) {
	auth := Auth{}
	status := auth.Register(ctx, true)

	if !status {
		resp.Error(ctx, 400, "注册失败")
	}

	resp.Success(ctx, "ok")
}

func (m *Admin)  updateAdmin(filed Admin, requestJson map[string]interface{}) error {
	db := helper.Db()
	return db.Transaction(func(tx *gorm.DB) error {
		delete(requestJson, "password")
		delete(requestJson, "createdAt")
		delete(requestJson, "updatedAt")
		if err := tx.Table("admin").Model(&filed).Updates(requestJson).Error; err != nil {
			return err
		}
		if err := tx.Table("admin_role").Where("admin_id = ?", filed.ID).Delete(AdminRole{}).Error; err != nil {
			return err
		}

		roleIds := strings.Split(requestJson["role_ids"].(string), ",")
		fmt.Println(roleIds)
		for _, v := range roleIds {
			if err := tx.Table("admin_role").Create(&AdminRole{AdminId: uint64(filed.ID), RoleId: uint64(str.ToUint(v))}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *Admin) ResetPassword(ctx *gin.Context) {
	requestJson := helper.GetRequestJson(ctx)
	id := str.ToUint(ctx.Param("id"))
	password := requestJson["password"]
	password2 := requestJson["password2"]

	if password != password2 || password == "" {
		resp.Error(ctx, 400, "重置失败，密码不能为空或两次输入不一致")
		return
	}

	db := helper.Db()
	newPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
	if err != nil {
		resp.Error(ctx, 400, err.Error())
		return
	}
	if err := db.Table("admin").Where("id = ?", id).Update("password", newPassword).Error; err != nil {
		resp.Error(ctx, 400, err.Error())
	}

	resp.Success(ctx, "重置成功")
}

func (m *Admin) Update(ctx *gin.Context) {
	var filed Admin
	requestJson := helper.GetRequestJson(ctx)
	filed.ID = str.ToUint(ctx.Param("id"))

	if err := m.updateAdmin(filed, requestJson); err != nil {
		fmt.Println(err)
		resp.Error(ctx, 400, "更新失败")
		return
	}

	resp.Success(ctx,  "更新成功")
}

func (m *Admin) Destroy(ctx *gin.Context) {
	db := helper.Db()
	var field Admin
	field.ID = str.ToUint(ctx.Param("id"))
	if err := db.Table("admin").Delete(&field).Error; err != nil {
		resp.Error(ctx, 400, err.Error())
		return
	}

	resp.Success(ctx,  "删除成功")
}
