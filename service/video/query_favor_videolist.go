package video

import (
	"concise_douyin/model"
	"errors"
)

type FavorList struct {
	Videos []*model.Video `json:"video_list"`
}

func QueryFavorVideoList(userID int64) (*FavorList, error) {
	return NewQueryFavorVideoListFlow(userID).Do()
}

type QueryFavorVideoListFlow struct {
	userID int64

	videos []*model.Video

	videoList *FavorList
}

func NewQueryFavorVideoListFlow(userID int64) *QueryFavorVideoListFlow {
	return &QueryFavorVideoListFlow{userID: userID}
}

func (q *QueryFavorVideoListFlow) Do() (*FavorList, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryFavorVideoListFlow) checkNum() error {
	if !model.NewUserInfoDAO().IsUserExistByID(q.userID) {
		return errors.New("用户状态异常")
	}
	return nil
}

func (q *QueryFavorVideoListFlow) prepareData() error {
	err := model.NewVideoDAO().QueryFavorVideoListByUserID(q.userID, &q.videos)
	if err != nil {
		return err
	}
	//填充信息(Author和IsFavorite字段，由于是点赞列表，故所有的都是点赞状态
	for i := range q.videos {
		//作者信息查询
		var userInfo model.UserInfo
		err = model.NewUserInfoDAO().QueryUserInfoByID(q.videos[i].UserInfoID, &userInfo)
		if err == nil { //若查询未出错则更新，否则不更新作者信息
			q.videos[i].Author = userInfo
		}
		q.videos[i].IsFavorite = true
	}
	return nil
}

func (q *QueryFavorVideoListFlow) packData() error {
	q.videoList = &FavorList{Videos: q.videos}
	return nil
}
