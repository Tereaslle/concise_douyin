package video

import (
	"concise_douyin/config"
	"concise_douyin/model"
	"concise_douyin/service/video"
	"concise_douyin/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

// PublishVideoHandler 发布视频，并截取一帧画面作为封面
func PublishVideoHandler(c *gin.Context) {
  // util.GetServerIP(c.Request.Host)
	//准备参数
	rawId, _ := c.Get("user_id")

	userID, ok := rawId.(int64)
	if !ok {
		PublishVideoError(c, "解析UserId出错")
		return
	}

	title := c.PostForm("title")

	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}

	//支持多文件上传
	files := form.File["data"]
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)    //得到后缀
		if _, ok := videoIndexMap[suffix]; !ok { //判断是否为视频格式
			PublishVideoError(c, "不支持的视频格式")
			continue
		}
		//使用NewFileName生成“userID+workCount”的文件名
		coverName := util.NewFileName(userID)
		videoName := coverName + suffix //根据userID得到唯一的文件名
		savePath := filepath.Join(config.Global.StaticSourcePath, videoName)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		//截取一帧画面作为封面,到这里会生成完整的封面名（带后缀）
		err = util.SaveImageFromVideo(&coverName, false)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		//数据库持久化
		err := video.PostVideo(userID, videoName, coverName, title)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		PublishVideoOK(c, file.Filename+"上传成功")
	}
}

func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, model.CommonResponse{StatusCode: 1,
		StatusMsg: msg})
}

func PublishVideoOK(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, model.CommonResponse{StatusCode: 0, StatusMsg: msg})
}
