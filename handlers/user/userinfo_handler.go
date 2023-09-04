package user

import (
	"concise_douyin/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserResponse struct {
	model.CommonResponse
	User *model.UserInfo `json:"user"`
}

func UserInfoHandler(c *gin.Context) {
  // util.GetServerIP(c.Request.Host)
	p := NewProxyUserInfo(c)
	//得到上层中间件根据token解析的userId
	rawId, ok := c.Get("user_id")
	if !ok {
		p.UserInfoError("解析userId出错")
		return
	}
	err := p.DoQueryUserInfoByUserId(rawId)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}

type ProxyUserInfo struct {
	c *gin.Context
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userID, ok := rawId.(int64)
	if !ok {
		return errors.New("解析userId失败")
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := model.NewUserInfoDAO()

	var userInfo model.UserInfo
	err := userinfoDAO.QueryUserInfoByID(userID, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOK(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: model.CommonResponse{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOK(user *model.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: model.CommonResponse{StatusCode: 0},
		User:           user,
	})
}
