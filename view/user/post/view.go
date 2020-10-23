/**
 * @Author chenjun
 * @Date 2020/9/11
 **/
package post

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	"ops-backend/pkg/util"
)

func GetPostList(c *gin.Context) {
	return
}

/**
 * @Author chenjun
 * @Description //获取岗位信息
 * @Date 5:50 下午 2020/9/12
 **/
func GetPost(c *gin.Context) {
	appG := app.Gin{C: c}
	var post Post
	if err := c.MustBindWith(&post, binding.JSON); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, "获取成功")

}

/**
 * @Author chenjun
 * @Description //添加岗位
 * @Date 6:14 下午 2020/9/12
 **/
func InsertPost(c *gin.Context) {
	appG := app.Gin{C: c}
	var post Post
	err := c.Bind(&post)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
	}
	result, err := post.Create()
	appG.Response(http.StatusOK, e.SUCCESS, result)
}

/**
 * @Author chenjun
 * @Description //TODO
 * @Date 6:34 下午 2020/9/12
 **/
func UpdatePost(c *gin.Context) {
	appG := app.Gin{C: c}
	var post Post
	err := c.Bind(&post)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
	}
	result, err := post.Update(post.PostId)
	appG.Response(http.StatusOK, e.SUCCESS, result)
}
func DeletePost(c *gin.Context) {
	appG := app.Gin{C: c}
	var post Post
	IDS:=util.IdsStrToIdsIntGroup("postId",c)
	result, err :=post.BatchDelete(IDS)
	if err!=nil{
		appG.Response(http.StatusBadRequest, e.ERROR_USER_DELETE, err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, result)
	return
}
