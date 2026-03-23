package service

import (
	"fmt"
	"sync"
	"time"

	"sz-realestate-crawler-analyzer/internal/model"
)

// CacheEntry 缓存条目
type CacheEntry struct {
	Response  *model.AnalyzeResponse
	Timestamp time.Time
}

// CacheService 缓存服务
type CacheService struct {
	cache map[string]*CacheEntry
	mu    sync.RWMutex
}

// NewCacheService 创建缓存服务
func NewCacheService() *CacheService {
	return &CacheService{
		cache: make(map[string]*CacheEntry),
	}
}

// Get 获取缓存
func (cs *CacheService) Get(key string) *model.AnalyzeResponse {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	
	entry, exists := cs.cache[key]
	if !exists {
		return nil
	}
	
	// 检查缓存是否过期（30分钟）
	if time.Since(entry.Timestamp) > 30*time.Minute {
		return nil
	}
	
	return entry.Response
}

// Set 设置缓存
func (cs *CacheService) Set(key string, response *model.AnalyzeResponse) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	cs.cache[key] = &CacheEntry{
		Response:  response,
		Timestamp: time.Now(),
	}
}

// Clear 清除缓存
func (cs *CacheService) Clear() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	cs.cache = make(map[string]*CacheEntry)
}

// CacheKey 生成缓存键
func CacheKey(projectName, buildingName, houseType string) string {
	return fmt.Sprintf("%s_%s_%s", projectName, buildingName, houseType)
}