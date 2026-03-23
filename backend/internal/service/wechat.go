package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/internal/repository"
)

// WechatService 微信推送服务
type WechatService struct {
	config     *config.Config
	pushRepo   *repository.PushRecordRepository
	favoriteRepo *repository.FavoriteRepository
}

// NewWechatService 创建微信推送服务
func NewWechatService(cfg *config.Config) *WechatService {
	pushRepo := repository.NewPushRecordRepository(cfg)
	favoriteRepo := repository.NewFavoriteRepository(cfg)
	
	return &WechatService{
		config:       cfg,
		pushRepo:     pushRepo,
		favoriteRepo: favoriteRepo,
	}
}

// WechatMessage 微信消息结构
type WechatMessage struct {
	MsgType  string            `json:"msgtype"`
	Markdown *MarkdownContent  `json:"markdown,omitempty"`
}

type MarkdownContent struct {
	Content string `json:"content"`
}

// SendMarkdownMessage 发送Markdown消息
func (s *WechatService) SendMarkdownMessage(content string) error {
	message := WechatMessage{
		MsgType: "markdown",
		Markdown: &MarkdownContent{
			Content: content,
		},
	}
	
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}
	
	resp, err := http.Post(s.config.Wechat.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("微信接口返回错误: %s", resp.Status)
	}
	
	return nil
}

// BuildPriceAlertMessage 构建价格提醒消息
func (s *WechatService) BuildPriceAlertMessage(projectName, buildingName, houseType string, currentPrice, previousPrice, changePct float64) string {
	var emoji, trend string
	if changePct > 0 {
		emoji = "📈"
		trend = "上涨"
	} else {
		emoji = "📉"
		trend = "下降"
	}
	
	return fmt.Sprintf(`%s **价格变化提醒**
**楼盘**: %s %s %s
**当前价格**: %.0f 元/㎡
**上次价格**: %.0f 元/㎡
**变化幅度**: %s %.2f%%
**时间**: %s`, 
		emoji, projectName, buildingName, houseType, 
		currentPrice, previousPrice, trend, changePct,
		time.Now().Format("2006-01-02 15:04:05"))
}

// BuildSaleAlertMessage 构建销售状态提醒消息
func (s *WechatService) BuildSaleAlertMessage(projectName, buildingName, houseType string, currentRate, previousRate, changePct float64) string {
	var emoji, trend string
	if changePct > 0 {
		emoji = "🟢"
		trend = "增加"
	} else {
		emoji = "🔴"
		trend = "减少"
	}
	
	return fmt.Sprintf(`%s **销售状态提醒**
**楼盘**: %s %s %s
**当前在售率**: %.1f%%
**上次在售率**: %.1f%%
**变化幅度**: %s %.1f%%
**时间**: %s`, 
		emoji, projectName, buildingName, houseType,
		currentRate, previousRate, trend, changePct,
		time.Now().Format("2006-01-02 15:04:05"))
}

// BuildNewFavoriteMessage 构建新收藏通知消息
func (s *WechatService) BuildNewFavoriteMessage(projectName, buildingName, houseType, zone string) string {
	return fmt.Sprintf(`⭐ **新收藏楼盘**
**楼盘**: %s %s %s
**区域**: %s
**时间**: %s
**推送设置**: 价格变化≥3%% 销售状态≥5%%`, 
		projectName, buildingName, houseType, zone,
		time.Now().Format("2006-01-02 15:04:05"))
}

// CheckFavoriteAlerts 检查收藏楼盘是否需要推送提醒
func (s *WechatService) CheckFavoriteAlerts() error {
	favorites, err := s.favoriteRepo.Load()
	if err != nil {
		return fmt.Errorf("加载收藏失败: %v", err)
	}
	
	// 过滤启用推送的收藏
	activeFavorites := make([]model.FavoriteItem, 0)
	for _, fav := range favorites {
		if fav.EnablePush {
			activeFavorites = append(activeFavorites, fav)
		}
	}
	
	if len(activeFavorites) == 0 {
		return nil
	}
	
	// TODO: 实现价格和销售状态变化的检查逻辑
	// 这里需要结合历史记录数据进行分析
	
	fmt.Printf("[推送检查] 检查 %d 个收藏楼盘的变化...\n", len(activeFavorites))
	
	return nil
}

// StartPushMonitor 启动推送监控
func (s *WechatService) StartPushMonitor() {
	fmt.Println("[推送监控] 启动推送监控任务...")
	
	// 立即检查一次
	s.CheckFavoriteAlerts()
	
	// 每30分钟检查一次
	ticker := time.NewTicker(30 * time.Minute)
	
	go func() {
		for range ticker.C {
			fmt.Printf("[推送监控] 开始检查收藏楼盘变化 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
			if err := s.CheckFavoriteAlerts(); err != nil {
				fmt.Printf("[推送监控] 检查失败: %v\n", err)
			}
		}
	}()
}