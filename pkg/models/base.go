package models

import (
)

type ListResponse struct {
	Total       int64           `json:"total"`
	Data        interface{}     `json:"data"`

}

func List (page int, pageSize int, data interface{}, fieldFilter interface{}, filter interface{}) (ListResponse, error) {
	responce := ListResponse{Total: 0, Data: data}
	start := (page - 1) * pageSize
	err := DB.Where(fieldFilter).Where(filter).Limit(pageSize).Offset(start).Find(data).Error
	if err != nil {
		return responce, err 
	}
	DB.Model(data).Count(&responce.Total)
	return responce, nil 
}

func Create (data interface{}) error {
	return DB.Create(data).Error
}

func Get (data interface{}) error {
	return DB.Where(data).First(data).Error
}

func Update (data interface{}) error {
	return DB.Save(data).Error	
}

func Delete (data interface{}) error {
	return DB.Delete(data).Error	
}

func Find (data interface{}, fieldFilter interface{}, filter interface{}) error {
	return DB.Where(fieldFilter).Where(filter).Find(data).Error	
}

func FindAll (data interface{}) error {
	return DB.Find(data).Error	
}