package login

import (
	"concise_douyin/model"
	"concise_douyin/service/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	model.CommonResponse
	*login.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")
	password, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: model.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
	}
	userLoginResponse, err := login.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: model.CommonResponse{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: model.CommonResponse{StatusCode: 0},
		LoginResponse:  userLoginResponse,
	})
}
