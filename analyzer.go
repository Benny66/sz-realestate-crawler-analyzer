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
		fmt.Println("无房源数据可分析")
		return AnalysisResult{}
	}

	result := AnalysisResult{}

	// ---- 一、价格类指标（仅统计有挂牌价的在售房源）----
	unitPrices := make([]float64, 0)
	recordedPrices := make([]float64, 0)
	totalPrices := make([]float64, 0)
	totalPriceSum := 0.0

	for _, item := range items {
		// 过滤掉挂牌价为 0 的房源（已售/认购状态价格字段为 null）
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
		fmt.Println("无有效挂牌价数据")
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

	// ---- 二、面积类指标 ----
	result.BuildingArea = items[0].YsbuildingArea
	result.ExpandArea = items[0].YsExpandArea
	result.InsideArea = items[0].YsInsideArea
	result.ActualUseArea = result.BuildingArea + result.ExpandArea
	if result.BuildingArea > 0 {
		result.HousingRate = math.Round(result.InsideArea/result.BuildingArea*10000) / 100
	}

	// ---- 三、楼层分布指标（仅统计在售房源）----
	floorCountMap := make(map[int]int)
	for _, g := range groups {
		floorNum := 0
		fmt.Sscanf(g.Floor, "%d", &floorNum)
		for _, item := range g.List {
			// 仅统计挂牌价大于 0 的在售房源
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

	// TotalCount 也改为仅统计在售套数，与价格统计保持一致
	result.TotalCount = len(unitPrices)

	// ---- 四、性价比衍生指标 ----
	avgUnitPrice := totalPriceSum / float64(len(unitPrices)) // 复用已有求和
	// 注意：这里应用单价均值，重新计算
	unitPriceSum := 0.0
	for _, p := range unitPrices {
		unitPriceSum += p
	}
	avgUnitPrice = unitPriceSum / float64(len(unitPrices))

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

	// 遍历所有房源，按 lastStatusName 分组统计
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
		// 期房待售 = 在售；其余状态均视为已售
		if status == "期房待售" {
			summary.ForSaleCount += count
		} else {
			summary.SoldCount += count
		}
	}

	// 按套数降序排列状态明细
	sort.Slice(summary.StatusDetails, func(i, j int) bool {
		return summary.StatusDetails[i].Count > summary.StatusDetails[j].Count
	})

	if summary.TotalCount > 0 {
		summary.SoldRate = math.Round(float64(summary.SoldCount)/float64(summary.TotalCount)*10000) / 100
		summary.ForSaleRate = math.Round(float64(summary.ForSaleCount)/float64(summary.TotalCount)*10000) / 100
	}

	return summary
}

// PrintReport 格式化打印分析报告
func PrintReport(r AnalysisResult, s SaleSummary) {
	fmt.Println("\n========== 房源数据分析报告 ==========")
	fmt.Printf("分析时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("\n【一、价格类核心指标】")
	fmt.Println("  ▶ 单价相关")
	fmt.Printf("    最低挂牌单价:   %.2f 元/㎡\n", r.MinUnitPrice)
	fmt.Printf("    最高挂牌单价:   %.2f 元/㎡\n", r.MaxUnitPrice)
	fmt.Printf("    挂牌单价中位数: %.2f 元/㎡\n", r.MedianUnitPrice)
	fmt.Printf("    备案单价区间:   %.2f ~ %.2f 元/㎡\n", r.MinRecordedUnitPrice, r.MaxRecordedUnitPrice)
	fmt.Println("  ▶ 总价相关")
	fmt.Printf("    最低总价:     %.0f 元（约 %.2f 万）\n", r.MinTotalPrice, r.MinTotalPrice/10000)
	fmt.Printf("    最高总价:     %.0f 元（约 %.2f 万）\n", r.MaxTotalPrice, r.MaxTotalPrice/10000)
	fmt.Printf("    总价均值:     约 %.2f 万\n", r.AvgTotalPrice/10000)
	fmt.Printf("    总价区间跨度: %.2f 万\n", r.TotalPriceSpan/10000)

	fmt.Println("\n【二、面积类核心指标】")
	fmt.Printf("    产权建筑面积:   %.2f ㎡\n", r.BuildingArea)
	fmt.Printf("    赠送/拓展面积:  %.2f ㎡\n", r.ExpandArea)
	fmt.Printf("    套内面积:       %.2f ㎡\n", r.InsideArea)
	fmt.Printf("    实际使用面积:   %.2f ㎡（建筑面积 + 拓展面积）\n", r.ActualUseArea)
	fmt.Printf("    得房率:         %.2f%%\n", r.HousingRate)

	fmt.Println("\n【三、房源供应与楼层分布】")
	fmt.Printf("    总在售套数:       %d 套\n", r.TotalCount)
	fmt.Printf("    在售楼层范围:     %d ~ %d 楼\n", r.MinFloor, r.MaxFloor)
	fmt.Printf("    低楼层(2-9楼):   %d 套（占比 %.1f%%）\n", r.LowFloorCount, float64(r.LowFloorCount)/float64(r.TotalCount)*100)
	fmt.Printf("    中楼层(10-19楼): %d 套（占比 %.1f%%）\n", r.MidFloorCount, float64(r.MidFloorCount)/float64(r.TotalCount)*100)
	fmt.Printf("    高楼层(20+楼):   %d 套（占比 %.1f%%）\n", r.HighFloorCount, float64(r.HighFloorCount)/float64(r.TotalCount)*100)
	fmt.Printf("    单楼层最多套数:   %d 套\n", r.MaxPerFloor)
	fmt.Printf("    单楼层最少套数:   %d 套\n", r.MinPerFloor)

	fmt.Println("\n【四、性价比衍生指标】")
	fmt.Printf("    单价/拓展面积比均值:      %.2f 元/㎡\n", r.UnitPricePerExpandArea)
	fmt.Printf("    总价/实际使用面积比均值:  %.2f 元/㎡\n", r.CostPerActualArea)
	fmt.Printf("    楼层价格梯度（高低涨幅）: %.2f%%\n", r.FloorPricePremium)

	fmt.Println("\n【五、销售状态统计】")
	fmt.Printf("    总套数:   %d 套\n", s.TotalCount)
	fmt.Printf("    在售套数: %d 套（占比 %.2f%%）\n", s.ForSaleCount, s.ForSaleRate)
	fmt.Printf("    已售套数: %d 套（占比 %.2f%%）\n", s.SoldCount, s.SoldRate)
	fmt.Println("  ▶ 各状态明细")
	for _, detail := range s.StatusDetails {
		rate := math.Round(float64(detail.Count)/float64(s.TotalCount)*10000) / 100
		fmt.Printf("    %-12s %d 套（占比 %.2f%%）\n", detail.StatusName, detail.Count, rate)
	}

	fmt.Println("\n======================================")
}
