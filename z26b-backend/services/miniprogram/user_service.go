package miniprogram

import (
	"errors"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserServiceInterface {
	return &UserService{db: db}
}

// GetOrCreateUser 根据OpenID获取或创建用户
func (s *UserService) GetOrCreateUser(openID string) (*internal.User, error) {
	var user internal.User
	err := s.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，创建新用户
			user = internal.User{
				ID:     internal.GenerateUUID(),
				OpenID: openID,
			}
			if err := s.db.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id string) (*internal.User, error) {
	var user internal.User
	err := s.db.First(&user, "id = ?", id).Error
	return &user, err
}

// UpdateUserInfo 更新用户信息
func (s *UserService) UpdateUserInfo(userID string, nickName, avatar string) error {
	return s.db.Model(&internal.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nick_name": nickName,
		"avatar":    avatar,
	}).Error
}

// GetUserInfo 获取用户信息（直接调用GetUserByID，避免冗余）
func (s *UserService) GetUserInfo(userID string) (*internal.User, error) {
	return s.GetUserByID(userID)
}
