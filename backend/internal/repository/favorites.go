package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
)

// FavoriteRepository 收藏仓库
type FavoriteRepository struct {
	config *config.Config
}

// NewFavoriteRepository 创建收藏仓库
func NewFavoriteRepository(cfg *config.Config) *FavoriteRepository {
	return &FavoriteRepository{
		config: cfg,
	}
}

// Load 加载收藏列表
func (r *FavoriteRepository) Load() ([]model.FavoriteItem, error) {
	data, err := os.ReadFile(r.config.Database.FavoritesFile)
	if os.IsNotExist(err) {
		return []model.FavoriteItem{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取收藏列表失败: %v", err)
	}
	
	var favorites []model.FavoriteItem
	if err = json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("解析收藏列表失败: %v", err)
	}
	
	return favorites, nil
}

// Save 保存收藏列表
func (r *FavoriteRepository) Save(favorites []model.FavoriteItem) error {
	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化收藏列表失败: %v", err)
	}
	
	return os.WriteFile(r.config.Database.FavoritesFile, data, 0644)
}

// Create 创建收藏
func (r *FavoriteRepository) Create(favorite *model.FavoriteItem) error {
	favorites, err := r.Load()
	if err != nil {
		favorites = []model.FavoriteItem{}
	}
	
	// 检查是否已存在
	for _, fav := range favorites {
		if fav.ProjectName == favorite.ProjectName && 
		   fav.BuildingName == favorite.BuildingName && 
		   fav.HouseType == favorite.HouseType {
			return fmt.Errorf("该楼盘已存在于收藏列表中")
		}
	}
	
	// 设置ID和时间戳
	if favorite.ID == "" {
		favorite.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	if favorite.CreatedAt == "" {
		favorite.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	
	// 默认开启推送
	if !favorite.EnablePush {
		favorite.EnablePush = true
	}
	if !favorite.PriceAlert {
		favorite.PriceAlert = true
	}
	if !favorite.SaleAlert {
		favorite.SaleAlert = true
	}
	
	favorites = append(favorites, *favorite)
	
	return r.Save(favorites)
}

// Delete 删除收藏
func (r *FavoriteRepository) Delete(id string) error {
	favorites, err := r.Load()
	if err != nil {
		return fmt.Errorf("加载收藏列表失败: %v", err)
	}
	
	filtered := make([]model.FavoriteItem, 0)
	found := false
	
	for _, fav := range favorites {
		if fav.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, fav)
	}
	
	if !found {
		return fmt.Errorf("未找到ID为 %s 的收藏记录", id)
	}
	
	return r.Save(filtered)
}

// Update 更新收藏
func (r *FavoriteRepository) Update(id string, updateData map[string]interface{}) error {
	favorites, err := r.Load()
	if err != nil {
		return fmt.Errorf("加载收藏列表失败: %v", err)
	}
	
	updated := false
	for i, fav := range favorites {
		if fav.ID == id {
			if enablePush, ok := updateData["enablePush"].(bool); ok {
				favorites[i].EnablePush = enablePush
			}
			if priceAlert, ok := updateData["priceAlert"].(bool); ok {
				favorites[i].PriceAlert = priceAlert
			}
			if saleAlert, ok := updateData["saleAlert"].(bool); ok {
				favorites[i].SaleAlert = saleAlert
			}
			updated = true
			break
		}
	}
	
	if !updated {
		return fmt.Errorf("未找到对应的收藏记录")
	}
	
	return r.Save(favorites)
}