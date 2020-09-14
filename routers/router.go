package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kataras/iris"
	"github.com/kukayyou/commonlib/mylog"
	"runtime"
	"userserver/controllers"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<12)
				runtime.Stack(buf, true)
				// 记录一个错误的日志
				mylog.Error("panic error : %s", string(buf))
				return
			}
		}()
		c.Next()
	}
}

func InitGinRouters() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(ErrHandler())
	root := ginRouter.Group("/userserver")
	//andriod、pc、mac用户
	userGroup := root.Group("/user")
	userGroup.POST("/registry", controllers.UserRegisterController{}.UserRegisterApi)
	userGroup.POST("/login", controllers.UserLoginController{}.UserLoginApi)
	userGroup.POST("/infos", controllers.UserInfosController{}.GetUserInfoApi)
	userGroup.POST("/update", controllers.UserInfosController{}.UpdateUserInfoApi)
	//微信小程序用户
	mini := root.Group("/miniuser")
	mini.POST("/login", controllers.MiniUserLoginController{}.MiniUserLoginApi)
	mini.POST("/infos", controllers.MiniUserInfosController{}.MiniUserInfosApi)
	return ginRouter
}

func InitIrisRouters() *iris.Application {
	irisRouter := iris.Default()
	root := irisRouter.Party("/userserver")
	//andriod、pc、mac用户
	userGroup := root.Party("/user")
	userGroup.Post("/registry", controllers.UserRegisterController{}.UserRegisterIrisApi)
	userGroup.Post("/login", controllers.UserLoginController{}.UserLoginIrisApi)
	//userGroup.Post("/infos", controllers.UserInfosController{}.GetUserInfoApi)
	//userGroup.Post("/update", controllers.UserInfosController{}.UpdateUserInfoApi)
	//微信小程序用户
	//mini := root.Party("/miniuser")
	//mini.Post("/login", controllers.MiniUserLoginController{}.MiniUserLoginApi)
	//mini.Post("/infos", controllers.MiniUserInfosController{}.MiniUserInfosApi)
	return irisRouter
}
