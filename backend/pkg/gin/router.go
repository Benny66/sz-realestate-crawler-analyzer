package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sz-realestate-crawler-analyzer/internal/handler"
	"sz-realestate-crawler-analyzer/internal/model"
)

// SetupRouter 设置路由
func SetupRouter(analyzeHandler *handler.AnalyzeHandler, favoriteHandler *handler.FavoriteHandler) *gin.Engine {
	r := gin.Default()
	
	// 中间件
	r.Use(CORSMiddleware())
	r.Use(LoggerMiddleware())
	
	// API路由分组
	api := r.Group("/api")
	{
		// 搜索相关
		api.GET("/search", analyzeHandler.SearchProjects)
		api.GET("/buildings", analyzeHandler.GetBuildings)
		api.POST("/analyze", analyzeHandler.Analyze)
		api.POST("/compare", analyzeHandler.Compare)
		api.GET("/export/csv", analyzeHandler.ExportCSV)
		
		// 历史记录
		api.GET("/history", analyzeHandler.GetHistory)
		api.DELETE("/history", analyzeHandler.DeleteHistory)
		
		// 收藏相关
		favorites := api.Group("/favorites")
		{
			favorites.GET("", favoriteHandler.List)
			favorites.POST("", favoriteHandler.Create)
			favorites.DELETE("", favoriteHandler.Delete)
			favorites.PUT("/:id", favoriteHandler.Update)
		}
		
		// 缓存管理
		api.POST("/cache/clear", analyzeHandler.ClearCache)
	}
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.APIResponse{
			Code:    0,
			Message: "服务运行正常",
			Data:    nil,
		})
	})
	
	return r
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求日志
		c.Next()
	}
}