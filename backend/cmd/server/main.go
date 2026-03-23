package main

import (
	"fmt"
	"log"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/handler"
	"sz-realestate-crawler-analyzer/internal/service"
	"sz-realestate-crawler-analyzer/pkg/gin"
)

func main() {
	fmt.Println("========================================")
	fmt.Println(" 深圳房产信息智能分析平台 - 后端服务 (Gin)")
	fmt.Println("========================================")

	// 1. 加载配置
	cfg := config.Load()
	fmt.Printf("监听端口: %s\n", cfg.Server.Addr)
	fmt.Printf("定时刷新间隔: %v\n", cfg.Crawler.Interval)

	// 2. 初始化服务
	analyzerService := service.NewAnalyzerService(cfg)
	favoriteService := service.NewFavoriteService(cfg)

	// 3. 启动后台任务
	analyzerService.StartBackgroundTasks()

	// 4. 初始化处理器
	analyzeHandler := handler.NewAnalyzeHandler(analyzerService)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	// 5. 设置路由
	r := gin.SetupRouter(analyzeHandler, favoriteHandler)

	// 6. 启动 HTTP Server
	fmt.Printf("API 服务已启动: http://localhost%s\n", cfg.Server.Addr)
	fmt.Println("前端访问地址: http://localhost:5173")
	fmt.Println("健康检查: http://localhost:8080/health")
	fmt.Println("后台任务: 定时刷新 + 微信推送监控已启动")

	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
