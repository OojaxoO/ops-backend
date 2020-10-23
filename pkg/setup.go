package pkg

import (
	"ops-backend/pkg/setting"
	"ops-backend/pkg/models"
	"ops-backend/pkg/util"
)

func Setup() {
	setting.Setup()
	models.Setup()
	util.Setup()
}