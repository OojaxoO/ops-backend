package project

import (
	"time"
	"errors"
	//"reflect"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"

	"ops-backend/view/setting/auth"
	"ops-backend/view/user"
	"ops-backend/view/kube/cluster"
	"ops-backend/pkg/models/base"
	"ops-backend/pkg/models"
)

type Env struct {
	base.CRUD
	gorm.Model
	Name      string     `json:"name"     gorm:"type:varchar(64);comment:'名字'"`
	ChName    string     `json:"ch_name"  gorm:"type:varchar(64);comment:'中文名字'"`
	Step      int        `json:"step"     gorm:"type:int(2);commnet:'阶段'"`
	Projects  []Project  `json:"projects" gorm:"many2many:project_env"`
}

type EnvShort struct {
	ID        int	
	Name      string     `json:"name"     gorm:"type:varchar(64);comment:'名字'"`
	ChName    string     `json:"ch_name"  gorm:"type:varchar(64);comment:'中文名字'"`
	Step      int        `json:"step"     gorm:"type:int(2);commnet:'阶段'"`
}

func (EnvShort) TableName () string {
	return "env"
}

type Project struct {
	base.CRUD
	gorm.Model
	Name      string     	  `json:"name"       gorm:"type:varchar(64);comment:'名字'"` 
	ChName    string     	  `json:"ch_name"    gorm:"type:varchar(64);comment:'中文名字'"`
	GitID     int        	  `json:"git_id"     gorm:"comment:'git认证ID'"` 
	Git 	  auth.Auth  	  `json:"git"        gorm:"foreignKey:GitID"`             					           
	JenkinsID int        	  `json:"jenkins_id" gorm:"comment:'Jenkins认证ID'"`
	Jenkins   auth.Auth  	  `json:"jenkins"    gorm:"foreignKey:JenkinsID"` 
	Envs      []Env           `json:"envs"       gorm:"many2many:project_env"`
	ProjectEnv []ProjectEnv   `json:"project_env"` 
	Apps      []App           `json:"apps"`
}

// 若重写
func (this *Project) List (data interface{}, page int, pageSize int, filter string, filterField interface{}, preload []string) (base.ListResponse, error) {
	responce := base.ListResponse{Total: 0, Data: data}
	start := (page - 1) * pageSize
	if err := models.DB.Joins("Jenkins").
						Joins("Git").
						Preload("Apps"). 
						Where(filterField).
						Where(filter).
						Limit(pageSize).
						Offset(start).
						Find(data).
						Error; 
						err != nil {
		return responce, err 
	}
	models.DB.Model(this).Count(&responce.Total)
	return responce, nil 
}

type ProjectEnv struct {
	// 若需要configmap后续可加该字段
	base.CRUD
	gorm.Model
	ProjectID     int    			`json:"project_id" gorm:"comment:'项目id'"`        
	EnvID         int    			`json:"env_id"     gorm:"comment:'环境id'"` 
	Env           Env               `json:"env"        gorm:"<-:false"`
	ClusterID     int    			`json:"cluster_id" gorm:"comment:'集群id'"`
	Cluster       cluster.Cluster   `json:"cluster"    gorm:"<-:false"` 
	Namespace     string            `json:"namespace"  gorm:"type:varchar(32);comment:'命名空间'"` 
	Admin         []user.User       `json:"admin"      gorm:"many2many:project_env_admin"`
}

// func (this *ProjectEnv) Create (data interface{}) error {
// 	if err := models.DB.Create(data).Error; err != nil {
// 		return err
// 	}
// 	models.DB.Model(this).Association("Admin").Replace(this.Admin)
// 	return nil
// }
// 
// func (this *ProjectEnv) Update (data interface{}) error {
// 	// models.DB.Model(data).Association("Admin").Clear()
// 	if err := models.DB.Save(data).Error; err != nil {
// 		return err
// 	}	
// 	return nil
// }

// func (this *ProjectEnv) BeforeSave (tx *gorm.DB) (err error) {
//	return models.DB.Delete(&this.Admin).Error
// }

type ProjectEnvShort struct {
	ID              int	
	Env             EnvShort    `json:"env"` 
}

func (ProjectEnvShort) TableName () string {
	return "project_env"
} 

type ProjectIDName struct {
	ID   	  int 
	Name      string     	  `json:"name"       gorm:"type:varchar(64);comment:'名字'"` 
	ChName    string     	  `json:"ch_name"    gorm:"type:varchar(64);comment:'中文名字'"`
}

func (ProjectIDName) TableName () string {
	return "project"
}

type App struct {
	base.CRUD
	gorm.Model
	Name                string         		 `json:"name"       gorm:"type:varchar(64);comment:'名字'"` 
	ChName              string         		 `json:"ch_name"  gorm:"type:varchar(64);comment:'中文名字'"`
	ProjectID           int            		 `json:"project_id"  gorm:"comment:'项目id'"`
	Project             ProjectIDName        `json:"project"`
	GitUri              string         		 `json:"git_uri"     gorm:"type:varchar(64);comment:'git子路径'"`
	JenkinsJob          string         		 `json:"jenkins_job" gorm:"type:varchar(64);comment:'jenkins job名'"`
	ProjectEnv          []ProjectEnv   		 `json:"project_env"      gorm:"many2many:config;comment:'配置'"` 
	Config              []Config             `json:"config"`
}

func (this *App) AfterCreate(db *gorm.DB) (err error) {
	// 创建应用后, 添加配置
	var projectEnvs []ProjectEnv
	var configs  []Config
	db.Select("id").Where("project_id = ?", this.ProjectID).Find(&projectEnvs)
	for _, projectEnv := range projectEnvs {
		configs = append(configs, Config{AppID: this.ID, ProjectEnvID: projectEnv.ID})
	}
	if err := db.Create(&configs).Error; err != nil {
		return err
	}
	return
}

type AppDeploy struct {
	base.CRUD
	gorm.Model
	Version        Version        `json:"version"`
	VersionID      int            `json:"version_id"  gorm:"comment:'版本id'"` 
	Branch         string         `json:"branch"       gorm:"type:varchar(64);comment:'打包分支'"`
	Release        int            `json:"release"      gorm:"type:int(3);comment:'发型号'"`
	App            App            `json:"app"`
	AppID          int            `json:"app_id"       gorm:"comment:'应用id'"`
	User           user.User      `json:"user"`
	UserID         int            `json:"user_id"      gorm:"comment:'用户id'"`
}

type AppDeployAction struct {
	base.CRUD
	gorm.Model
	Type           string         `json:"type"           gorm:"type:varchar(32);comment:'update,build,rollback,reload,stop'"`
	AppDeploy      AppDeploy      `json:"app_deploy"`
	AppDeployID    int            `json:"app_deploy_id"  gorm:"comment:'应用部署id'"`
	User           user.User      `json:"user"`
	UserID         int            `json:"user_id"        gorm:"comment:'执行者'"`
	Result         string         `json:"result"         gorm:"type:text;comment:'结果'"`
	ProjectEnv     ProjectEnv     `json:"project_env"`
	ProjectEnvID   int      	  `json:"project_evn_id" gorm:"comment:'项目环境id'"`
}

type AppConfig struct {
	ID              int	
	ProjectEnv      ProjectEnv      `json:"project_env"`
	Chart           string          `json:"chart"          gorm:"type:varchar(32);comment:'chart名'"`
	Values          string          `json:"values"         gorm:"type:text;comment:'values.yaml'"` 
}

func (AppConfig) TableName () string {
	return "config"
} 

type Config struct {
	base.CRUD
	gorm.Model
	AppID         uint      	  `json:"app_id"         gorm:"comment:'应用id'"` 
	ProjectEnvID  uint      	  `json:"project_evn_id" gorm:"comment:'项目环境id'"`
	Chart         string          `json:"chart"          gorm:"type:varchar(32);comment:'chart名'"`
	ProjectEnv    ProjectEnv      `json:"project_env"    gorm:"<-:false"`
	Values        string   		  `json:"values"         gorm:"type:text;comment:'values.yaml'"` 
}

type Version struct {
	base.CRUD
	gorm.Model	
	Name         		string      	`json:"name"        gorm:"type:varchar(32);comment:'名字'"`  
	ProjectID    		int         	`json:"project_id"  gorm:"comment:'项目id'"`
	Project      		Project     	`json:"project"     gorm:"<-:false"` 
	Desc         		string   	 	`json:"desc"        gorm:"type:text;comment:'描述'"`   
	DeployTime   		time.Time  	 	`json:"deploy_time" gorm:"comment:'上线预估时间'"` 
	User         		user.User   	`json:"user"        gorm:"<-:false"` 
	UserID       		int         	`json:"user_id"     gorm:"comment:'创建者'"` 
	// LastVersionDeploy   *VersionDeploy  `json:"last_version_deploy" gorm:"<-:false"` 
}

type VersionDetail struct {
	ID      		uint
	Name            string      	     `json:"name"`	
	ProjectID       int
	Project 		Project  			 `json:"project"` 
}

func (VersionDetail) TableName () string {
	return "version"
} 

func (this *Version) Detail() (vd VersionDetail, err error) {
	db := models.DB
	err = db.Preload("Project.ProjectEnv.Cluster").
			 Preload("Project.ProjectEnv.Env").
			 Preload("Project.Apps.Config").
			 Where(this).
			 Find(&vd).
			 Error
	return
}

func (this *Version) LastVersionDeploy() (*VersionDeploy, error) {
	var vd *VersionDeploy = new(VersionDeploy)
	db := models.DB
	err := db.Preload("AppDeploy").
			 Preload("VersionDeployAction").
			 Where("version_id = ?", this.ID).
			 Last(vd).
			 Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return vd, err
}

type VersionDeployDetail struct {
	base.CRUD
	gorm.Model	
	Current      		 string      			`json:"current"     gorm:"comment:'当前阶段:test,green,approve,blue'"`
	AppDeploy    		 []AppDeploy 			`json:"app_deploy"  gorm:"many2many:version_deploy_app_deploy"`
	VersionDeployAction  []VersionDeployAction  `json:"version_deploy_action" gorm:"<-:false;foreignKey:VersionDeployID"`
}

func (VersionDeployDetail) TableName () string {
	return "version_deploy"
} 

type VersionDeploy struct {
	base.CRUD
	gorm.Model	
	Version      Version     `json:"version"`
	VersionID    int         `json:"version_id"  gorm:"comment:'版本id'"` 
	Current      string      `json:"current"     gorm:"comment:'当前阶段:test,green,approve,blue'"`
	AppDeploy    []AppDeploy `json:"app_deploy"  gorm:"many2many:version_deploy_app_deploy"`
	VersionDeployAction  []VersionDeployAction  `json:"version_deploy_action" gorm:"<-:false"`
}

type VersionDeployAction struct {
	base.CRUD
	gorm.Model
	VersionDeploy   VersionDeploy  `json:"version_deploy"`
	VersionDeployID int            `json:"version_deploy_id" gorm:"comment:'版本部署id'"`
	Step            string         `json:"step"              gorm:"type:varchar(16);comment:'部署阶段: dev,test,green,approve,blue'"`
	Status          string         `json:"status"            gorm:"type:varchar(16);comment:'状态: wait、process、finish、error'"`
	User            user.User      `json:"user"`
	UserID          int            `json:"user_id"           gorm:"comment:'执行者'"` 
}