package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const historyFile = "history.json"

// CacheEntry 单次分析结果缓存
type CacheEntry struct {
	Response  AnalyzeResponse
	CachedAt  time.Time
}

// analysisCache 内存缓存：key = "projectName_buildingName_houseType"
var (
	analysisCache   = make(map[string]*CacheEntry)
	analysisCacheMu sync.RWMutex
)

// cacheKey 生成缓存 key
func cacheKey(projectName, buildingName, houseType string) string {
	return fmt.Sprintf("%s_%s_%s", projectName, buildingName, houseType)
}

// SetCache 写入缓存
func SetCache(projectName, buildingName, houseType string, resp AnalyzeResponse) {
	key := cacheKey(projectName, buildingName, houseType)
	analysisCacheMu.Lock()
	defer analysisCacheMu.Unlock()
	analysisCache[key] = &CacheEntry{
		Response: resp,
		CachedAt: time.Now(),
	}
}

// GetCache 读取缓存，返回 (entry, found)
func GetCache(projectName, buildingName, houseType string) (*CacheEntry, bool) {
	key := cacheKey(projectName, buildingName, houseType)
	analysisCacheMu.RLock()
	defer analysisCacheMu.RUnlock()
	entry, ok := analysisCache[key]
	return entry, ok
}

// IsCacheValid 判断缓存是否在有效期内（默认 30 分钟）
func IsCacheValid(entry *CacheEntry, ttl time.Duration) bool {
	return time.Since(entry.CachedAt) < ttl
}

// ---- 历史记录持久化 ----

// LoadHistory 从文件加载历史记录
func LoadHistory() ([]HistoryRecord, error) {
	data, err := os.ReadFile(historyFile)
	if os.IsNotExist(err) {
		return []HistoryRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取历史文件失败: %v", err)
	}
	var records []HistoryRecord
	if err = json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析历史文件失败: %v", err)
	}
	return records, nil
}

// SaveHistory 追加一条历史记录到文件
func SaveHistory(record HistoryRecord) error {
	records, err := LoadHistory()
	if err != nil {
		records = []HistoryRecord{}
	}
	// 最多保留 500 条，超出则删除最旧的
	records = append(records, record)
	if len(records) > 500 {
		records = records[len(records)-500:]
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化历史记录失败: %v", err)
	}
	return os.WriteFile(historyFile, data, 0644)
}

// DeleteHistory 删除指定 ID 的历史记录
func DeleteHistory(id string) error {
	records, err := LoadHistory()
	if err != nil {
		return err
	}
	newRecords := make([]HistoryRecord, 0, len(records))
	for _, r := range records {
		if r.ID != id {
			newRecords = append(newRecords, r)
		}
	}
	data, err := json.MarshalIndent(newRecords, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(historyFile, data, 0644)
}
