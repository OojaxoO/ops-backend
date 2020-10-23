package user

import (
	"errors"
	"gorm.io/gorm"

	models "ops-backend/pkg/models"
	baseModel "ops-backend/pkg/models/base"
)

type User struct {
	baseModel.CRUD
	gorm.Model
	Username 	 string `json:"username" gorm:"unique"`
	Password     string `json:"password" gorm:"->:false;<-:create"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	IsSuperadmin bool   `json:"is_superadmin"`
	NickName     string `json:"nick_name"` //昵称
	Avatar       string `json:"avatar"`    //头像
	DeptId       int    `json:"detp_id"`   //部门编码
	PostId       int    `json:"post_id"`   //职位编码
	Status       string `json:"status"`
	RoleId       int    `json:"role_id"` //角色编码
}

type Auth struct {
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password" gorm:"->:false;<-:create"`
}

func (this *Auth) Check() (bool, error) {
	var user User
	err := models.DB.Where("username = ? and password = ?", this.Username, this.Password).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func (this *User) Detail() error {
	err := models.DB.Where("username = ? and password = ?", this.Username, this.Password).First(&this).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if this.ID > 0 {
		return nil
	}
	return errors.New("用户未找到")
}
//创建用户
func (this User) CreateUser() (id uint, err error) {
	//check
	var count int64
	models.DB.Table("user").Where("username = ?", this.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}
	//添加数据
	if err = models.DB.Table("user").Create(&this).Error; err != nil {
		return
	}
	id = this.ID
	return
}
//更新用户信息
func (this *User) UpdateUser(id uint) (err error) {
	if err := models.DB.Model(&this).Where("id = ?", id).Save(&this).Error; err != nil {
		err = errors.New("更新失败！")
	}
	return
}

func (user *User) DeleteUser(id []int) (Result bool, err error) {
	if err = models.DB.Table("user").Where("id in (?)", id).Delete(&User{}).Error; err != nil {
		return
	}
	Result = true
	return
}
//获取用户数据
func (this *User) SelectUser()(user User,err error) {
	table := models.DB.Table("user").Select([]string{"user.*", "role.role_name"}).Joins("left join role on user.role_id=role.role_id")
	if this.ID != 0 {
		table = table.Where("user.id = ?", this.ID)
	}

	if this.Username != "" {
		table = table.Where("user.username = ?", this.Username)
	}

	if this.Password != "" {
		table = table.Where("user.password = ?", this.Password)
	}

	if this.RoleId != 0 {
		table = table.Where("user.role_id = ?", this.RoleId)
	}

	if this.DeptId != 0 {
		table = table.Where("user.dept_id = ?", this.DeptId)
	}

	if this.PostId != 0 {
		table = table.Where("user.post_id = ?", this.PostId)
	}

	if err = table.First(&user).Error; err != nil {
		return
	}
	return
}
