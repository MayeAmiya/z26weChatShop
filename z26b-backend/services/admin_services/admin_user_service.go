package admin_services

import (
	"errors"
	"time"

	"z26b-backend/internal"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminUserManagementService struct {
	db *gorm.DB
}

func NewAdminUserManagementService(db *gorm.DB) AdminUserManagementServiceInterface {
	return &AdminUserManagementService{db: db}
}

// GetAdminUsers 获取管理员用户列表
func (s *AdminUserManagementService) GetAdminUsers(page, pageSize int, keyword, status string) ([]internal.Admin, int64, error) {
	var adminUsers []internal.Admin
	var total int64

	query := s.db.Model(&internal.Admin{})

	if keyword != "" {
		query = query.Where("username LIKE ?", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&adminUsers).Error; err != nil {
		return nil, 0, err
	}

	return adminUsers, total, nil
}

// GetAdminUserByID 根据ID获取管理员用户详情
func (s *AdminUserManagementService) GetAdminUserByID(adminUserID string) (*internal.Admin, error) {
	var admin internal.Admin
	if err := s.db.Where("id = ?", adminUserID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("管理员用户不存在")
		}
		return nil, err
	}
	return &admin, nil
}

// UpdateAdminUser 更新管理员用户信息
func (s *AdminUserManagementService) UpdateAdminUser(adminUserID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	if err := s.db.Model(&internal.Admin{}).Where("id = ?", adminUserID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAdminUser 删除管理员用户
func (s *AdminUserManagementService) DeleteAdminUser(adminUserID string) error {
	if err := s.db.Where("id = ?", adminUserID).Delete(&internal.Admin{}).Error; err != nil {
		return err
	}
	return nil
}

// CreateAdminUser 创建管理员用户
func (s *AdminUserManagementService) CreateAdminUser(username, password, role string) (*internal.Admin, error) {
	// 检查用户名是否已存在
	var existingAdmin internal.Admin
	if err := s.db.Where("username = ?", username).First(&existingAdmin).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := internal.Admin{
		Username:  username,
		Password:  string(hashedPassword),
		Role:      role,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

// AuthenticateAdmin 管理员登录验证
func (s *AdminUserManagementService) AuthenticateAdmin(username, password string) (*internal.Admin, error) {
	var admin internal.Admin
	if err := s.db.Where("username = ? AND status = ?", username, "active").First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return &admin, nil
}