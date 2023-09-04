package video

import (
	"concise_douyin/model"
	"concise_douyin/service/video"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FavorVideoListResponse struct {
	model.CommonResponse
	*video.FavorList
}

func QueryFavorVideoListHandler(c *gin.Context) {
	NewProxyFavorVideoListHandler(c).Do()
}

type ProxyFavorVideoListHandler struct {
	*gin.Context

	userID int64
}

func NewProxyFavorVideoListHandler(c *gin.Context) *ProxyFavorVideoListHandler {
  // util.GetServerIP(c.Request.Host)
	return &ProxyFavorVideoListHandler{Context: c}
}

func (p *ProxyFavorVideoListHandler) Do() {
	//解析参数
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}

	//正式调用
	favorVideoList, err := video.QueryFavorVideoList(p.userID)
	if err != nil {
		p.SendError(err.Error())
		return
	}

	//成功返回
	p.SendOK(favorVideoList)
}

func (p *ProxyFavorVideoListHandler) parseNum() error {
	rawUserID, _ := p.Get("user_id")
	userID, ok := rawUserID.(int64)
	if !ok {
		return errors.New("userID解析出错")
	}
	p.userID = userID
	return nil
}

func (p *ProxyFavorVideoListHandler) SendError(msg string) {
	p.JSON(http.StatusOK, FavorVideoListResponse{
		CommonResponse: model.CommonResponse{StatusCode: 1, StatusMsg: msg}})
}

func (p *ProxyFavorVideoListHandler) SendOK(favorList *video.FavorList) {
	p.JSON(http.StatusOK, FavorVideoListResponse{CommonResponse: model.CommonResponse{StatusCode: 0},
		FavorList: favorList,
	})
}
