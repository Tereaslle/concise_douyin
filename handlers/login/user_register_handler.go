package login

import (
	"concise_douyin/model"
	user_login2 "concise_douyin/service/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterResponse struct {
	model.CommonResponse
	*user_login2.LoginResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	rawVal, _ := c.Get("password")
	password, ok := rawVal.(string)
	if !ok {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: model.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析出错",
			},
		})
		return
	}
	registerResponse, err := user_login2.PostUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: model.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: model.CommonResponse{StatusCode: 0},
		LoginResponse:  registerResponse,
	})
}
