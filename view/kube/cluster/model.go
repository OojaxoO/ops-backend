package cluster

import (
	"gorm.io/gorm"

	baseView "ops-backend/pkg/view/base"
	baseModel "ops-backend/pkg/models/base"
)

type Cluster struct {
	baseModel.CRUD
	gorm.Model
	Name     string 		`gorm:"type:varchar(64)"`
	Config   string 		`gorm:"type:text"`
}

func View () (data []baseView.DBView) {
	searchField := []string{"Name"} 
	data = append(data, &baseView.BaseView{Data: &Cluster{}, Uri: "/kube/cluster", SearchField: searchField})
	return data
}


