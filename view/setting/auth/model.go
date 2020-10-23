package auth

import (
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Name   string  `json:"name" gorm:"type:varchar(64);comment:'名字'"`
	Type   string  `json:"type" gorm:"type:varchar(16);comment:'类型: jenkins,ssh,git等'"`  // ssh, http, tcp, udp, mysql, redis, https
	Arg    string  `json:"arg"  gorm:"type:text;comment:'类型参数'"`  	      // 参数	
}