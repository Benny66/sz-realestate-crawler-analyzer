package service

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/internal/repository"
)

// AnalyzerService 分析服务
type AnalyzerService struct {
	config      *config.Config
	crawler     *CrawlerService
	cache       *CacheService
	wechat      *WechatService
	export      *ExportService
	historyRepo *repository.HistoryRepository
}

// NewAnalyzerService 创建分析服务
func NewAnalyzerService(cfg *config.Config) *AnalyzerService {
	crawler := NewCrawlerService(cfg)
	cache := NewCacheService()
	wechat := NewWechatService(cfg)
	export := NewExportService()
	historyRepo := repository.NewHistoryRepository(cfg)

	return &AnalyzerService{
		config:      cfg,
		crawler:     crawler,
		cache:       cache,
		wechat:      wechat,
		export:      export,
		historyRepo: historyRepo,
	}
}

// Analyze 执行房源数据分析（核心算法）
func (s *AnalyzerService) Analyze(groups []model.FloorGroup) model.AnalysisResult {
	var result model.AnalysisResult

	// 汇总所有房源
	allHouses := make([]model.HouseItem, 0)
	floorMap := make(map[string]int) // 楼层房源数量统计

	for _, group := range groups {
		allHouses = append(allHouses, group.List...)
		floorMap[group.Floor] = len(group.List)
	}

	if len(allHouses) == 0 {
		return result
	}

	// 价格相关计算（仅统计有挂牌价的在售房源，与旧模块保持一致）
	unitPrices := make([]float64, 0)
	recordedPrices := make([]float64, 0)
	totalPrices := make([]float64, 0)
	totalPriceSum := 0.0

	// 筛选有挂牌价的在售房源
	onSaleHouses := make([]model.HouseItem, 0)
	for _, house := range allHouses {
		if house.AskPriceEachB > 0 {
			onSaleHouses = append(onSaleHouses, house)
			unitPrices = append(unitPrices, house.AskPriceEachB)
			totalPrices = append(totalPrices, house.AskPriceTotalB)
			totalPriceSum += house.AskPriceTotalB

			if house.RecordedPricePerUnitInside > 0 {
				recordedPrices = append(recordedPrices, house.RecordedPricePerUnitInside)
			}
		}
	}

	if len(unitPrices) == 0 {
		return result
	}

	// 排序价格数组用于计算中位数
	sort.Float64s(unitPrices)
	sort.Float64s(totalPrices)
	sort.Float64s(recordedPrices)

	// 计算统计指标
	result.MinUnitPrice = s.min(unitPrices)
	result.MaxUnitPrice = s.max(unitPrices)
	result.MedianUnitPrice = s.median(unitPrices)

	result.MinRecordedUnitPrice = s.min(recordedPrices)
	result.MaxRecordedUnitPrice = s.max(recordedPrices)

	result.MinTotalPrice = s.min(totalPrices)
	result.MaxTotalPrice = s.max(totalPrices)
	result.AvgTotalPrice = s.average(totalPrices)
	result.TotalPriceSpan = result.MaxTotalPrice - result.MinTotalPrice

	// 面积相关（使用单个房源面积，不是总和）
	if len(allHouses) > 0 {
		// 使用第一个房源的面积作为基准（旧模块逻辑）
		firstHouse := allHouses[0]
		result.BuildingArea = firstHouse.YsbuildingArea
		result.ExpandArea = firstHouse.YsExpandArea
		result.InsideArea = firstHouse.YsInsideArea
		result.ActualUseArea = firstHouse.YsInsideArea + firstHouse.YsExpandArea

		if firstHouse.YsbuildingArea > 0 {
			result.HousingRate = (firstHouse.YsInsideArea / firstHouse.YsbuildingArea) * 100
		}
	}

	// 楼层统计（只统计在售房源）
	result.TotalCount = len(onSaleHouses)
	result.MinFloor, result.MaxFloor = s.calculateFloorRangeForOnSale(groups, onSaleHouses)

	// 楼层分布统计（只统计在售房源）
	low, mid, high := s.calculateFloorDistributionForOnSale(groups, result.MinFloor, result.MaxFloor, onSaleHouses)
	result.LowFloorCount = low
	result.MidFloorCount = mid
	result.HighFloorCount = high

	// 每层房源数量统计（只统计在售房源）
	onSaleFloorMap := s.buildOnSaleFloorMap(groups, onSaleHouses)
	result.MaxPerFloor, result.MinPerFloor = s.calculateFloorDensity(onSaleFloorMap)

	// 衍生指标
	if result.ExpandArea > 0 {
		result.UnitPricePerExpandArea = result.MedianUnitPrice / result.ExpandArea
	}
	if result.ActualUseArea > 0 {
		result.CostPerActualArea = result.AvgTotalPrice / result.ActualUseArea
	}

	// 楼层溢价计算
	result.FloorPricePremium = s.calculateFloorPremium(groups)

	return result
}

// AnalyzeSaleStatus 分析销售状态
func (s *AnalyzerService) AnalyzeSaleStatus(groups []model.FloorGroup) model.SaleSummary {
	var summary model.SaleSummary
	statusCounts := make(map[string]int)

	// 统计所有房源
	for _, group := range groups {
		for _, house := range group.List {
			summary.TotalCount++
			statusCounts[house.LastStatusName]++

			// 判断是否在售（这里需要根据实际状态判断）
			if strings.Contains(house.LastStatusName, "在售") || strings.Contains(house.LastStatusName, "可售") {
				summary.ForSaleCount++
			} else if strings.Contains(house.LastStatusName, "已售") || strings.Contains(house.LastStatusName, "已备案") {
				summary.SoldCount++
			}
		}
	}

	// 计算比例
	if summary.TotalCount > 0 {
		summary.SoldRate = (float64(summary.SoldCount) / float64(summary.TotalCount)) * 100
		summary.ForSaleRate = (float64(summary.ForSaleCount) / float64(summary.TotalCount)) * 100
	}

	// 状态详情
	for status, count := range statusCounts {
		summary.StatusDetails = append(summary.StatusDetails, model.SaleStatusStats{
			StatusName: status,
			Count:      count,
		})
	}

	return summary
}

// 工具函数
func (s *AnalyzerService) min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return values[0]
}

func (s *AnalyzerService) max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return values[len(values)-1]
}

func (s *AnalyzerService) median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	mid := len(values) / 2
	if len(values)%2 == 0 {
		return (values[mid-1] + values[mid]) / 2
	}
	return values[mid]
}

func (s *AnalyzerService) average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (s *AnalyzerService) calculateFloorRange(groups []model.FloorGroup) (min, max int) {
	minFloor, maxFloor := 100, -1

	for _, group := range groups {
		floor, err := strconv.Atoi(group.Floor)
		if err == nil {
			if floor < minFloor {
				minFloor = floor
			}
			if floor > maxFloor {
				maxFloor = floor
			}
		}
	}

	if minFloor == 100 && maxFloor == -1 {
		return 0, 0
	}
	return minFloor, maxFloor
}

func (s *AnalyzerService) calculateFloorDistribution(groups []model.FloorGroup, minFloor, maxFloor int) (low, mid, high int) {
	if maxFloor <= minFloor {
		return 0, 0, 0
	}

	// 动态楼层分布计算（更合理，适应不同高度的楼盘）
	totalFloors := maxFloor - minFloor + 1
	lowThreshold := minFloor + totalFloors/3
	highThreshold := minFloor + totalFloors*2/3

	for _, group := range groups {
		floor, err := strconv.Atoi(group.Floor)
		if err == nil {
			if floor <= lowThreshold {
				low += len(group.List)
			} else if floor <= highThreshold {
				mid += len(group.List)
			} else {
				high += len(group.List)
			}
		}
	}

	return low, mid, high
}

func (s *AnalyzerService) calculateFloorDensity(floorMap map[string]int) (max, min int) {
	if len(floorMap) == 0 {
		return 0, 0
	}

	maxDensity, minDensity := 0, 1000
	for _, count := range floorMap {
		if count > maxDensity {
			maxDensity = count
		}
		if count < minDensity {
			minDensity = count
		}
	}

	return maxDensity, minDensity
}

// calculateFloorRangeForOnSale 计算在售房源的楼层范围
func (s *AnalyzerService) calculateFloorRangeForOnSale(groups []model.FloorGroup, onSaleHouses []model.HouseItem) (min, max int) {
	minFloor, maxFloor := 100, -1

	// 创建在售房源的楼层映射
	onSaleFloorMap := make(map[int]bool)
	for _, house := range onSaleHouses {
		floor, err := strconv.Atoi(house.Floor)
		if err == nil {
			onSaleFloorMap[floor] = true
		}
	}

	// 只统计有在售房源的楼层
	for floor := range onSaleFloorMap {
		if floor < minFloor {
			minFloor = floor
		}
		if floor > maxFloor {
			maxFloor = floor
		}
	}

	if minFloor == 100 && maxFloor == -1 {
		return 0, 0
	}
	return minFloor, maxFloor
}

// calculateFloorDistributionForOnSale 计算在售房源的楼层分布
func (s *AnalyzerService) calculateFloorDistributionForOnSale(groups []model.FloorGroup, minFloor, maxFloor int, onSaleHouses []model.HouseItem) (low, mid, high int) {
	if maxFloor <= minFloor {
		return 0, 0, 0
	}

	// 动态楼层分布计算（更合理，适应不同高度的楼盘）
	totalFloors := maxFloor - minFloor + 1
	lowThreshold := minFloor + totalFloors/3
	highThreshold := minFloor + totalFloors*2/3

	// 统计在售房源的楼层分布
	for _, house := range onSaleHouses {
		floor, err := strconv.Atoi(house.Floor)
		if err == nil {
			if floor <= lowThreshold {
				low++
			} else if floor <= highThreshold {
				mid++
			} else {
				high++
			}
		}
	}

	return low, mid, high
}

// buildOnSaleFloorMap 构建在售房源的楼层数量统计
func (s *AnalyzerService) buildOnSaleFloorMap(groups []model.FloorGroup, onSaleHouses []model.HouseItem) map[string]int {
	floorMap := make(map[string]int)

	for _, house := range onSaleHouses {
		floorMap[house.Floor]++
	}

	return floorMap
}

func (s *AnalyzerService) calculateFloorPremium(groups []model.FloorGroup) float64 {
	if len(groups) == 0 {
		return 0.0
	}

	// 计算不同楼层的平均价格
	floorPrices := make(map[int][]float64)

	for _, group := range groups {
		floor, err := strconv.Atoi(group.Floor)
		if err != nil {
			continue
		}

		for _, house := range group.List {
			if house.AskPriceEachB > 0 {
				floorPrices[floor] = append(floorPrices[floor], house.AskPriceEachB)
			}
		}
	}

	// 计算楼层溢价
	if len(floorPrices) < 2 {
		return 0.0
	}

	// 计算底层和高层的平均价格差异
	var lowFloorPrices, highFloorPrices []float64

	for floor, prices := range floorPrices {
		if floor <= 10 {
			lowFloorPrices = append(lowFloorPrices, prices...)
		} else {
			highFloorPrices = append(highFloorPrices, prices...)
		}
	}

	if len(lowFloorPrices) == 0 || len(highFloorPrices) == 0 {
		return 0.0
	}

	lowAvg := s.average(lowFloorPrices)
	highAvg := s.average(highFloorPrices)

	if lowAvg > 0 {
		return ((highAvg - lowAvg) / lowAvg) * 100
	}

	return 0.0
}

// DoAnalyze 执行完整的分析流程
func (s *AnalyzerService) DoAnalyze(req model.AnalyzeRequest) (*model.AnalyzeResponse, error) {
	// 检查缓存
	cacheKey := CacheKey(req.Keyword, req.BuildingName, req.HouseType)
	if cached := s.cache.Get(cacheKey); cached != nil {
		return cached, nil
	}

	// 解析楼盘参数
	params, err := s.crawler.ResolveProjectParams(req.Keyword, req.Zone, req.BuildingName)
	if err != nil {
		return nil, fmt.Errorf("解析楼盘参数失败: %v", err)
	}

	// 爬取房源数据
	houses, err := s.crawler.FetchHouses(params, req.HouseType)
	if err != nil {
		return nil, fmt.Errorf("爬取房源数据失败: %v", err)
	}

	// 分析数据
	analysis := s.Analyze(houses)
	sale := s.AnalyzeSaleStatus(houses)

	// 构建响应
	result := &model.AnalyzeResponse{
		Analysis:    analysis,
		Sale:        sale,
		Houses:      s.flattenHouses(houses),
		FloorGroups: houses,
		Params:      *params,
		Alerts:      s.buildAlerts(params.ProjectName, req.BuildingName, req.HouseType, analysis, sale),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	// 缓存结果
	s.cache.Set(cacheKey, result)

	// 保存历史记录
	record := &model.HistoryRecord{
		ID:           fmt.Sprintf("%d", time.Now().UnixMilli()),
		ProjectName:  params.ProjectName,
		BuildingName: req.BuildingName,
		HouseType:    req.HouseType,
		Zone:         req.Zone,
		Analysis:     analysis,
		Sale:         sale,
		CreatedAt:    result.UpdatedAt,
	}

	if err := s.historyRepo.Save(record); err != nil {
		fmt.Printf("保存历史记录失败（不影响主流程）: %v\n", err)
	}

	return result, nil
}

// SearchProjects 搜索楼盘
func (s *AnalyzerService) SearchProjects(keyword, zone string, pageIndex, pageSize int) (*model.ProjectListResponse, error) {
	return s.crawler.SearchProjects(keyword, zone, pageIndex, pageSize)
}

// GetBuildings 获取楼栋列表
func (s *AnalyzerService) GetBuildings(ysProjectId, preSellId int) (*model.BuildingListResponse, error) {
	return s.crawler.GetBuildings(ysProjectId, preSellId)
}

// ExportHousesCSV 导出房源 CSV
func (s *AnalyzerService) ExportHousesCSV(w http.ResponseWriter, resp *model.AnalyzeResponse) error {
	return s.export.ExportHousesCSV(w, resp)
}

// ExportAnalysisCSV 导出分析结果 CSV
func (s *AnalyzerService) ExportAnalysisCSV(w http.ResponseWriter, resp *model.AnalyzeResponse) error {
	return s.export.ExportAnalysisCSV(w, resp)
}

// ClearCache 清除缓存
func (s *AnalyzerService) ClearCache() {
	s.cache.Clear()
	fmt.Println("[缓存] 缓存已清除")
}

// StartBackgroundTasks 启动后台任务
func (s *AnalyzerService) StartBackgroundTasks() {
	// 启动微信推送监控
	s.wechat.StartPushMonitor()

	// 启动定时刷新任务
	go s.scheduleAutoRefresh()
}

func (s *AnalyzerService) flattenHouses(groups []model.FloorGroup) []model.HouseItem {
	var houses []model.HouseItem
	for _, group := range groups {
		houses = append(houses, group.List...)
	}
	return houses
}

func (s *AnalyzerService) buildAlerts(projectName, buildingName, houseType string, analysis model.AnalysisResult, sale model.SaleSummary) []model.PriceAlert {
	// 从历史记录中获取对比数据
	records, err := s.historyRepo.FindByProject(projectName)
	if err != nil || len(records) == 0 {
		return []model.PriceAlert{{
			Level:     "info",
			Title:     "首次记录",
			Message:   "当前楼盘暂无可对比历史记录，已建立首条分析基线。",
			Direction: "flat",
			ChangePct: 0,
		}}
	}

	// 获取最近的历史记录
	last := records[len(records)-1]
	alerts := make([]model.PriceAlert, 0)

	// 价格变化预警
	if last.Analysis.MedianUnitPrice > 0 && analysis.MedianUnitPrice > 0 {
		change := (analysis.MedianUnitPrice - last.Analysis.MedianUnitPrice) / last.Analysis.MedianUnitPrice * 100
		if change >= 3 {
			alerts = append(alerts, model.PriceAlert{
				Level:     "danger",
				Title:     "价格上涨预警",
				Message:   fmt.Sprintf("中位单价较上次上涨 %.2f%%", change),
				Direction: "up",
				ChangePct: change,
			})
		} else if change <= -3 {
			alerts = append(alerts, model.PriceAlert{
				Level:     "success",
				Title:     "价格回落提醒",
				Message:   fmt.Sprintf("中位单价较上次下降 %.2f%%", -change),
				Direction: "down",
				ChangePct: change,
			})
		}
	}

	// 销售状态预警
	if last.Sale.ForSaleRate > 0 && sale.ForSaleRate > 0 {
		rateChange := sale.ForSaleRate - last.Sale.ForSaleRate
		if rateChange <= -5 {
			alerts = append(alerts, model.PriceAlert{
				Level:     "warning",
				Title:     "去化加速提醒",
				Message:   fmt.Sprintf("在售比例较上次下降 %.2f%%，说明去化速度可能加快", -rateChange),
				Direction: "down",
				ChangePct: rateChange,
			})
		}
	}

	if len(alerts) == 0 {
		alerts = append(alerts, model.PriceAlert{
			Level:     "info",
			Title:     "走势平稳",
			Message:   "相比上次记录，价格与去化变化较小。",
			Direction: "flat",
			ChangePct: 0,
		})
	}

	return alerts
}

// scheduleAutoRefresh 定时自动刷新
func (s *AnalyzerService) scheduleAutoRefresh() {
	fmt.Println("[定时任务] 预热默认楼盘缓存...")

	// 预热默认楼盘
	req := model.AnalyzeRequest{
		Keyword:      s.config.Crawler.RequestConfig.ProjectKeyword,
		BuildingName: s.config.Crawler.RequestConfig.BuildingName,
		HouseType:    s.config.Crawler.RequestConfig.HouseType,
	}

	_, err := s.DoAnalyze(req)
	if err != nil {
		fmt.Printf("[定时任务] 预热失败（不影响服务启动）: %v\n", err)
	} else {
		fmt.Println("[定时任务] 预热成功")
	}

	// 定时刷新
	ticker := time.NewTicker(s.config.Crawler.Interval)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("[定时任务] 开始刷新缓存 - %s\n", time.Now().Format("2006-01-02 15:04:05"))
		_, err := s.DoAnalyze(req)
		if err != nil {
			fmt.Printf("[定时任务] 刷新失败: %v\n", err)
		} else {
			fmt.Println("[定时任务] 刷新成功")
		}
	}
}
