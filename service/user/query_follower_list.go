package user

import (
	"concise_douyin/cache"
	"concise_douyin/model"
)

type FollowerList struct {
	UserList []*model.UserInfo `json:"user_list"`
}

func QueryFollowerList(userId int64) (*FollowerList, error) {
	return NewQueryFollowerListFlow(userId).Do()
}

type QueryFollowerListFlow struct {
	userID int64

	userList []*model.UserInfo

	*FollowerList
}

func NewQueryFollowerListFlow(userID int64) *QueryFollowerListFlow {
	return &QueryFollowerListFlow{userID: userID}
}

func (q *QueryFollowerListFlow) Do() (*FollowerList, error) {
	var err error
	if err = q.checkNum(); err != nil {
		return nil, err
	}
	if err = q.prepareData(); err != nil {
		return nil, err
	}
	if err = q.packData(); err != nil {
		return nil, err
	}
	return q.FollowerList, nil
}

func (q *QueryFollowerListFlow) checkNum() error {
	if !model.NewUserInfoDAO().IsUserExistByID(q.userID) {
		return ErrUserNotExist
	}
	return nil
}

func (q *QueryFollowerListFlow) prepareData() error {

	err := model.NewUserInfoDAO().GetFollowerListByUserID(q.userID, &q.userList)
	if err != nil {
		return err
	}
	//填充is_follow字段
	for _, v := range q.userList {
		v.IsFollow = cache.NewProxyIndexMap().GetUserRelation(q.userID, v.ID)
	}
	return nil
}

func (q *QueryFollowerListFlow) packData() error {
	q.FollowerList = &FollowerList{UserList: q.userList}

	return nil
}
