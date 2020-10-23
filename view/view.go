package view

import (
	"github.com/gin-gonic/gin"

	"ops-backend/view/kube/cluster"
	"ops-backend/view/asset/region"
	"ops-backend/view/project"
	"ops-backend/view/setting/auth"
	"ops-backend/pkg/view/base"
)

func RegistRouter (v base.DBView, private *gin.RouterGroup) {
	uri := v.GetUri()
	private.GET(uri, v.List)
	private.POST(uri, v.Create)
	private.PUT(uri+"/:ID", v.Update)
	private.GET(uri+"/:ID", v.Get)
	private.DELETE(uri+"/:ID", v.Delete)	
}

func RegistAllRouter (private *gin.RouterGroup) {
	// view := append(cluster.View(), region.View()...)
	view := append(region.View(), auth.View()...)
	for k, _ := range view {
		v := &view[k]
		RegistRouter(v, private)
	}
	clusterView := cluster.View()
	for _, v := range clusterView {
		RegistRouter(v, private)
	}
	projectView := project.View()
	for _, v := range projectView {
		RegistRouter(v, private)
	}
	private.GET("/project/version/:ID/last_version_deploy", project.LastVersionDeploy)
}

