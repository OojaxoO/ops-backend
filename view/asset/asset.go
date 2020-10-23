package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	"ops-backend/pkg/models"
	"ops-backend/view/asset/host"
	"ops-backend/view/kube/cluster"
)

type Count struct {
	HostCount    int64 `json:"host_count"`
	ClusterCount int64 `json:"cluster_count"`
	ProjectCount int64 `json:"project_count"`
	ServiceCount int64 `json:"service_count"`
}

func Info(c *gin.Context) {
	appG := app.Gin{C: c}
	var count Count

	if err := models.DB.Model(&host.Host{}).Count(&count.HostCount).Error; err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
		return
	}
	if err := models.DB.Model(&cluster.Cluster{}).Count(&count.ClusterCount).Error; err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, count)
}
