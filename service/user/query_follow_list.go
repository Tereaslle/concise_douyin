package user

import (
	"concise_douyin/model"
	"errors"
)

var (
	ErrUserNotExist = errors.New("用户不存在或已注销")
)

type FollowList struct {
	UserList []*model.UserInfo `json:"user_list"`
}

func QueryFollowList(userId int64) (*FollowList, error) {
	return NewQueryFollowListFlow(userId).Do()
}

type QueryFollowListFlow struct {
	userId int64

	userList []*model.UserInfo

	*FollowList
}

func NewQueryFollowListFlow(userId int64) *QueryFollowListFlow {
	return &QueryFollowListFlow{userId: userId}
}

func (q *QueryFollowListFlow) Do() (*FollowList, error) {
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

	return q.FollowList, nil
}

func (q *QueryFollowListFlow) checkNum() error {
	if !model.NewUserInfoDAO().IsUserExistByID(q.userId) {
		return ErrUserNotExist
	}
	return nil
}

func (q *QueryFollowListFlow) prepareData() error {
	var userList []*model.UserInfo
	err := model.NewUserInfoDAO().GetFollowListByUserID(q.userId, &userList)
	if err != nil {
		return err
	}
	for i, _ := range userList {
		userList[i].IsFollow = true //当前用户的关注列表，故isFollow定为true
	}
	q.userList = userList
	return nil
}

func (q *QueryFollowListFlow) packData() error {
	q.FollowList = &FollowList{UserList: q.userList}

	return nil
}
