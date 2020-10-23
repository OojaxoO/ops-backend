package base

import (
	"ops-backend/pkg/models"
)

type CRUDInterface interface {
	List(data interface{}, page int, pageSize int, filter string, filterField interface{}, preload []string) (ListResponse, error)
	Create (data interface{}) error 	
	// Get (data interface{}) error	
	Update (data interface{}) error	
	Delete (data interface{}) error
	Find (data interface{}, filterField interface{}, filter string, preload []string) error	
	FindAll (data interface{}) error 	
	Get (data interface{}, preload []string) error
}

type CRUD struct {
}

type ListResponse struct {
	Total       int64           `json:"total"`
	Data        interface{}     `json:"data"`
}

func (this *CRUD) List (data interface{}, page int, pageSize int, filter string, filterField interface{}, preload []string) (ListResponse, error) {
	responce := ListResponse{Total: 0, Data: data}
	start := (page - 1) * pageSize
	db := models.DB
	for _, reload := range preload {
		db = db.Preload(reload)
	}
	err := db.Where(filterField).Where(filter).Limit(pageSize).Offset(start).Find(data).Error
	if err != nil {
		return responce, err 
	}
	models.DB.Model(data).Count(&responce.Total)
	return responce, nil 
}

func (this *CRUD) Get (data interface{}, preload []string) error {
	db := models.DB
	for _, reload := range preload {
		db = db.Preload(reload)
	}
	return db.Where(data).Find(data).Error
}

func (this *CRUD) Create (data interface{}) error {
	return models.DB.Create(data).Error
}

// func (this *CRUD) Get (data interface{}) error {
// 	return models.DB.Where(data).First(data).Error
// }

func (this *CRUD) Update (data interface{}) error {
	return models.DB.Save(data).Error	
}

func (this *CRUD) Delete (data interface{}) error {
	return models.DB.Delete(data).Error	
}

func (this *CRUD) Find (data interface{}, filterField interface{}, filter string, preload []string) error {
	db := models.DB
	for _, reload := range preload {
		db = db.Preload(reload)
	}
	return db.Where(filterField).Where(filter).Find(data).Error	
}

func (this *CRUD) FindAll (data interface{}) error {
	return models.DB.Find(data).Error	
}