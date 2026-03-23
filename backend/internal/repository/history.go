package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
)

// HistoryRepository 历史记录仓库
type HistoryRepository struct {
	config *config.Config
}

// NewHistoryRepository 创建历史记录仓库
func NewHistoryRepository(cfg *config.Config) *HistoryRepository {
	return &HistoryRepository{
		config: cfg,
	}
}

// Load 加载历史记录
func (r *HistoryRepository) Load() ([]model.HistoryRecord, error) {
	data, err := os.ReadFile(r.config.Database.HistoryFile)
	if os.IsNotExist(err) {
		return []model.HistoryRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取历史记录失败: %v", err)
	}
	
	var records []model.HistoryRecord
	if err = json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析历史记录失败: %v", err)
	}
	
	// 按创建时间排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt < records[j].CreatedAt
	})
	
	return records, nil
}

// Save 保存历史记录
func (r *HistoryRepository) Save(record *model.HistoryRecord) error {
	records, err := r.Load()
	if err != nil {
		records = []model.HistoryRecord{}
	}
	
	// 设置ID和时间戳
	if record.ID == "" {
		record.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	if record.CreatedAt == "" {
		record.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	
	records = append(records, *record)
	
	// 只保留最近1000条记录
	if len(records) > 1000 {
		records = records[len(records)-1000:]
	}
	
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史记录失败: %v", err)
	}
	
	return os.WriteFile(r.config.Database.HistoryFile, data, 0644)
}

// Delete 删除历史记录
func (r *HistoryRepository) Delete(id string) error {
	records, err := r.Load()
	if err != nil {
		return fmt.Errorf("加载历史记录失败: %v", err)
	}
	
	filtered := make([]model.HistoryRecord, 0)
	found := false
	
	for _, rec := range records {
		if rec.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, rec)
	}
	
	if !found {
		return fmt.Errorf("未找到ID为 %s 的历史记录", id)
	}
	
	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史记录失败: %v", err)
	}
	
	return os.WriteFile(r.config.Database.HistoryFile, data, 0644)
}

// FindByProject 根据楼盘名查找历史记录
func (r *HistoryRepository) FindByProject(projectName string) ([]model.HistoryRecord, error) {
	records, err := r.Load()
	if err != nil {
		return nil, err
	}
	
	filtered := make([]model.HistoryRecord, 0)
	for _, rec := range records {
		if rec.ProjectName == projectName {
			filtered = append(filtered, rec)
		}
	}
	
	return filtered, nil
}