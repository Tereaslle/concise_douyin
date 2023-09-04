package video

import (
	"concise_douyin/model"
)

// PostVideo 投稿视频
func PostVideo(userID int64, videoName, coverName, title string) error {
	return NewPostVideoFlow(userID, videoName, coverName, title).Do()
}

func NewPostVideoFlow(userID int64, videoName, coverName, title string) *PostVideoFlow {
	return &PostVideoFlow{
		videoName: videoName,
		coverName: coverName,
		userID:    userID,
		title:     title,
	}
}

type PostVideoFlow struct {
	videoName string
	coverName string
	title     string
	userID    int64

	video *model.Video
}

func (f *PostVideoFlow) Do() error {
	//f.prepareParam()

	if err := f.publish(); err != nil {
		return err
	}
	return nil
}

//// 准备好参数
//func (f *PostVideoFlow) prepareParam() {
//	f.videoName = util.GetFileURL(f.videoName)
//	f.coverName = util.GetFileURL(f.coverName)
//}

// 组合并添加到数据库
func (f *PostVideoFlow) publish() error {
	video := &model.Video{
		UserInfoID: f.userID,
		PlayURL:    f.videoName,
		CoverURL:   f.coverName,
		Title:      f.title,
	}
	return model.NewVideoDAO().AddVideo(video)
}
