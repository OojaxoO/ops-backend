package auth

import (
	"ops-backend/pkg/view"
)

func View () (data []view.DBObject) {
	searchField := []string{"Name"}
	data = append(data, view.DBObject{Data: Auth{}, Uri: "/setting/auth", SearchField: searchField})
	return data
}