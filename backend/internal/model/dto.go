package model

// APIResponse 统一 API 响应包装
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// AnalyzeRequest 前端发起分析请求体
type AnalyzeRequest struct {
	Keyword      string `json:"keyword"`
	BuildingName string `json:"buildingName"`
	HouseType    string `json:"houseType"`
	Zone         string `json:"zone"`
	YsProjectId  int    `json:"ysProjectId"`
	PreSellId    int    `json:"preSellId"`
	FybId        int    `json:"fybId"`
}

// AnalyzeResponse 分析结果完整响应
type AnalyzeResponse struct {
	Analysis    AnalysisResult        `json:"analysis"`
	Sale        SaleSummary           `json:"sale"`
	Houses      []HouseItem           `json:"houses"`
	FloorGroups []FloorGroup          `json:"floorGroups"`
	Params      ResolvedProjectParams `json:"params"`
	Alerts      []PriceAlert          `json:"alerts"`
	UpdatedAt   string                `json:"updatedAt"`
}

// ProjectListResponse 楼盘搜索结果响应
type ProjectListResponse struct {
	Total    int           `json:"total"`
	Projects []ProjectItem `json:"projects"`
}

// BuildingListResponse 楼栋列表响应
type BuildingListResponse struct {
	Buildings []BuildingNameItem `json:"buildings"`
}

// CompareRequest 多楼盘对比请求体
type CompareRequest struct {
	Items []AnalyzeRequest `json:"items"`
}

// CompareItem 单个楼盘对比数据
type CompareItem struct {
	ProjectName  string         `json:"projectName"`
	BuildingName string         `json:"buildingName"`
	HouseType    string         `json:"houseType"`
	Zone         string         `json:"zone"`
	Analysis     AnalysisResult `json:"analysis"`
	Sale         SaleSummary    `json:"sale"`
}

// CompareResponse 多楼盘对比响应
type CompareResponse struct {
	Items     []CompareItem `json:"items"`
	UpdatedAt string        `json:"updatedAt"`
}

// FavoriteRequest 收藏请求
type FavoriteRequest struct {
	ProjectName  string `json:"projectName"`
	BuildingName string `json:"buildingName"`
	HouseType    string `json:"houseType"`
	Zone         string `json:"zone"`
	EnablePush   bool   `json:"enablePush"`
	PriceAlert   bool   `json:"priceAlert"`
	SaleAlert    bool   `json:"saleAlert"`
}

// FavoriteUpdateRequest 收藏更新请求
type FavoriteUpdateRequest struct {
	EnablePush *bool `json:"enablePush"`
	PriceAlert *bool `json:"priceAlert"`
	SaleAlert  *bool `json:"saleAlert"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword   string `json:"keyword"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	Zone      string `json:"zone"`
}

// BuildingsRequest 楼栋列表请求
type BuildingsRequest struct {
	YsProjectId int `json:"ysProjectId"`
	PreSellId   int `json:"preSellId"`
}