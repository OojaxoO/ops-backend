package project

import (
	"net/http"
	// "reflect"

	"github.com/gin-gonic/gin"

	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	"ops-backend/pkg/view"
	baseView "ops-backend/pkg/view/base"
)

type ProjectView struct {
	baseView.BaseView
}

type VersionView struct {
	baseView.BaseView
}

// 重写
func (this *VersionView) Get (c *gin.Context) {
	data := &Version{}
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	detail, err := data.Detail()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, detail) 		
}

func LastVersionDeploy (c *gin.Context) {
	data := &Version{}
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	versionDeploy, err := data.LastVersionDeploy()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if versionDeploy != nil {
		appG.Response(http.StatusOK, e.SUCCESS, *versionDeploy) 		
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil) 		
}

func View () (data []baseView.DBView) {
	searchField := []string{"Name"}
	projectView := ProjectView{baseView.BaseView{
		Data: &Project{}, 
		Uri: "/project/project", 
		SearchField: searchField, 
		Preload: []string{"Jenkins", "Git", "Envs"}}} 
    versionView := VersionView{
		baseView.BaseView{
			Data: &Version{}, 
			Uri: "/project/version", 
			SearchField: searchField, 
			Preload: []string{"Project", "User"}}}
	data = append(data, &projectView) 
	data = append(data, &versionView) 
	data = append(data, &baseView.BaseView{
		Data: &ProjectEnv{}, 
		Uri: "/project/project_env",
		Preload: []string{"Env", "Cluster", "Admin"}}) 
	data = append(data, &baseView.BaseView{
		Data: &App{}, 
		Uri: "/project/app",
		SearchField: searchField,
		Preload: []string{"Project"}}) 
	data = append(data, &baseView.BaseView{
		Data: &VersionDeploy{}, 
		Uri: "/project/version_deploy",
		SearchField: searchField,
		Preload: []string{"Project"}}) 
	data = append(data, &view.DBObject{Data: Env{}, Uri: "/project/env", SearchField: searchField})
	data = append(data, &baseView.BaseView{Data: &Config{}, Uri: "/project/config", SearchField: searchField, Preload: []string{"ProjectEnv.Cluster", "ProjectEnv.Env"}})
	return data
}