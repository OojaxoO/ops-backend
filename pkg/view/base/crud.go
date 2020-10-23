package base 

import (
	"strconv"
	"strings"
	"net/http"
	"reflect"
	"fmt"

	"github.com/gin-gonic/gin"

	"ops-backend/pkg/app"
	"ops-backend/pkg/e"
	. "ops-backend/pkg/models/base"
)

type DBView interface {
	List(*gin.Context)
	Update(*gin.Context)
	Create(*gin.Context)
	Delete(*gin.Context)
	Get(*gin.Context)
	GetUri() string
}

type BaseView struct {
	Data        CRUDInterface 
	Uri         string
	SearchField []string
	Preload     []string
}

func (this *BaseView) GetUri() string {
	return this.Uri 
}

func (this *BaseView) GetSearch(search string) string {
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

func (this *BaseView) List(c *gin.Context) {
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
		if err := this.Data.Find(data, fieldFilter, filter, this.Preload); err != nil { 
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		} 
		appG.Response(http.StatusOK, e.SUCCESS, data) 
		return
	}

	// 分页
	response, err := this.Data.List(data, page, pageSize, filter, fieldFilter, this.Preload)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, response) 		
}
 
func (this *BaseView) Create(c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	fmt.Print(_type)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	err := c.BindJSON(data)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if err := this.Data.Create(data); err != nil { 
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if err := this.Data.Get(data, this.Preload); err != nil { 
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

func (this *BaseView) Update(c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := this.Data.Get(data, nil); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := c.BindJSON(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := this.Data.Update(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	if err := this.Data.Get(data, this.Preload); err != nil { 
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

func (this *BaseView) Delete (c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := this.Data.Get(data, nil); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := this.Data.Delete(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, "删除成功") 		
}

func (this *BaseView) Get (c *gin.Context) {
	_type := reflect.TypeOf(this.Data)
	data := reflect.New(_type).Interface()
	appG := app.Gin{C: c}
	if err := c.ShouldBindUri(data); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := this.Data.Get(data, this.Preload); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, data) 		
}

