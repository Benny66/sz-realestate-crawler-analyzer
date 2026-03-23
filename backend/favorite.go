package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const favoriteFile = "favorites.json"

// LoadFavorites 加载收藏列表
func LoadFavorites() ([]FavoriteItem, error) {
	data, err := os.ReadFile(favoriteFile)
	if os.IsNotExist(err) {
		return []FavoriteItem{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取收藏文件失败: %v", err)
	}
	var items []FavoriteItem
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("解析收藏文件失败: %v", err)
	}
	return items, nil
}

// SaveFavorites 保存收藏列表
func SaveFavorites(items []FavoriteItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化收藏数据失败: %v", err)
	}
	return os.WriteFile(favoriteFile, data, 0644)
}

// AddFavorite 新增收藏
func AddFavorite(item FavoriteItem) error {
	items, err := LoadFavorites()
	if err != nil {
		items = []FavoriteItem{}
	}
	for _, it := range items {
		if it.ProjectName == item.ProjectName && it.BuildingName == item.BuildingName && it.HouseType == item.HouseType {
			return nil
		}
	}
	items = append(items, item)
	return SaveFavorites(items)
}

// DeleteFavorite 删除收藏
func DeleteFavorite(id string) error {
	items, err := LoadFavorites()
	if err != nil {
		return err
	}
	newItems := make([]FavoriteItem, 0, len(items))
	for _, it := range items {
		if it.ID != id {
			newItems = append(newItems, it)
		}
	}
	return SaveFavorites(newItems)
}
