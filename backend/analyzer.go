package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Analyze 分析房源数据，返回量化指标
func Analyze(groups []FloorGroup) AnalysisResult {
	var items []HouseItem
	for _, g := range groups {
		items = append(items, g.List...)
	}
	if len(items) == 0 {
		return AnalysisResult{}
	}

	result := AnalysisResult{}

	// 一、价格类指标（仅统计有挂牌价的在售房源）
	unitPrices := make([]float64, 0)
	recordedPrices := make([]float64, 0)
	totalPrices := make([]float64, 0)
	totalPriceSum := 0.0

	for _, item := range items {
		if item.AskPriceEachB <= 0 {
			continue
		}
		unitPrices = append(unitPrices, item.AskPriceEachB)
		totalPrices = append(totalPrices, item.AskPriceTotalB)
		totalPriceSum += item.AskPriceTotalB
		if item.RecordedPricePerUnitInside > 0 {
			recordedPrices = append(recordedPrices, item.RecordedPricePerUnitInside)
		}
	}

	if len(unitPrices) == 0 {
		return AnalysisResult{}
	}

	sort.Float64s(unitPrices)
	sort.Float64s(recordedPrices)
	sort.Float64s(totalPrices)

	result.MinUnitPrice = unitPrices[0]
	result.MaxUnitPrice = unitPrices[len(unitPrices)-1]
	if len(unitPrices)%2 == 0 {
		result.MedianUnitPrice = (unitPrices[len(unitPrices)/2-1] + unitPrices[len(unitPrices)/2]) / 2
	} else {
		result.MedianUnitPrice = unitPrices[len(unitPrices)/2]
	}
	if len(recordedPrices) > 0 {
		result.MinRecordedUnitPrice = recordedPrices[0]
		result.MaxRecordedUnitPrice = recordedPrices[len(recordedPrices)-1]
	}
	result.MinTotalPrice = totalPrices[0]
	result.MaxTotalPrice = totalPrices[len(totalPrices)-1]
	result.AvgTotalPrice = totalPriceSum / float64(len(unitPrices))
	result.TotalPriceSpan = result.MaxTotalPrice - result.MinTotalPrice

	// 二、面积类指标
	result.BuildingArea = items[0].YsbuildingArea
	result.ExpandArea = items[0].YsExpandArea
	result.InsideArea = items[0].YsInsideArea
	result.ActualUseArea = result.BuildingArea + result.ExpandArea
	if result.BuildingArea > 0 {
		result.HousingRate = math.Round(result.InsideArea/result.BuildingArea*10000) / 100
	}

	// 三、楼层分布指标
	floorCountMap := make(map[int]int)
	for _, g := range groups {
		floorNum := 0
		fmt.Sscanf(g.Floor, "%d", &floorNum)
		for _, item := range g.List {
			if item.AskPriceEachB > 0 {
				floorCountMap[floorNum]++
			}
		}
	}

	result.MinFloor = math.MaxInt32
	result.MaxFloor = 0
	result.MinPerFloor = math.MaxInt32
	result.MaxPerFloor = 0
	for floor, cnt := range floorCountMap {
		if floor < result.MinFloor {
			result.MinFloor = floor
		}
		if floor > result.MaxFloor {
			result.MaxFloor = floor
		}
		switch {
		case floor >= 2 && floor <= 9:
			result.LowFloorCount += cnt
		case floor >= 10 && floor <= 19:
			result.MidFloorCount += cnt
		case floor >= 20:
			result.HighFloorCount += cnt
		}
		if cnt > result.MaxPerFloor {
			result.MaxPerFloor = cnt
		}
		if cnt < result.MinPerFloor {
			result.MinPerFloor = cnt
		}
	}
	result.TotalCount = len(unitPrices)

	// 四、性价比衍生指标
	unitPriceSum := 0.0
	for _, p := range unitPrices {
		unitPriceSum += p
	}
	avgUnitPrice := unitPriceSum / float64(len(unitPrices))
	if result.ExpandArea > 0 {
		result.UnitPricePerExpandArea = math.Round(avgUnitPrice/result.ExpandArea*100) / 100
	}
	if result.ActualUseArea > 0 {
		result.CostPerActualArea = math.Round(result.AvgTotalPrice/result.ActualUseArea*100) / 100
	}
	if result.MinUnitPrice > 0 {
		result.FloorPricePremium = math.Round((result.MaxUnitPrice-result.MinUnitPrice)/result.MinUnitPrice*10000) / 100
	}

	return result
}

// AnalyzeSaleStatus 统计各销售状态的套数及比例
func AnalyzeSaleStatus(groups []FloorGroup) SaleSummary {
	statusCountMap := make(map[string]int)
	for _, g := range groups {
		for _, item := range g.List {
			statusCountMap[item.LastStatusName]++
		}
	}

	summary := SaleSummary{}
	for status, count := range statusCountMap {
		summary.TotalCount += count
		summary.StatusDetails = append(summary.StatusDetails, SaleStatusStats{
			StatusName: status,
			Count:      count,
		})
		if status == "期房待售" {
			summary.ForSaleCount += count
		} else {
			summary.SoldCount += count
		}
	}

	sort.Slice(summary.StatusDetails, func(i, j int) bool {
		return summary.StatusDetails[i].Count > summary.StatusDetails[j].Count
	})

	if summary.TotalCount > 0 {
		summary.SoldRate = math.Round(float64(summary.SoldCount)/float64(summary.TotalCount)*10000) / 100
		summary.ForSaleRate = math.Round(float64(summary.ForSaleCount)/float64(summary.TotalCount)*10000) / 100
	}

	return summary
}

// PrintReport 格式化打印分析报告（保留终端输出，方便调试）
func PrintReport(r AnalysisResult, s SaleSummary) {
	fmt.Println("\n========== 房源数据分析报告 ==========")
	fmt.Printf("分析时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("在售套数: %d 套 | 总套数: %d 套 | 售出比例: %.2f%%\n",
		s.ForSaleCount, s.TotalCount, s.SoldRate)
	fmt.Printf("单价区间: %.0f ~ %.0f 元/㎡ | 中位数: %.0f 元/㎡\n",
		r.MinUnitPrice, r.MaxUnitPrice, r.MedianUnitPrice)
	fmt.Printf("总价区间: %.2f ~ %.2f 万 | 均值: %.2f 万\n",
		r.MinTotalPrice/10000, r.MaxTotalPrice/10000, r.AvgTotalPrice/10000)
	fmt.Println("======================================")
}