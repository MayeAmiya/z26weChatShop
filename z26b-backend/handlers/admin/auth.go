package admin

import (
	"log"
	"net/http"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLogin 管理员登录
func (h *Handler) AdminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写邮箱和密码"})
		return
	}

	var admin internal.Admin
	if err := h.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	if admin.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "账号已被禁用"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	token, _ := internal.GenerateAdminToken(admin.ID, admin.Email, admin.Username, admin.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"admin": gin.H{"id": admin.ID, "email": admin.Email, "username": admin.Username, "role": admin.Role},
	})
}

// AdminGetProfile 获取管理员信息
func (h *Handler) AdminGetProfile(c *gin.Context) {
	adminID := c.GetString("adminID")
	var admin internal.Admin
	if err := h.DB.First(&admin, "id = ?", adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": admin.ID, "email": admin.Email, "username": admin.Username, "role": admin.Role})
}

// AdminChangePassword 修改密码
func (h *Handler) AdminChangePassword(c *gin.Context) {
	adminID := c.GetString("adminID")
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入有效的密码"})
		return
	}

	var admin internal.Admin
	if err := h.DB.First(&admin, "id = ?", adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "管理员不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "原密码错误"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	admin.Password = string(hashedPassword)
	h.DB.Save(&admin)

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// InitDefaultAdmin 初始化默认管理员
func InitDefaultAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&internal.Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := internal.Admin{
		ID: internal.GenerateUUID(), Email: "admin@z26b.com",
		Password: string(hashedPassword), Username: "超级管理员",
		Role: "super_admin", Status: "active",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Println("Default admin initialized (admin@z26b.com / admin123)")
	return nil
}
