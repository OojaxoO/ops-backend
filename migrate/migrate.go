package main 

import (
	"log"
	
	"ops-backend/view/kube/cluster"
	"ops-backend/view/asset/region"
	"ops-backend/view/setting/auth"
	"ops-backend/view/user"
	"ops-backend/view/project"
	"ops-backend/pkg/models"
	"ops-backend/pkg/setting"
)

func main () {
	log.Println("migrate starting...")

	setting.Setup()
	models.Setup()

	models.DB.AutoMigrate(&region.Region{})
	models.DB.AutoMigrate(&user.User{})
	models.DB.SetupJoinTable(&project.Project{}, "Envs", &project.ProjectEnv{})
	models.DB.SetupJoinTable(&project.App{}, "Config", &project.Config{})
	models.DB.AutoMigrate(&cluster.Cluster{}, &auth.Auth{}, &project.Env{}, &project.Project{}, &project.ProjectEnv{},
						  &project.App{}, &project.Version{}, &project.AppDeploy{}, &project.AppDeployAction{},
						  &project.Config{}, &project.VersionDeploy{}, &project.VersionDeployAction{})
}
