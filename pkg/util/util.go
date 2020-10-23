package util

import (
	"github.com/gin-gonic/gin"
	"ops-backend/pkg/setting"
	"strconv"
	"strings"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

//获取URL中批量id并解析
func IdsStrToIdsIntGroup(key string, c *gin.Context) []int {
	return idsStrToIdsIntGroup(c.Param(key))
}

func idsStrToIdsIntGroup(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := strconv.Atoi(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}