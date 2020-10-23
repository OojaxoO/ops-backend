package user

import (
	"strconv"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"ops-backend/middleware/jwt"
	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	"ops-backend/pkg/util"
	"ops-backend/view/user/post"
	"ops-backend/view/user/role"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	appG := app.Gin{C: c}
	auth := Auth{Username: username, Password: password}
	isExist, err := auth.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

func Info(c *gin.Context) {
	appG := app.Gin{C: c}
	token, err := jwt.GetToken(c)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, err)
		return
	}
	claims, err := util.ParseToken(token)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, err)
		return
	}
	var user User
	username, err := util.Base64Decode(claims.Username)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, err)
		return
	}
	password, err := util.Base64Decode(claims.Password)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, err)
		return
	}
	user.Username = username
	user.Password = password
	if err = user.Detail(); err != nil {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

/**
 * @Author chenjun
 * @Description //创建用户信息
 * @Date 3:16 下午 2020/9/11
 **/
func CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	if err := c.MustBindWith(&user, binding.JSON); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}
	id, err := user.CreateUser()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_USER_REPEAT, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, id)
	return
}

/**
 * @Author chenjun
 * @Description //更新用户信息
 * @Date 3:16 下午 2020/9/11
 **/
func UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	if err := c.MustBindWith(&user, binding.JSON); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}
	err := user.UpdateUser(user.ID)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_USER_UPDATE, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "更新成功")
	return
}

/**
 * @Author chenjun
 * @Description //删除用户信息
 * @Date 3:15 下午 2020/9/11
 **/
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	IdS := util.IdsStrToIdsIntGroup("Id", c)
	fmt.Println(c)
	fmt.Println(IdS)
	result, err := user.DeleteUser(IdS)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_USER_DELETE, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, result)
	return
}

/**
 * @Author chenjun
 * @Description //获取用户信息
 * @Date 3:15 下午 2020/9/11
 **/
func GetUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var user User
	if err := c.MustBindWith(&user, binding.JSON); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}
	result, err := user.SelectUser()
	if err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR_USER_SELECT, err)
		return
	}
	var Role role.Role
	var Post post.Post
	roles, err := Role.GetList()
	posts, err := Post.GetList()
	postIds := make([]int, 0)
	postIds = append(postIds, result.PostId)
	roleIds := make([]int, 0)
	roleIds = append(roleIds, result.RoleId)
	appG.Response(http.StatusOK, e.SUCCESS,gin.H{
		"data":    result,
		"postIds": postIds,
		"roleIds": roleIds,
		"roles":   roles,
		"posts":   posts,

	})
	return
}
/**
 * @Author chenjun
 * @Description //获取用户角色和职位
 * @Date 3:52 下午 2020/9/12
 **/
func GetUserInit(c *gin.Context)  {
	appG := app.Gin{C: c}
	var Role role.Role
	var Post post.Post
	roles, err := Role.GetList()
	posts, err := Post.GetList()
	if err!=nil{
		appG.Response(http.StatusBadRequest, e.ERROR_USER_SELECT, err)
	}
	mp:=make(map[string]interface{},2)
	mp["roles"] = roles
	mp["posts"] = posts
	appG.Response(http.StatusOK, e.ERROR_USER_SELECT, mp)
	
}
func GetSearch(search string, searchField []string) string {
	if len(search) == 0 {
		return ""
	}
	newField := []string{}
	for _, field := range searchField {
		like := field + " like '%" + search + "%'" 
		newField = append(newField, like)
	}
	return strings.Join(newField, " or ") 
}

func List(c *gin.Context) {
	// _type := reflect.TypeOf(this.Data)
	searchField := []string{"username"}
	preload := []string{}
	fieldFilter := new(User) 
	// _type_list := reflect.MakeSlice(reflect.SliceOf(_type), 0, 0)
	data := new([]User) 

	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(fieldFilter); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	queryPage := c.DefaultQuery("page", "0")
	queryPageSize := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")
	filter := GetSearch(search, searchField)

	page, err := strconv.Atoi(queryPage)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	pageSize, err := strconv.Atoi(queryPageSize)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

    // 返回全部
	if page == 0 {
		if err := fieldFilter.Find(data, fieldFilter, filter, preload); err != nil { 
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		} 
		appG.Response(http.StatusOK, e.SUCCESS, data) 
		return
	}

	// 分页
	response, err := fieldFilter.List(data, page, pageSize, filter, fieldFilter, preload)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, response) 		
}