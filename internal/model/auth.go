package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"rain/library/go-str"
	"rain/library/helper"
	"strings"
)

type Auth struct {
}

func (m *Auth) Login(ctx *gin.Context) (user Admin, status bool) {
	db := helper.Db()
	requestMap := helper.GetRequestJson(ctx)
	result := db.Table("admin").
		Where("username = ?", requestMap["username"]).
		First(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
		return user, false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestMap["password"].(string))); err != nil {
		fmt.Println(err)
		return user, false
	}
	fmt.Println("登录成功")
	return user, true
}

func CreateUser(username, password, rolesId string, isAdmin bool) error {
	db := helper.Db()
	return db.Transaction(func(tx *gorm.DB) error {
		admin := Admin{Username: username, Password: password}
		if err := tx.Table("admin").Create(&admin).Error; err != nil {
			return err
		}
		rolesIdList := strings.Split(rolesId, ",")
		if isAdmin {
			for _, v := range rolesIdList {
				adminRole := AdminRole{AdminId: uint64(admin.ID), RoleId: uint64(str.ToUint(v))}
				if err := tx.Table("admin_role").Create(&adminRole).Error; err != nil {
					return err
				}
			}
		} else {
			role := Role{}

			if err := tx.Table("role").FirstOrCreate(&role, Role{Name: "front_user", Memo: "前端用户", Sequence: 5}).Error; err != nil {
				return err
			}

			adminRole := AdminRole{AdminId: uint64(admin.ID), RoleId: uint64(role.ID)}
			if err := tx.Table("admin_role").Create(&adminRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *Auth) Register(ctx *gin.Context, isAdmin bool) (status bool) {
	requestMap := helper.GetRequestJson(ctx)
	var err error
	requestMap["password"], err = bcrypt.GenerateFromPassword([]byte(requestMap["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	password := string(requestMap["password"].([]byte))
	username := requestMap["username"].(string)
	rolesId, ok := requestMap["rolesId"].(string)
	if !ok {
		rolesId = ""
	}

	if err := CreateUser(username, password, rolesId, isAdmin); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
