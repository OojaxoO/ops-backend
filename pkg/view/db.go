package view

import (
	"strconv"
	"strings"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	"ops-backend/pkg/models"
)

type DBObject struct {
	Data interface{}
	Uri string
	SearchField []string
}

func (this *DBObject) GetUri() string {
	return this.Uri 
}

func (this *DBObject) GetSearch(search string) string {
	if len(search) == 0 {
		return ""
	}
	newField := []string{}
	for _, field := range this.SearchField {
		like := field + " like '%" + search + "%'" 
		newField = append(newField, like)
	}
	return strings.Join(newField, " or ") 
}

func (this *DBObject) List(c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	fieldFilter := reflect.New(_type).Interface()
	_type_list := reflect.MakeSlice(reflect.SliceOf(_type), 0, 0)
	data := reflect.New(_type_list.Type()).Interface()

	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(fieldFilter); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	queryPage := c.DefaultQuery("page", "0")
	queryPageSize := c.DefaultQuery("pageSize", "10")
	search := c.DefaultQuery("search", "")
	filter := this.GetSearch(search)

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
		if err := models.Find(data, fieldFilter, filter); err != nil { 
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		} 
		appG.Response(http.StatusOK, e.SUCCESS, data) 
		return
	}

	// 分页
	response, err := models.List(page, pageSize, data, fieldFilter, filter)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, response) 		
}
 
func (this *DBObject) Create(c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	err := c.BindJSON(data)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if err := models.Create(data); err != nil { 
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

func (this *DBObject) Update(c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := models.Get(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := c.BindJSON(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := models.Update(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

func (this *DBObject) Delete (c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := models.Get(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := models.Delete(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除成功") 		
}

func (this *DBObject) Get (c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := models.Get(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}
