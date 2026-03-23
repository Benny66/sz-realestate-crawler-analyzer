package service

import (
	"fmt"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/internal/repository"
)

// FavoriteService 收藏服务
type FavoriteService struct {
	config *config.Config
	repo   *repository.FavoriteRepository
	wechat *WechatService
}

// NewFavoriteService 创建收藏服务
func NewFavoriteService(cfg *config.Config) *FavoriteService {
	repo := repository.NewFavoriteRepository(cfg)
	wechat := NewWechatService(cfg)

	return &FavoriteService{
		config: cfg,
		repo:   repo,
		wechat: wechat,
	}
}

// List 获取收藏列表
func (s *FavoriteService) List() ([]model.FavoriteItem, error) {
	return s.repo.Load()
}

// Create 创建收藏
func (s *FavoriteService) Create(favorite *model.FavoriteItem) error {
	if err := s.repo.Create(favorite); err != nil {
		return err
	}

	// 发送新收藏通知
	go func() {
		message := s.wechat.BuildNewFavoriteMessage(
			favorite.ProjectName,
			favorite.BuildingName,
			favorite.HouseType,
			favorite.Zone,
		)
		if err := s.wechat.SendMarkdownMessage(message); err != nil {
			fmt.Printf("发送新收藏通知失败: %v\n", err)
		} else {
			fmt.Printf("已发送新收藏通知: %s\n", favorite.ProjectName)
		}
	}()

	return nil
}

// Delete 删除收藏
func (s *FavoriteService) Delete(id string) error {
	return s.repo.Delete(id)
}

// Update 更新收藏
func (s *FavoriteService) Update(id string, req model.FavoriteUpdateRequest) error {
	updateData := make(map[string]interface{})

	if req.EnablePush != nil {
		updateData["enablePush"] = *req.EnablePush
	}
	if req.PriceAlert != nil {
		updateData["priceAlert"] = *req.PriceAlert
	}
	if req.SaleAlert != nil {
		updateData["saleAlert"] = *req.SaleAlert
	}

	return s.repo.Update(id, updateData)
}
