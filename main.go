package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// crawlHouseInfo 爬取并分析房源信息
func crawlHouseInfo() error {
	fmt.Printf("\n开始爬取房源信息 - %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 1. 自动解析楼盘参数（通过关键字搜索获取 ysProjectId / preSellId / fybId）
	params, err := ResolveProjectParams(RequestConfig.ProjectKeyword)
	if err != nil {
		return fmt.Errorf("解析楼盘参数失败: %v", err)
	}

	// 2. 爬取楼栋字典（确认楼栋信息）
	_, err = postForm(DefaultConfig.BuildingDictURL, map[string]string{
		"ysProjectId": fmt.Sprintf("%d", params.YsProjectId),
		"preSellId":   fmt.Sprintf("%d", params.PreSellId),
		"fybId":       fmt.Sprintf("%d", params.FybId),
	})
	if err != nil {
		return fmt.Errorf("爬取楼栋字典失败: %v", err)
	}

	// 3. 爬取房源列表
	houseInfoResp, err := postJSON(DefaultConfig.HouseInfoURL, HouseInfoRequest{
		Buildingbranch: "",
		Floor:          "",
		FybId:          fmt.Sprintf("%d", params.FybId),
		Housenb:        "",
		Status:         "-1",
		Type:           RequestConfig.HouseType,
		YsProjectId:    params.YsProjectId,
		PreSellId:      params.PreSellId,
	})
	if err != nil {
		return fmt.Errorf("爬取房源列表失败: %v", err)
	}

	// 4. 解析响应
	var houseInfoResult HouseInfoResponse
	if err = json.Unmarshal(houseInfoResp, &houseInfoResult); err != nil {
		return fmt.Errorf("房源列表数据JSON解析失败: %v", err)
	}
	fmt.Printf("成功获取 %d 个楼层的房源数据\n", len(houseInfoResult.Data))

	// 5. 分析并打印报告
	result := Analyze(houseInfoResult.Data)
	saleSummary := AnalyzeSaleStatus(houseInfoResult.Data)
	PrintReport(result, saleSummary)
	return nil
}

// scheduleTask 定时任务
func scheduleTask() {
	if err := crawlHouseInfo(); err != nil {
		fmt.Printf("首次爬取失败: %v\n", err)
	}
	ticker := time.NewTicker(time.Duration(DefaultConfig.Interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := crawlHouseInfo(); err != nil {
			fmt.Printf("定时爬取失败: %v\n", err)
		}
	}
}

func main() {
	fmt.Println("深圳房产信息爬取脚本启动！")
	fmt.Printf("爬取间隔: %d 秒\n", DefaultConfig.Interval)
	fmt.Println("==============================")

	// 启动前先初始化 Cookie
	fmt.Println("正在自动获取 Cookie...")
	if err := initCookieClient(); err != nil {
		fmt.Printf("⚠️  自动获取 Cookie 失败，将使用 config.go 中的静态 Cookie: %v\n", err)
	}

	scheduleTask()
}
