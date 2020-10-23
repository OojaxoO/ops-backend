package host

import (
	"ops-backend/pkg/models"
)

type Host struct {
	InstanceId         string       `json:"InstanceId" gorm:"primary_key"`
	RegionId     	   string 		`json:"RegionId"`
	HostName     	   string 		`json:"HostName"`
	Memory       	   int    		`json:"Memory"`
	Cpu          	   int    		`json:"Cpu"`
	OSName       	   string       `json:"OSName" gorm:"column:os_name"`
	PublicIpAddress    string   	`json:"PublicIpAddress"`
	PrimaryIpAddress   string   	`json:"PrimaryIpAddress"`
	Status             string   	`json:"Status"`
	ExpiredTime        string       `json:"ExpiredTime"`
}

type ListResponse struct {
	Total       int       `json:"total"`
	Hosts       []Host    `json:"data"`

}

func (this *Host) CreateOrUpdate () {
	if models.DB.First(this).RowsAffected == 0 {
		models.DB.Create(&this)
		return
	}
	models.DB.Save(&this)
}
