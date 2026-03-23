package handler

import (
	"net/http"
	"strconv"

	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/internal/service"

	"github.com/gin-gonic/gin"
)

// AnalyzeHandler 分析处理器
type AnalyzeHandler struct {
	service *service.AnalyzerService
}

// NewAnalyzeHandler 创建分析处理器
func NewAnalyzeHandler(svc *service.AnalyzerService) *AnalyzeHandler {
	return &AnalyzeHandler{
		service: svc,
	}
}

// SearchProjects 搜索楼盘
func (h *AnalyzeHandler) SearchProjects(c *gin.Context) {
	// 手动获取查询参数，避免 ShouldBindQuery 的问题
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "keyword 参数不能为空",
			Data:    nil,
		})
		return
	}

	pageIndex, _ := strconv.Atoi(c.Query("pageIndex"))
	if pageIndex <= 0 {
		pageIndex = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize <= 0 {
		pageSize = 12
	}

	zone := c.Query("zone")

	// 调用爬虫服务进行搜索
	result, err := h.service.SearchProjects(keyword, zone, pageIndex, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "搜索失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// GetBuildings 获取楼栋列表
func (h *AnalyzeHandler) GetBuildings(c *gin.Context) {
	// 手动获取查询参数
	ysProjectId, err := strconv.Atoi(c.Query("ysProjectId"))
	if err != nil || ysProjectId == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "ysProjectId 参数无效",
			Data:    nil,
		})
		return
	}

	preSellId, err := strconv.Atoi(c.Query("preSellId"))
	if err != nil || preSellId == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "preSellId 参数无效",
			Data:    nil,
		})
		return
	}

	// 调用爬虫服务获取楼栋列表
	result, err := h.service.GetBuildings(ysProjectId, preSellId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "获取楼栋列表失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// Analyze 执行分析
func (h *AnalyzeHandler) Analyze(c *gin.Context) {
	var req model.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "请求体解析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if req.Keyword == "" && req.YsProjectId == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "keyword 或 ysProjectId 不能同时为空",
			Data:    nil,
		})
		return
	}

	result, err := h.service.DoAnalyze(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// Compare 多楼盘对比
func (h *AnalyzeHandler) Compare(c *gin.Context) {
	var req model.CompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "请求体解析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "对比列表不能为空",
			Data:    nil,
		})
		return
	}

	if len(req.Items) > 5 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "最多支持 5 个楼盘同时对比",
			Data:    nil,
		})
		return
	}

	// TODO: 实现对比逻辑
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    model.CompareResponse{Items: []model.CompareItem{}, UpdatedAt: ""},
	})
}

// GetHistory 获取历史记录
func (h *AnalyzeHandler) GetHistory(c *gin.Context) {
	_ = c.Query("projectName") // 暂时忽略过滤参数

	// TODO: 实现获取历史记录逻辑
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    []model.HistoryRecord{},
	})
}

// DeleteHistory 删除历史记录
func (h *AnalyzeHandler) DeleteHistory(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "id 参数不能为空",
			Data:    nil,
		})
		return
	}

	// TODO: 实现删除历史记录逻辑
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

// ExportCSV 导出CSV
func (h *AnalyzeHandler) ExportCSV(c *gin.Context) {
	var req model.AnalyzeRequest

	// 解析查询参数
	req.Keyword = c.Query("keyword")
	req.BuildingName = c.Query("buildingName")
	req.HouseType = c.Query("houseType")
	req.Zone = c.Query("zone")
	req.YsProjectId, _ = strconv.Atoi(c.Query("ysProjectId"))
	req.PreSellId, _ = strconv.Atoi(c.Query("preSellId"))
	req.FybId, _ = strconv.Atoi(c.Query("fybId"))

	if req.Keyword == "" && req.YsProjectId == 0 {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "keyword 或 ysProjectId 不能为空",
			Data:    nil,
		})
		return
	}

	// 执行分析获取数据
	result, err := h.service.DoAnalyze(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "分析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 根据导出类型选择不同的导出方式
	exportType := c.Query("type")
	if exportType == "analysis" {
		// 导出分析结果
		err = h.service.ExportAnalysisCSV(c.Writer, result)
	} else {
		// 默认导出房源列表
		err = h.service.ExportHousesCSV(c.Writer, result)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "导出CSV失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
}

// ClearCache 清除缓存
func (h *AnalyzeHandler) ClearCache(c *gin.Context) {
	h.service.ClearCache()

	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "缓存已清除",
		Data:    nil,
	})
}
