package router

import (
	"concise_douyin/config"
	"concise_douyin/handlers/comment"
	"concise_douyin/handlers/login"
	"concise_douyin/handlers/user"
	"concise_douyin/handlers/video"
	"concise_douyin/middleware"
	"concise_douyin/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	//初始化数据库
	model.InitDB()
	//创建一个默认的路由
	r := gin.Default()

	r.Static("static", config.Global.StaticSourcePath)
	r.LoadHTMLGlob("./template/*")
	
	// 为主页面加载html文件
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"IP": config.Global.IP})
	})
	
	baseGroup := r.Group("/douyin")
	// 绑定路由规则和路由函数，访问对应的路由，将由对应的函数去处理
	// 根据灵活性考虑是否加入JWT中间件来进行鉴权，还是在之后再做鉴权
	// 基础api定义
	//视频流
	baseGroup.GET("/feed/", video.FeedVideoListHandler)
	baseGroup.GET("/user/", middleware.JWTMiddleWare(), user.UserInfoHandler)
	//用户登陆/注册功能
	baseGroup.POST("/user/login/", middleware.SHAMiddleWare(), login.UserLoginHandler)
	baseGroup.POST("/user/register/", middleware.SHAMiddleWare(), login.UserRegisterHandler)
	//视频上传功能
	baseGroup.POST("/publish/action/", middleware.JWTMiddleWare(), video.PublishVideoHandler)
	baseGroup.GET("/publish/list/", middleware.NoAuthToGetUserId(), video.QueryVideoListHandler)
	//视频点赞功能
	baseGroup.POST("/favorite/action/", middleware.JWTMiddleWare(), video.PostFavorHandler)
	baseGroup.GET("/favorite/list/", middleware.NoAuthToGetUserId(), video.QueryFavorVideoListHandler)
	//评论功能
	baseGroup.POST("/comment/action/", middleware.JWTMiddleWare(), comment.PostCommentHandler)
	baseGroup.GET("/comment/list/", middleware.JWTMiddleWare(), comment.QueryCommentListHandler)
	//用户关注功能
	baseGroup.POST("/relation/action/", middleware.JWTMiddleWare(), user.PostFollowActionHandler)
	baseGroup.GET("/relation/follow/list/", middleware.NoAuthToGetUserId(), user.QueryFollowListHandler)
	baseGroup.GET("/relation/follower/list/", middleware.NoAuthToGetUserId(), user.QueryFollowerHandler)

	//返回路由至main函数，由main函数启动监听
	return r
}
