package admin_services

import (
	"errors"
	"time"

	"z26b-backend/internal"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminAuthService struct {
	db *gorm.DB
}

func NewAdminAuthService(db *gorm.DB) AdminAuthServiceInterface {
	return &AdminAuthService{db: db}
}

// Login 管理员登录
func (s *AdminAuthService) Login(username, password string) (*internal.Admin, error) {
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

	// 更新最后登录时间
	now := time.Now()
	if err := s.db.Model(&admin).Update("last_login_at", now).Error; err != nil {
		// 记录错误但不影响登录
		// 可以考虑记录到日志中
	}

	return &admin, nil
}

// Logout 管理员登出（清理会话等）
func (s *AdminAuthService) Logout(adminUserID string) error {
	// 这里可以清理相关的会话或token
	// 目前暂时不需要实现
	return nil
}

// ChangePassword 修改密码
func (s *AdminAuthService) ChangePassword(adminUserID, oldPassword, newPassword string) error {
	// 验证新密码强度
	if err := internal.ValidatePassword(newPassword); err != nil {
		return err
	}

	var admin internal.Admin
	if err := s.db.Where("id = ?", adminUserID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	if err := s.db.Model(&admin).Updates(map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword 重置密码（管理员功能）
func (s *AdminAuthService) ResetPassword(adminUserID, newPassword string) error {
	// 验证新密码强度
	if err := internal.ValidatePassword(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	if err := s.db.Model(&internal.Admin{}).Where("id = ?", adminUserID).Updates(map[string]interface{}{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

// ValidateToken 验证token（这里简化处理，实际项目中可能需要JWT等）
func (s *AdminAuthService) ValidateToken(token string) (*internal.Admin, error) {
	// 这里应该实现token验证逻辑
	// 暂时返回错误，表示需要实现
	return nil, errors.New("token验证功能待实现")
}

// GetCurrentUser 获取当前登录用户信息
func (s *AdminAuthService) GetCurrentUser(adminUserID string) (*internal.Admin, error) {
	var admin internal.Admin
	if err := s.db.Where("id = ? AND status = ?", adminUserID, "active").First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或已被禁用")
		}
		return nil, err
	}
	return &admin, nil
}

// CheckPermission 检查用户权限
func (s *AdminAuthService) CheckPermission(adminUserID, permission string) (bool, error) {
	var admin internal.Admin
	if err := s.db.Where("id = ?", adminUserID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("用户不存在")
		}
		return false, err
	}

	// 超级管理员拥有所有权限
	if admin.Role == "super_admin" {
		return true, nil
	}

	// 根据角色检查权限
	// 这里可以实现更复杂的权限控制逻辑
	switch permission {
	case "user_management":
		return admin.Role == "admin" || admin.Role == "super_admin", nil
	case "product_management":
		return admin.Role == "admin" || admin.Role == "super_admin", nil
	case "order_management":
		return admin.Role == "admin" || admin.Role == "super_admin", nil
	case "stats_view":
		return admin.Role == "admin" || admin.Role == "super_admin", nil
	default:
		return false, nil
	}
}

// GetUserPermissions 获取用户权限列表
func (s *AdminAuthService) GetUserPermissions(adminUserID string) ([]string, error) {
	var admin internal.Admin
	if err := s.db.Where("id = ?", adminUserID).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	var permissions []string

	// 超级管理员拥有所有权限
	if admin.Role == "super_admin" {
		permissions = []string{
			"user_management",
			"product_management",
			"order_management",
			"stats_view",
			"system_management",
		}
	} else if admin.Role == "admin" {
		permissions = []string{
			"user_management",
			"product_management",
			"order_management",
			"stats_view",
		}
	}

	return permissions, nil
}