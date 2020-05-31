package permission

import (
	"app/bootstrap/Table"
	"app/utils/helper"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"strconv"
)

const (
	PrefixUserId = "u"
	PrefixRoleId = "r"
)

var Enforcer *casbin.Enforcer

func Init() {

	enforcer, err := casbin.NewEnforcerSafe("conf/rbac_model.conf")
	if err != nil {
		fmt.Println(err)
	}

	var roles []Table.Role
	db := helper.Db()
	if err = db.Table("role").Find(&roles).Error; err != nil {
		fmt.Println(err)
	}

	if len(roles) == 0 {
		Enforcer = enforcer
	}

	for _, role := range roles {
		setRolePermission(db, enforcer, uint64(role.ID))
	}

	Enforcer = enforcer
}

// 设置角色权限
func setRolePermission(db *gorm.DB, enforcer *casbin.Enforcer, roleId uint64) {
	var roleMenus []Table.RoleMenu
	if err := db.Model(&Table.RoleMenu{RoleId: roleId}).Find(&roleMenus).Error; err != nil {
		fmt.Println(err)
	}
	for _, roleMenu := range roleMenus {
		var menu Table.Menu
		if err := db.Table("menu").Where("id = ?", roleMenu.MenuId).First(&menu).Error; err != nil {
			fmt.Println(err)
		}

		//
		if menu.Type == 3 {
			enforcer.AddPermissionForUser(
				PrefixRoleId+strconv.FormatInt(int64(roleId), 10),
				"/backend"+menu.Url,
				"GET|POST|PUT|DELETE")
		}
	}
}

// 重置角色权限
func resetRolePermission(roleId uint64) {
	if Enforcer == nil {
		return
	}
	Enforcer.DeletePermissionsForUser(PrefixRoleId + strconv.FormatInt(int64(roleId), 10))
	setRolePermission(helper.Db(), Enforcer, roleId)
}

// 设置用户角色之间的关系
func AddRoleForUser(userId uint64) (err error) {
	if Enforcer == nil {
		return
	}
	uid := PrefixUserId + strconv.FormatInt(int64(userId), 10)

	Enforcer.DeleteRolesForUser(uid)
	var adminRoles []Table.AdminRole
	db := helper.Db()
	if err = db.Table("admin_role").Model(&Table.AdminRole{AdminId: userId}).Find(&adminRoles).Error; err != nil {
		return
	}
	for _, adminRole := range adminRoles {
		Enforcer.AddRoleForUser(uid, PrefixRoleId + strconv.FormatInt(int64(adminRole.RoleId), 10))
	}
	return
}

// 删除角色
func DeleteRole(roleIds []int) {
	if Enforcer == nil {
		return
	}

	for _, roleId := range roleIds  {
		Enforcer.DeletePermissionsForUser(PrefixRoleId + strconv.FormatInt(int64(roleId), 10))
		Enforcer.DeleteRole(PrefixRoleId + strconv.FormatInt(int64(roleId), 10))
	}

}

// 检查用户是否拥有权限
func CheckPermission(userId, url, method string) (bool, error) {
	row := helper.Db().Table("admin_role").
		Where("admin_role.admin_id = ?", userId).
		Joins("left join role on admin_role.role_id = role.id").
		Select("role.name").
		Row()
	var name string
	err := row.Scan(&name)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println("name start")
	fmt.Println(name)
	fmt.Println("name end")

	if name == "super_admin" {
		return true, nil
	}

	return Enforcer.EnforceSafe(PrefixUserId + userId, url, method)
}

