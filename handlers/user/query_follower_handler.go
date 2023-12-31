package user

import (
	"concise_douyin/model"
	"concise_douyin/service/user"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FollowerListResponse struct {
	model.CommonResponse
	*user.FollowerList
}

func QueryFollowerHandler(c *gin.Context) {
	NewProxyQueryFollowerHandler(c).Do()
}

type ProxyQueryFollowerHandler struct {
	*gin.Context

	userId int64

	*user.FollowerList
}

func NewProxyQueryFollowerHandler(c *gin.Context) *ProxyQueryFollowerHandler {
  // util.GetServerIP(c.Request.Host)
	return &ProxyQueryFollowerHandler{Context: c}
}

func (p *ProxyQueryFollowerHandler) Do() {
	var err error
	if err = p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.prepareData(); err != nil {
		if errors.Is(err, user.ErrUserNotExist) {
			p.SendError(err.Error())
		} else {
			p.SendError("准备数据出错")
		}
		return
	}
	p.SendOK("成功")
}

func (p *ProxyQueryFollowerHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId
	return nil
}

func (p *ProxyQueryFollowerHandler) prepareData() error {
	list, err := user.QueryFollowerList(p.userId)
	if err != nil {
		return err
	}
	p.FollowerList = list
	return nil
}

func (p *ProxyQueryFollowerHandler) SendError(msg string) {
	p.JSON(http.StatusOK, FollowerListResponse{
		CommonResponse: model.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *ProxyQueryFollowerHandler) SendOK(msg string) {
	p.JSON(http.StatusOK, FollowerListResponse{
		CommonResponse: model.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FollowerList: p.FollowerList,
	})
}
