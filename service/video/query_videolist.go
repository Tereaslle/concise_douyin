package video

import (
	"concise_douyin/cache"
	"concise_douyin/model"
	"errors"
)

type List struct {
	Videos []*model.Video `json:"video_list,omitempty"`
}

func QueryVideoListByUserID(userId int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(userId).Do()
}

func NewQueryVideoListByUserIdFlow(userID int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{userID: userID}
}

type QueryVideoListByUserIdFlow struct {
	userID int64
	videos []*model.Video

	videoList *List
}

func (q *QueryVideoListByUserIdFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryVideoListByUserIdFlow) checkNum() error {
	//检查userId是否存在
	if !model.NewUserInfoDAO().IsUserExistByID(q.userID) {
		return errors.New("用户不存在")
	}

	return nil
}

// 注意：Video由于在数据库中没有存储作者信息，所以需要手动填充
func (q *QueryVideoListByUserIdFlow) packData() error {
	err := model.NewVideoDAO().QueryVideoListByUserID(q.userID, &q.videos)
	if err != nil {
		return err
	}
	//作者信息查询
	var userInfo model.UserInfo
	err = model.NewUserInfoDAO().QueryUserInfoByID(q.userID, &userInfo)
	p := cache.NewProxyIndexMap()
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段
	for i := range q.videos {
		q.videos[i].Author = userInfo
		q.videos[i].IsFavorite = p.GetVideoFavorState(q.userID, q.videos[i].ID)
	}

	q.videoList = &List{Videos: q.videos}

	return nil
}
