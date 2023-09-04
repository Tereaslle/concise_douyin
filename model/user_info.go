package model

import (
	"concise_douyin/config"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	ErrIvdPtr        = errors.New("空指针错误")
	ErrEmptyUserList = errors.New("用户列表为空")
)

type UserInfo struct {
	ID              int64       `json:"id" gorm:"omitempty"`                                                      // 用户id
	Name            string      `json:"name" gorm:"omitempty"`                                                    // 用户名
	FollowCount     int64       `json:"follow_count" gorm:"omitempty"`                                            // 关注总数
	FollowerCount   int64       `json:"follower_count" gorm:"omitempty"`                                          // 粉丝总数
	Avatar          string      `json:"avatar" gorm:"default:'avatar-default.jpg'"`                     // 用户头像
	BackgroundImage string      `json:"background_image" gorm:"default:'background_image-default.jpg'"` //用户个人页顶部大图
	Signature       string      `json:"signature" gorm:"default:'这个人很懒，什么都没有写'"`                        //个人简介
	IsFollow        bool        `json:"is_follow" gorm:"omitempty"`                                               // true-已关注，false-未关注
	TotalFavorited  int64       `json:"total_favorited" gorm:"omitempty,default:0"`                               //获赞数量
	WorkCount       int64       `json:"work_count" gorm:"omitempty,default:0"`                                    //作品数量
	FavoriteCount   int64       `json:"favorite_count" gorm:"omitempty,default:0"`                                //点赞数量
	User            *UserLogin  `json:"-"`                                                                        //用户与账号密码之间的一对一
	Videos          []*Video    `json:"-"`                                                                        //用户与投稿视频的一对多
	Follows         []*UserInfo `json:"-" gorm:"many2many:user_relations;"`                                       //用户之间的多对多
	FavorVideos     []*Video    `json:"-" gorm:"many2many:user_favor_videos;"`                                    //用户与点赞视频之间的多对多
	Comments        []*Comment  `json:"-"`                                                                        //用户与评论的一对多
}

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

func GenerateCompleteUserInfoURL(user *UserInfo) {
	user.Avatar = fmt.Sprintf("https://%s/static/%s", config.Global.IP,  user.Avatar)
	user.BackgroundImage = fmt.Sprintf("https://%s/static/%s", config.Global.IP, user.BackgroundImage)
}

// UpdateTotalFavorited 更新用户的获赞数量
func (u *UserInfoDAO) UpdateTotalFavorited(userID int64) error {
	var videos []Video
	count := int64(0)
	// 找到该用户的所有视频，把每个视频收到的赞favorite_count求和即可
	database.Model(&Video{}).Where("user_info_id = ?", 1).Find(&videos)
	//求和
	for _, v := range videos {
		count += v.FavoriteCount
		//fmt.Printf("%v\n", v.FavoriteCount)
	}
	result := database.Model(&UserInfo{}).Where("id=?", userID).Update("total_favorited", count)
	return result.Error
}

// UpdateWorkCount 更新用户的作品数量
func (u *UserInfoDAO) UpdateWorkCount(userID int64) error {
	var count int64
	err := NewVideoDAO().QueryVideoCountByUserID(userID, &count) // 获取该用户ID的作品数量
	if err != nil {
		return err
	}
	result := database.Model(&UserInfo{}).Where("id=?", userID).Update("work_count", count)
	return result.Error
}

func (u *UserInfoDAO) QueryUserInfoByID(userID int64, userInfo *UserInfo) error {
	if userInfo == nil {
		return ErrIvdPtr
	}
	*userInfo = UserInfo{ID: userID}
	database.First(userInfo)
	GenerateCompleteUserInfoURL(userInfo)
	//id为零值，说明sql执行失败
	if userInfo.ID == 0 {
		return errors.New("该用户不存在")
	}
	return nil
}

func (u *UserInfoDAO) AddUserInfo(userinfo *UserInfo) error {
	if userinfo == nil {
		return ErrIvdPtr
	}
	return database.Create(userinfo).Error
}

func (u *UserInfoDAO) IsUserExistByID(id int64) bool {
	var userinfo UserInfo
	if err := database.Where("id=?", id).Select("id").First(&userinfo).Error; err != nil {
		log.Println(err)
	}
	if userinfo.ID == 0 {
		return false
	}
	return true
}
func (u *UserInfoDAO) AddUserFollow(userId, userToId int64) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count+1 WHERE id = ?", userToId).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO `user_relations` (`user_info_id`,`follow_id`) VALUES (?,?)", userId, userToId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserInfoDAO) CancelUserFollow(userId, userToId int64) error {
	return database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", userToId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_relations` WHERE user_info_id=? AND follow_id=?", userId, userToId).Error; err != nil {
			return err
		}
		return nil
	})
}

func (u *UserInfoDAO) GetFollowListByUserID(userId int64, userList *[]*UserInfo) error {
	if userList == nil {
		return ErrIvdPtr
	}
	var err error
	if err = database.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.user_info_id = ? AND r.follow_id = u.id", userId).Scan(userList).Error; err != nil {
		return err
	}
	if len(*userList) == 0 || (*userList)[0].ID == 0 {
		return ErrEmptyUserList
	}
	for _, user := range *userList {
		GenerateCompleteUserInfoURL(user)
	}
	return nil
}

func (u *UserInfoDAO) GetFollowerListByUserID(userId int64, userList *[]*UserInfo) error {
	if userList == nil {
		return ErrIvdPtr
	}
	var err error
	if err = database.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.follow_id = ? AND r.user_info_id = u.id", userId).Scan(userList).Error; err != nil {
		return err
	}
	//if len(*userList) == 0 || (*userList)[0].Id == 0 {
	//	return ErrEmptyUserList
	//}
	for _, user := range *userList {
		GenerateCompleteUserInfoURL(user)
	}
	return nil
}
