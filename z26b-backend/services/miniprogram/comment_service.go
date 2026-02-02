package miniprogram

import (
	"z26b-backend/internal"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService(db *gorm.DB) CommentServiceInterface {
	return &CommentService{db: db}
}

// GetGoodsComments 获取商品评论
func (s *CommentService) GetGoodsComments(spuID string, page, pageSize int) ([]internal.Comment, int64, error) {
	var comments []internal.Comment
	var total int64

	query := s.db.Where("spu_id = ?", spuID)
	err := query.Model(&internal.Comment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(comment *internal.Comment) error {
	return s.db.Create(comment).Error
}

// GetCommentByID 根据ID获取评论
func (s *CommentService) GetCommentByID(id string) (*internal.Comment, error) {
	var comment internal.Comment
	err := s.db.First(&comment, "id = ?", id).Error
	return &comment, err
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(id string, updates map[string]interface{}) error {
	return s.db.Model(&internal.Comment{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id string) error {
	return s.db.Delete(&internal.Comment{}, "id = ?", id).Error
}

// GetCommentsByUser 获取用户的评论
func (s *CommentService) GetCommentsByUser(userID string, page, pageSize int) ([]internal.Comment, int64, error) {
	var comments []internal.Comment
	var total int64

	query := s.db.Where("user_id = ?", userID)
	err := query.Model(&internal.Comment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&comments).Error
	return comments, total, err
}
