package host

import (
	"github.com/gin-gonic/gin"

	"ops-backend/pkg/view"
)

func List (c *gin.Context) {
	searchField := []string{"region_id", "host_name", "os_name", "public_ip_address", "primary_ip_address"}
	ob := view.DBObject {Data: Host{}, SearchField: searchField}
	ob.List(c)
}