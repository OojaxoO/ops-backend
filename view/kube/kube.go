package kube 

import (
	"github.com/gin-gonic/gin"
	viewKube "ops-backend/pkg/kube"
)

func List (c *gin.Context) {
	viewKube.List(c)	
}