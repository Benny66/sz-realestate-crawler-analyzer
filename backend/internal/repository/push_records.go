package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
)

// PushRecordRepository 推送记录仓库
type PushRecordRepository struct {
	config *config.Config
}

// NewPushRecordRepository 创建推送记录仓库
func NewPushRecordRepository(cfg *config.Config) *PushRecordRepository {
	return &PushRecordRepository{
		config: cfg,
	}
}

// Load 加载推送记录
func (r *PushRecordRepository) Load() ([]model.PushRecord, error) {
	data, err := os.ReadFile(r.config.Database.PushRecordsFile)
	if os.IsNotExist(err) {
		return []model.PushRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取推送记录失败: %v", err)
	}
	
	var records []model.PushRecord
	if err = json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析推送记录失败: %v", err)
	}
	
	return records, nil
}

// Save 保存推送记录
func (r *PushRecordRepository) Save(record *model.PushRecord) error {
	records, err := r.Load()
	if err != nil {
		records = []model.PushRecord{}
	}
	
	// 设置ID和时间戳
	if record.ID == "" {
		record.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	if record.PushedAt == "" {
		record.PushedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	
	records = append(records, *record)
	
	// 只保留最近1000条记录
	if len(records) > 1000 {
		records = records[len(records)-1000:]
	}
	
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化推送记录失败: %v", err)
	}
	
	return os.WriteFile(r.config.Database.PushRecordsFile, data, 0644)
}

// FindByFavoriteID 根据收藏ID查找推送记录
func (r *PushRecordRepository) FindByFavoriteID(favoriteID string) ([]model.PushRecord, error) {
	records, err := r.Load()
	if err != nil {
		return nil, err
	}
	
	filtered := make([]model.PushRecord, 0)
	for _, rec := range records {
		if rec.FavoriteID == favoriteID {
			filtered = append(filtered, rec)
		}
	}
	
	return filtered, nil
}

// FindRecent 查找最近的推送记录
func (r *PushRecordRepository) FindRecent(hours int) ([]model.PushRecord, error) {
	records, err := r.Load()
	if err != nil {
		return nil, err
	}
	
	cutoffTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	filtered := make([]model.PushRecord, 0)
	
	for _, rec := range records {
		recordTime, err := time.Parse("2006-01-02 15:04:05", rec.PushedAt)
		if err != nil {
			continue
		}
		
		if recordTime.After(cutoffTime) {
			filtered = append(filtered, rec)
		}
	}
	
	return filtered, nil
}