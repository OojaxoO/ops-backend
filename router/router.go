package router

import (
	"ops-backend/view/user/post"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ops-backend/pkg/setting"
	"ops-backend/view"
	"ops-backend/view/asset"
	"ops-backend/view/asset/host"
	"ops-backend/view/kube"
	"ops-backend/view/user"
)

func Setup() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTION", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	public := r.Group("/")

	public.POST("/login", user.Login)

	private := r.Group("/api")
	//private.Use(jwt.JWT())

	private.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	private.GET("/asset", asset.Info)
	private.GET("/asset/host", host.List)

	//用户系统路由
	private.GET("/user/info", user.Info)
	private.POST("/user", user.CreateUser)
	private.PUT("/user",user.UpdateUser)
	// private.GET("/user/:ID",user.GetUser)
	private.DELETE("/user/:ID",user.DeleteUser)
	// private.GET("/user",user.GetUserInit)
	private.GET("/user",user.List)

	//岗位系统路由
	private.GET("/post",post.UpdatePost)
	private.POST("/post", post.InsertPost)
	private.PUT("/post/:postId", post.UpdatePost)
	private.DELETE("/post/:postId", post.DeletePost)


	view.RegistAllRouter(private)
	private.GET("/kube/cluster/:ID/:resource", kube.List)
	port := setting.HttpSetting.Port
	r.Run(":" + port)
}
