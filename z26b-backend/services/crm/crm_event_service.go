package crm

import (
	"time"

	"z26b-backend/internal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CRMEventService CRM事件服务
type CRMEventService struct {
	db *gorm.DB
}

// NewCRMEventService 创建CRM事件服务实例
func NewCRMEventService(db *gorm.DB) *CRMEventService {
	return &CRMEventService{db: db}
}

// RecordEvent 记录CRM事件
func (s *CRMEventService) RecordEvent(event *internal.CRMEvent) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt == 0 {
		event.CreatedAt = time.Now().UnixMilli()
	}
	return s.db.Create(event).Error
}

// GetEventsByUser 获取用户的事件列表
func (s *CRMEventService) GetEventsByUser(userID string, eventType string, page, pageSize int) ([]internal.CRMEvent, int64, error) {
	var events []internal.CRMEvent
	var total int64

	query := s.db.Model(&internal.CRMEvent{}).Where("user_id = ?", userID)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetEventsBySPU 获取商品的事件列表
func (s *CRMEventService) GetEventsBySPU(spuID string, eventType string, page, pageSize int) ([]internal.CRMEvent, int64, error) {
	var events []internal.CRMEvent
	var total int64

	query := s.db.Model(&internal.CRMEvent{}).Where("spu_id = ?", spuID)
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// GetEventStats 获取事件统计
func (s *CRMEventService) GetEventStats(startTime, endTime int64) (map[string]int64, error) {
	stats := make(map[string]int64)

	type Result struct {
		EventType string
		Count     int64
	}
	var results []Result

	query := s.db.Model(&internal.CRMEvent{}).
		Select("event_type, COUNT(*) as count").
		Group("event_type")

	if startTime > 0 {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		query = query.Where("created_at <= ?", endTime)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	for _, r := range results {
		stats[r.EventType] = r.Count
	}

	return stats, nil
}

// GetDailyEventStats 获取每日事件统计
func (s *CRMEventService) GetDailyEventStats(days int) ([]map[string]interface{}, error) {
	now := time.Now()
	var stats []map[string]interface{}

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UnixMilli()
		dayEnd := dayStart + 86400000

		eventStats, _ := s.GetEventStats(dayStart, dayEnd)

		// 计算当日各类事件数
		views := eventStats[internal.CRMEventTypeView]
		carts := eventStats[internal.CRMEventTypeCart]
		purchases := eventStats[internal.CRMEventTypePurchase]

		stats = append(stats, map[string]interface{}{
			"date":      date.Format("01-02"),
			"views":     views,
			"carts":     carts,
			"purchases": purchases,
			"total":     views + carts + purchases,
		})
	}

	return stats, nil
}
