package main

import (
	"fmt"
	"net/http"
	"time"
)

// scheduleAutoRefresh 定时自动刷新默认楼盘缓存（修复原有 return 导致定时器不生效的 bug）
func scheduleAutoRefresh() {
	fmt.Println("[定时任务] 预热默认楼盘缓存...")
	_, err := doAnalyze(AnalyzeRequest{
		Keyword:      RequestConfig.ProjectKeyword,
		BuildingName: RequestConfig.BuildingName,
		HouseType:    RequestConfig.HouseType,
	})
	if err != nil {
		fmt.Printf("[定时任务] 预热失败（不影响服务启动）: %v\n", err)
	} else {
		fmt.Println("[定时任务] 预热成功")
	}

	ticker := time.NewTicker(time.Duration(DefaultConfig.Interval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		fmt.Printf("[定时任务] 开始刷新缓存 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
		_, err := doAnalyze(AnalyzeRequest{
			Keyword:      RequestConfig.ProjectKeyword,
			BuildingName: RequestConfig.BuildingName,
			HouseType:    RequestConfig.HouseType,
		})
		if err != nil {
			fmt.Printf("[定时任务] 刷新失败: %v\n", err)
		} else {
			fmt.Println("[定时任务] 刷新成功")
		}
	}
}

func main() {
	fmt.Println("========================================")
	fmt.Println(" 深圳房产信息智能分析平台 - 后端服务")
	fmt.Println("========================================")
	fmt.Printf("监听端口: %s\n", DefaultConfig.ServerAddr)
	fmt.Printf("定时刷新间隔: %d 秒\n", DefaultConfig.Interval)

	// 1. 初始化 Cookie
	fmt.Println("正在自动获取 Cookie...")
	if err := initCookieClient(); err != nil {
		fmt.Printf("自动获取 Cookie 失败，将使用静态 Cookie: %v\n", err)
	}

	// 2. 注册路由
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	// 3. 后台启动定时刷新任务
	go scheduleAutoRefresh()
	// 4. 启动推送监控（订阅收藏楼盘的价格/销售变化）
	go StartPushMonitor()

	// 4. 启动 HTTP Server
	fmt.Printf("API 服务已启动: http://localhost%s\n", DefaultConfig.ServerAddr)
	fmt.Println("前端访问地址: http://localhost:5173")
	if err := http.ListenAndServe(DefaultConfig.ServerAddr, mux); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
