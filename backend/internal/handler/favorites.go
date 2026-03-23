package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/internal/service"
)

// FavoriteHandler 收藏处理器
type FavoriteHandler struct {
	service *service.FavoriteService
}

// NewFavoriteHandler 创建收藏处理器
func NewFavoriteHandler(svc *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		service: svc,
	}
}

// List 获取收藏列表
func (h *FavoriteHandler) List(c *gin.Context) {
	favorites, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "加载收藏失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    favorites,
	})
}

// Create 新增收藏
func (h *FavoriteHandler) Create(c *gin.Context) {
	var req model.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "收藏请求解析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	if req.ProjectName == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "projectName 不能为空",
			Data:    nil,
		})
		return
	}
	
	favorite := model.FavoriteItem{
		ProjectName:  req.ProjectName,
		BuildingName: req.BuildingName,
		HouseType:    req.HouseType,
		Zone:         req.Zone,
		EnablePush:   req.EnablePush,
		PriceAlert:   req.PriceAlert,
		SaleAlert:    req.SaleAlert,
	}
	
	if err := h.service.Create(&favorite); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "保存收藏失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

// Delete 删除收藏
func (h *FavoriteHandler) Delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "id 参数不能为空",
			Data:    nil,
		})
		return
	}
	
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "删除收藏失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}

// Update 更新收藏
func (h *FavoriteHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "id 参数不能为空",
			Data:    nil,
		})
		return
	}
	
	var req model.FavoriteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Code:    1,
			Message: "更新请求解析失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	if err := h.service.Update(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Code:    1,
			Message: "更新收藏失败: " + err.Error(),
			Data:    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	})
}