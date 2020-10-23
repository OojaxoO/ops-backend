package region

import (
	
	"ops-backend/pkg/view"
)

type Region struct {
	ID       int            `gorm:"AUTO_INCREMENT" gorm:"PRIMARY_KEY"`
	Name     string 		`json:"name" gorm:"type:varchar(32)"`
	ChName   string 		`json:"ch_name" gorm:"type:varchar(32)"`
}

func View () (data []view.DBObject) {
	searchField := []string{"Name"}
	data = append(data, view.DBObject{Data: Region{}, Uri: "/asset/region", SearchField: searchField})
	return data
}