package model

// 实体定义，无需导入包

// HouseItem 单个房源数据结构
type HouseItem struct {
	Housenb                    string  `json:"housenb"`
	Floor                      string  `json:"floor"`
	YsbuildingArea             float64 `json:"ysbuildingarea"`
	YsExpandArea               float64 `json:"ysexpandarea"`
	YsInsideArea               float64 `json:"ysinsidearea"`
	AskPriceTotalB             float64 `json:"askpricetotalB"`
	AskPriceEachB              float64 `json:"askpriceeachB"`
	RecordedPricePerUnitInside float64 `json:"recordedPricePerUnitInside"`
	LastStatusName             string  `json:"lastStatusName"`
	BuildingName               string  `json:"buildingName"`
	UseAge                     string  `json:"useage"`
	BarrierFree                string  `json:"barrierFree"`
}

// FloorGroup 按楼层分组的房源数据
type FloorGroup struct {
	Floor string      `json:"floor"`
	List  []HouseItem `json:"list"`
}

// BuildingNameItem 单个楼栋信息
type BuildingNameItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Key   string `json:"key"`
	Acive bool   `json:"acive"`
}

// ProjectItem 单个楼盘信息
type ProjectItem struct {
	Id          string `json:"id"`
	SypId       string `json:"sypId"`
	SypeId      string `json:"sypeId"`
	Zone        string `json:"zone"`
	PreSellNo   string `json:"strpreprojectid"`
	ProjectName string `json:"project"`
	CompanyName string `json:"name"`
	SiteAddress string `json:"siteaddress"`
	PassDate    string `json:"passdate"`
	ImagePath   string `json:"imagePath"`
}

// AnalysisResult 房源数据量化分析结果
type AnalysisResult struct {
	MinUnitPrice           float64 `json:"min_unit_price"`
	MaxUnitPrice           float64 `json:"max_unit_price"`
	MedianUnitPrice        float64 `json:"median_unit_price"`
	MinRecordedUnitPrice   float64 `json:"min_recorded_unit_price"`
	MaxRecordedUnitPrice   float64 `json:"max_recorded_unit_price"`
	MinTotalPrice          float64 `json:"min_total_price"`
	MaxTotalPrice          float64 `json:"max_total_price"`
	AvgTotalPrice          float64 `json:"avg_total_price"`
	TotalPriceSpan         float64 `json:"total_price_span"`
	BuildingArea           float64 `json:"building_area"`
	ExpandArea             float64 `json:"expand_area"`
	InsideArea             float64 `json:"inside_area"`
	ActualUseArea          float64 `json:"actual_use_area"`
	HousingRate            float64 `json:"housing_rate"`
	TotalCount             int     `json:"total_count"`
	MinFloor               int     `json:"min_floor"`
	MaxFloor               int     `json:"max_floor"`
	LowFloorCount          int     `json:"low_floor_count"`
	MidFloorCount          int     `json:"mid_floor_count"`
	HighFloorCount         int     `json:"high_floor_count"`
	MaxPerFloor            int     `json:"max_per_floor"`
	MinPerFloor            int     `json:"min_per_floor"`
	UnitPricePerExpandArea float64 `json:"unit_price_per_expand_area"`
	CostPerActualArea      float64 `json:"cost_per_actual_area"`
	FloorPricePremium      float64 `json:"floor_price_premium"`
}

// SaleStatusStats 销售状态统计
type SaleStatusStats struct {
	StatusName string `json:"statusName"`
	Count      int    `json:"count"`
}

// SaleSummary 销售概况汇总
type SaleSummary struct {
	TotalCount    int               `json:"total_count"`
	ForSaleCount  int               `json:"for_sale_count"`
	SoldCount     int               `json:"sold_count"`
	StatusDetails []SaleStatusStats `json:"status_details"`
	SoldRate      float64           `json:"sold_rate"`
	ForSaleRate   float64           `json:"for_sale_rate"`
}

// ResolvedProjectParams 从搜索结果中解析出的目标楼盘参数
type ResolvedProjectParams struct {
	YsProjectId  int    `json:"ysProjectId"`
	PreSellId    int    `json:"preSellId"`
	FybId        int    `json:"fybId"`
	ProjectName  string `json:"projectName"`
	BuildingName string `json:"buildingName"`
	HouseType    string `json:"houseType"`
	AutoSelected bool   `json:"autoSelected"`
}

// PriceAlert 价格/销售变化预警
type PriceAlert struct {
	Level     string  `json:"level"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	Direction string  `json:"direction"`
	ChangePct float64 `json:"changePct"`
}

// HistoryRecord 历史爬取记录
type HistoryRecord struct {
	ID           string         `json:"id"`
	ProjectName  string         `json:"projectName"`
	BuildingName string         `json:"buildingName"`
	HouseType    string         `json:"houseType"`
	Zone         string         `json:"zone"`
	Analysis     AnalysisResult `json:"analysis"`
	Sale         SaleSummary    `json:"sale"`
	CreatedAt    string         `json:"createdAt"`
}

// FavoriteItem 收藏楼盘
type FavoriteItem struct {
	ID           string `json:"id"`
	ProjectName  string `json:"projectName"`
	BuildingName string `json:"buildingName"`
	HouseType    string `json:"houseType"`
	Zone         string `json:"zone"`
	EnablePush   bool   `json:"enablePush"`
	PriceAlert   bool   `json:"priceAlert"`
	SaleAlert    bool   `json:"saleAlert"`
	CreatedAt    string `json:"createdAt"`
}

// PushRecord 推送记录
type PushRecord struct {
	ID         string  `json:"id"`
	FavoriteID string  `json:"favoriteId"`
	Type       string  `json:"type"`
	ChangePct  float64 `json:"changePct"`
	Message    string  `json:"message"`
	PushedAt   string  `json:"pushedAt"`
}
