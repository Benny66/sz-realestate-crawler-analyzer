package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const pushRecordsFile = "push_records.json"

// LoadPushRecords 加载推送记录
func LoadPushRecords() ([]PushRecord, error) {
	data, err := os.ReadFile(pushRecordsFile)
	if os.IsNotExist(err) {
		return []PushRecord{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("读取推送记录失败: %v", err)
	}
	var records []PushRecord
	if err = json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析推送记录失败: %v", err)
	}
	return records, nil
}

// SavePushRecord 保存推送记录
func SavePushRecord(record PushRecord) error {
	records, err := LoadPushRecords()
	if err != nil {
		records = []PushRecord{}
	}
	
	record.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
	record.PushedAt = time.Now().Format("2006-01-02 15:04:05")
	records = append(records, record)
	
	// 只保留最近1000条记录
	if len(records) > 1000 {
		records = records[len(records)-1000:]
	}
	
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化推送记录失败: %v", err)
	}
	return os.WriteFile(pushRecordsFile, data, 0644)
}

// CheckFavoriteAlerts 检查收藏楼盘是否需要推送提醒
func CheckFavoriteAlerts() {
	favorites, err := LoadFavorites()
	if err != nil {
		fmt.Printf("[推送检查] 加载收藏失败: %v\n", err)
		return
	}
	
	records, err := LoadHistory()
	if err != nil {
		fmt.Printf("[推送检查] 加载历史记录失败: %v\n", err)
		return
	}
	
	for _, fav := range favorites {
		if !fav.EnablePush {
			continue
		}
		
		// 过滤该楼盘的历史记录
		filtered := make([]HistoryRecord, 0)
		for _, r := range records {
			if r.ProjectName == fav.ProjectName && r.BuildingName == fav.BuildingName && r.HouseType == fav.HouseType {
				filtered = append(filtered, r)
			}
		}
		
		if len(filtered) < 2 {
			continue
		}
		
		last := filtered[len(filtered)-1]
		prev := filtered[len(filtered)-2]
		
		// 检查价格变化
		if fav.PriceAlert && last.Analysis.MedianUnitPrice > 0 && prev.Analysis.MedianUnitPrice > 0 {
			changePct := (last.Analysis.MedianUnitPrice - prev.Analysis.MedianUnitPrice) / prev.Analysis.MedianUnitPrice * 100
			
			// 价格变化超过3%才推送
			if changePct >= 3 || changePct <= -3 {
				message := BuildPriceAlertMessage(
					fav.ProjectName, fav.BuildingName, fav.HouseType,
					last.Analysis.MedianUnitPrice, prev.Analysis.MedianUnitPrice, changePct,
				)
				
				if err := SendWechatMarkdownMessage(message); err != nil {
					fmt.Printf("[推送检查] 发送价格提醒失败: %v\n", err)
				} else {
					fmt.Printf("[推送检查] 已发送价格提醒: %s\n", fav.ProjectName)
					_ = SavePushRecord(PushRecord{
						FavoriteID: fav.ID,
						Type:       "price",
						ChangePct:  changePct,
						Message:    message,
					})
				}
			}
		}
		
		// 检查销售状态变化
		if fav.SaleAlert && last.Sale.ForSaleRate > 0 && prev.Sale.ForSaleRate > 0 {
			changePct := last.Sale.ForSaleRate - prev.Sale.ForSaleRate
			
			// 在售率变化超过5%才推送
			if changePct <= -5 || changePct >= 5 {
				message := BuildSaleAlertMessage(
					fav.ProjectName, fav.BuildingName, fav.HouseType,
					last.Sale.ForSaleRate, prev.Sale.ForSaleRate, changePct,
				)
				
				if err := SendWechatMarkdownMessage(message); err != nil {
					fmt.Printf("[推送检查] 发送销售提醒失败: %v\n", err)
				} else {
					fmt.Printf("[推送检查] 已发送销售提醒: %s\n", fav.ProjectName)
					_ = SavePushRecord(PushRecord{
						FavoriteID: fav.ID,
						Type:       "sale",
						ChangePct:  changePct,
						Message:    message,
					})
				}
			}
		}
	}
}

// StartPushMonitor 启动推送监控
func StartPushMonitor() {
	fmt.Println("[推送监控] 启动推送监控任务...")
	
	// 立即检查一次
	CheckFavoriteAlerts()
	
	// 每30分钟检查一次
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		fmt.Printf("[推送监控] 开始检查收藏楼盘变化 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
		CheckFavoriteAlerts()
	}
}
