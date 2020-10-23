/**
 * @Author chenjun
 * @Date 2020/9/11
 **/
package role

import "ops-backend/pkg/models"

//角色表
type Role struct {
	RoleId   int    `json:"role_id" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName string `json:"role_name" gorm:"size:128;"`                // 角色名称
	Status   string `json:"status" gorm:"size:4;"`                     //
	RoleKey  string `json:"role_key" gorm:"size:128;"`                 //角色代码
	RoleSort int    `json:"role_sort" gorm:""`                         //角色排序
	Flag     string `json:"flag" gorm:"size:128;"`                     //
	CreateBy string `json:"create_by" gorm:"size:128;"`                //
	UpdateBy string `json:"update_by" gorm:"size:128;"`                //
	Remark   string `json:"remark" gorm:"size:255;"`                   //备注
}

func (this *Role) GetList() (Role []Role, err error) {
	table := models.DB.Table("role")
	if this.RoleId != 0 {
		table = table.Where("role = ?", this.RoleId)
	}
	if this.RoleName != "" {
		table = table.Where("role_name =?", this.RoleName)
	}
	if err = table.Order("role_sort").Find(&Role).Error; err != nil {
		return
	}
	return
}
