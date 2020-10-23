/**
 * @Author chenjun
 * @Date 2020/9/11
 **/
package post

import "ops-backend/pkg/models"

//岗位表
type Post struct {
	PostId   int    `gorm:"primary_key;AUTO_INCREMENT" json:"post_id"` //岗位编号
	PostName string `gorm:"size:128;" json:"post_name"`                //岗位名称
	PostCode string `gorm:"size:128;" json:"post_code"`                //岗位代码
	Sort     int    `gorm:"" json:"sort"`                              //岗位排序
	Status   string `gorm:"size:4;" json:"status"`                     //状态
	Remark   string `gorm:"size:255;" json:"remark"`                   //描述
	CreateBy string `gorm:"size:128;" json:"create_by"`
	UpdateBy string `gorm:"size:128;" json:"update_by"`
}

func (Post) TableName() string {
	return "sys_post"
}

func (this *Post) GetList() ([]Post, error) {
	var doc []Post

	table := models.DB.Select("*").Table("post")
	if this.PostId != 0 {
		table = table.Where("post_id = ?", this.PostId)
	}
	if this.PostName != "" {
		table = table.Where("post_name = ?", this.PostName)
	}
	if this.PostCode != "" {
		table = table.Where("post_code = ?", this.PostCode)
	}
	if this.Status != "" {
		table = table.Where("status = ?", this.Status)
	}
	if err := table.Find(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}
func (this *Post) Get() (Post, error) {
	var doc Post
	table := models.DB.Table(this.TableName())
	if this.PostId != 0 {
		table = table.Where("pod_id = ?", this.PostId)
	}
	if this.PostName != "" {
		table = table.Where("post_name = ?", this.PostName)
	}
	if this.PostCode != "" {
		table = table.Where("post_code = ?", this.PostCode)
	}
	if this.Status != "" {
		table = table.Where("status = ?", this.Status)
	}
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil

}

func (this *Post) Create() (Post, error) {
	var doc Post
	result := models.DB.Table(this.TableName()).Create(&this)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *this
	return doc, nil
}

func (this *Post) Update(id int) (update Post, err error) {
	if err = models.DB.Table(this.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if err = models.DB.Table(this.TableName()).Model(&update).Save(&this).Error; err != nil {
		return
	}
	return
}
func (this *Post) BatchDelete(ids []int) (Result bool, err error) {
	if err = models.DB.Table(this.TableName()).Where("post_id in (?)", ids).Delete(&Post{}).Error; err != nil {
		return
	}
	Result = true
	return
}
