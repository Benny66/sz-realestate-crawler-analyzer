package main

// HouseItem 单个房源数据结构（对应实际JSON字段）
type HouseItem struct {
	Housenb                    string  `json:"housenb"`                    // 房号（如"2902"）
	Floor                      string  `json:"floor"`                      // 楼层（字符串类型，如"29"）
	YsbuildingArea             float64 `json:"ysbuildingarea"`             // 产权建筑面积（㎡）
	YsExpandArea               float64 `json:"ysexpandarea"`               // 赠送/拓展面积（㎡），不计入产权但可实际使用
	YsInsideArea               float64 `json:"ysinsidearea"`               // 套内面积（㎡），即产权建筑面积去除公摊后的面积
	AskPriceTotalB             float64 `json:"askpricetotalB"`             // 挂牌总价（元）
	AskPriceEachB              float64 `json:"askpriceeachB"`              // 挂牌单价（元/㎡），按产权建筑面积计算
	RecordedPricePerUnitInside float64 `json:"recordedPricePerUnitInside"` // 备案单价（元/㎡），按套内面积计算，通常高于挂牌单价
	LastStatusName             string  `json:"lastStatusName"`             // 房源当前状态名称（如"期房待售"）
	BuildingName               string  `json:"buildingName"`               // 所属楼栋名称（如"1栋"）
	UseAge                     string  `json:"useage"`                     // 用途类型（如"住宅"）
	BarrierFree                string  `json:"barrierFree"`                // 是否无障碍房源（"是"/"否"）
}

// FloorGroup 按楼层分组的房源数据，对应接口返回的 data 数组中每个元素
type FloorGroup struct {
	Floor string      `json:"floor"` // 楼层编号（字符串类型，如"29"）
	List  []HouseItem `json:"list"`  // 该楼层下的所有房源列表
}

// HouseInfoResponse 房源列表接口响应结构
type HouseInfoResponse struct {
	Data   []FloorGroup `json:"data"`   // 按楼层分组的房源数据列表
	Msg    string       `json:"msg"`    // 接口返回消息（如"成功"）
	Status int          `json:"status"` // 接口返回状态码（200表示成功）
}

// BuildingDictResponse 楼栋字典接口响应结构
type BuildingDictResponse struct {
	Data struct {
		BuildNo   string   `json:"buildNo"`   // 楼栋编号（全国房屋唯一编码）
		ChList    []string `json:"chList"`    // 可售楼层列表（字符串数组，如["2","3",...]）
		GnqmcList []string `json:"gnqmcList"` // 功能区名称列表（如["未知"]，对应楼栋单元/分区名称）
		Zrzid     int      `json:"zrzid"`     // 自然资源ID
	} `json:"data"`
	Msg    string `json:"msg"`    // 接口返回消息（如"成功"）
	Status int    `json:"status"` // 接口返回状态码（200表示成功）
}

// BuildingNameItem 单个楼栋信息（对应 getBuildingNameListToPublicity 接口返回字段）
type BuildingNameItem struct {
	Label string `json:"label"` // 楼栋名称（如"1栋"），用于匹配目标楼栋
	Value string `json:"value"` // 楼栋编号（全国房屋唯一编码）
	Key   string `json:"key"`   // 房源表ID（fybId），作为后续接口的请求参数
	Acive bool   `json:"acive"` // 是否激活状态（接口原始字段，拼写有误）
}

// BuildingNameResponse 楼栋名称列表接口响应结构
type BuildingNameResponse struct {
	Data   []BuildingNameItem `json:"data"`   // 楼栋列表
	Msg    string             `json:"msg"`    // 接口返回消息
	Status int                `json:"status"` // 接口返回状态码
}

// HouseInfoRequest 房源列表查询请求体结构
type HouseInfoRequest struct {
	Buildingbranch string `json:"buildingbranch"` // 楼栋单元/分区名称（如"未知"、"A座"），对应 gnqmcList 中的值
	Floor          string `json:"floor"`          // 楼层筛选条件，空字符串表示不限楼层
	FybId          string `json:"fybId"`          // 房源表ID，对应具体楼盘的唯一标识
	Housenb        string `json:"housenb"`        // 房号筛选条件，空字符串表示不限房号
	Status         string `json:"status"`         // 房源状态筛选（"0"表示可售/待售）
	Type           string `json:"type"`           // 户型类型筛选（如"三房"、"两房"）
	YsProjectId    int    `json:"ysProjectId"`    // 预售项目ID，对应楼盘项目的唯一标识
	PreSellId      int    `json:"preSellId"`      // 预售证ID，对应预售许可证的唯一标识
}

// AnalysisResult 房源数据量化分析结果
type AnalysisResult struct {
	// 一、价格类核心指标
	MinUnitPrice         float64 `json:"min_unit_price"`          // 最低挂牌单价（元/㎡）
	MaxUnitPrice         float64 `json:"max_unit_price"`          // 最高挂牌单价（元/㎡）
	MedianUnitPrice      float64 `json:"median_unit_price"`       // 挂牌单价中位数（元/㎡），所有房源单价排序后取中间值
	MinRecordedUnitPrice float64 `json:"min_recorded_unit_price"` // 最低备案单价（元/㎡），按套内面积计算
	MaxRecordedUnitPrice float64 `json:"max_recorded_unit_price"` // 最高备案单价（元/㎡），按套内面积计算
	MinTotalPrice        float64 `json:"min_total_price"`         // 最低挂牌总价（元）
	MaxTotalPrice        float64 `json:"max_total_price"`         // 最高挂牌总价（元）
	AvgTotalPrice        float64 `json:"avg_total_price"`         // 挂牌总价均值（元），所有房源总价求和后取平均
	TotalPriceSpan       float64 `json:"total_price_span"`        // 总价区间跨度（元），= 最高总价 - 最低总价

	// 二、面积类核心指标
	BuildingArea  float64 `json:"building_area"`   // 产权建筑面积（㎡），即合同面积
	ExpandArea    float64 `json:"expand_area"`     // 赠送/拓展面积（㎡），不计入产权但可实际使用
	InsideArea    float64 `json:"inside_area"`     // 套内面积（㎡），= 产权建筑面积 - 公摊面积
	ActualUseArea float64 `json:"actual_use_area"` // 实际使用面积（㎡），= 产权建筑面积 + 拓展面积
	HousingRate   float64 `json:"housing_rate"`    // 得房率（%），= 套内面积 / 产权建筑面积 × 100

	// 三、房源供应与楼层分布指标
	TotalCount     int `json:"total_count"`      // 总在售套数
	MinFloor       int `json:"min_floor"`        // 最低在售楼层
	MaxFloor       int `json:"max_floor"`        // 最高在售楼层
	LowFloorCount  int `json:"low_floor_count"`  // 低楼层（2-9楼）在售套数
	MidFloorCount  int `json:"mid_floor_count"`  // 中楼层（10-19楼）在售套数
	HighFloorCount int `json:"high_floor_count"` // 高楼层（20楼及以上）在售套数
	MaxPerFloor    int `json:"max_per_floor"`    // 单楼层最多在售套数
	MinPerFloor    int `json:"min_per_floor"`    // 单楼层最少在售套数

	// 四、性价比衍生指标
	UnitPricePerExpandArea float64 `json:"unit_price_per_expand_area"` // 单价/拓展面积比均值（元/㎡），= 挂牌单价均值 / 拓展面积，数值越低性价比越高
	CostPerActualArea      float64 `json:"cost_per_actual_area"`       // 总价/实际使用面积比均值（元/㎡），反映实际使用面积的购房成本
	FloorPricePremium      float64 `json:"floor_price_premium"`        // 楼层价格梯度（%），= (最高单价 - 最低单价) / 最低单价 × 100，反映高低楼层单价涨幅
}

// ProjectSearchRequest 楼盘搜索请求体
type ProjectSearchRequest struct {
	Project   string `json:"project"`   // 楼盘名称关键字（如"乐宸"）
	PageIndex int    `json:"pageIndex"` // 页码，从 1 开始
	PageSize  int    `json:"pageSize"`  // 每页条数
	Total     int    `json:"total"`     // 总条数（首次传 0，后续传接口返回值）
	Zone      string `json:"zone"`      // 行政区筛选，空字符串表示不限区域
}

// ProjectItem 单个楼盘信息（对应 getYsfYsPublicity 接口返回字段）
type ProjectItem struct {
	Id          string `json:"id"`              // 预售证ID（对应 preSellId）
	SypId       string `json:"sypId"`           // 预售项目ID（对应 ysProjectId）
	SypeId      string `json:"sypeId"`          // 预售证ID（与 id 相同）
	Zone        string `json:"zone"`            // 所在行政区（如"龙岗"）
	PreSellNo   string `json:"strpreprojectid"` // 预售证编号（如"深房许字（2025）龙岗020号"）
	ProjectName string `json:"project"`         // 楼盘项目名称（如"乐宸花园"）
	CompanyName string `json:"name"`            // 开发商名称
	SiteAddress string `json:"siteaddress"`     // 楼盘地址
	PassDate    string `json:"passdate"`        // 预售证批准日期
	ImagePath   string `json:"imagePath"`       // 楼盘封面图路径
}

// ProjectSearchResponse 楼盘搜索接口响应结构
type ProjectSearchResponse struct {
	Data struct {
		Total    int           `json:"total"`    // 总条数
		PageSize int           `json:"pageSize"` // 每页条数
		List     []ProjectItem `json:"list"`     // 楼盘列表
	} `json:"data"`
	Msg    string `json:"msg"`    // 接口返回消息
	Status int    `json:"status"` // 接口返回状态码
}

// ResolvedProjectParams 从搜索结果中解析出的目标楼盘参数
type ResolvedProjectParams struct {
	YsProjectId int    // 预售项目ID
	PreSellId   int    // 预售证ID
	FybId       int    // 房源表ID
	ProjectName string // 匹配到的楼盘名称
}

// SaleStatusStats 销售状态统计
type SaleStatusStats struct {
	StatusName string // 状态名称
	Count      int    // 套数
}

// SaleSummary 销售概况汇总
type SaleSummary struct {
	TotalCount    int               // 总套数
	ForSaleCount  int               // 在售套数（期房待售）
	SoldCount     int               // 已售套数（已签认购书 + 已录入合同等）
	StatusDetails []SaleStatusStats // 各状态明细
	SoldRate      float64           // 售出比例（%）
	ForSaleRate   float64           // 在售比例（%）
}
