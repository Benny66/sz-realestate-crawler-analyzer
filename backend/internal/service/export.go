package service

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"sz-realestate-crawler-analyzer/internal/model"
)

// ExportService 导出服务
type ExportService struct{}

// NewExportService 创建导出服务
func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportHousesCSV 导出房源 CSV
func (s *ExportService) ExportHousesCSV(w http.ResponseWriter, resp *model.AnalyzeResponse) error {
	filename := fmt.Sprintf("houses_%s_%d_%s.csv",
		resp.Params.ProjectName,
		resp.Params.PreSellId,
		time.Now().Format("20060102150405"),
	)

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(http.StatusOK)

	// 写入 UTF-8 BOM
	_, _ = w.Write([]byte("\xEF\xBB\xBF"))
	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{
		"楼盘", "房号", "楼层", "楼栋", "状态", "建筑面积(㎡)", "套内面积(㎡)", "拓展面积(㎡)",
		"挂牌单价(元/㎡)", "挂牌总价(元)", "备案单价(元/㎡)", "用途", "无障碍",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, h := range resp.Houses {
		record := []string{
			resp.Params.ProjectName,
			h.Housenb,
			h.Floor,
			h.BuildingName,
			h.LastStatusName,
			fmt.Sprintf("%.2f", h.YsbuildingArea),
			fmt.Sprintf("%.2f", h.YsInsideArea),
			fmt.Sprintf("%.2f", h.YsExpandArea),
			fmt.Sprintf("%.2f", h.AskPriceEachB),
			fmt.Sprintf("%.2f", h.AskPriceTotalB),
			fmt.Sprintf("%.2f", h.RecordedPricePerUnitInside),
			h.UseAge,
			h.BarrierFree,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// ExportAnalysisCSV 导出分析结果 CSV
func (s *ExportService) ExportAnalysisCSV(w http.ResponseWriter, resp *model.AnalyzeResponse) error {
	filename := fmt.Sprintf("analysis_%s_%d_%s.csv",
		resp.Params.ProjectName,
		resp.Params.PreSellId,
		time.Now().Format("20060102150405"),
	)

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(http.StatusOK)

	// 写入 UTF-8 BOM
	_, _ = w.Write([]byte("\xEF\xBB\xBF"))
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 分析结果头部
	headers := []string{
		"楼盘名称", "楼栋名称", "户型", "分析时间",
		"在售套数", "总套数", "售出比例",
		"最低单价(元/㎡)", "最高单价(元/㎡)", "中位单价(元/㎡)",
		"最低总价(元)", "最高总价(元)", "平均总价(元)",
		"建筑面积(㎡)", "套内面积(㎡)", "拓展面积(㎡)", "实际使用面积(㎡)",
		"低楼层套数", "中楼层套数", "高楼层套数",
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// 分析结果数据
	record := []string{
		resp.Params.ProjectName,
		resp.Params.BuildingName,
		resp.Params.HouseType,
		resp.UpdatedAt,
		fmt.Sprintf("%d", resp.Sale.ForSaleCount),
		fmt.Sprintf("%d", resp.Sale.TotalCount),
		fmt.Sprintf("%.2f%%", resp.Sale.SoldRate),
		fmt.Sprintf("%.0f", resp.Analysis.MinUnitPrice),
		fmt.Sprintf("%.0f", resp.Analysis.MaxUnitPrice),
		fmt.Sprintf("%.0f", resp.Analysis.MedianUnitPrice),
		fmt.Sprintf("%.0f", resp.Analysis.MinTotalPrice),
		fmt.Sprintf("%.0f", resp.Analysis.MaxTotalPrice),
		fmt.Sprintf("%.0f", resp.Analysis.AvgTotalPrice),
		fmt.Sprintf("%.2f", resp.Analysis.BuildingArea),
		fmt.Sprintf("%.2f", resp.Analysis.InsideArea),
		fmt.Sprintf("%.2f", resp.Analysis.ExpandArea),
		fmt.Sprintf("%.2f", resp.Analysis.ActualUseArea),
		fmt.Sprintf("%d", resp.Analysis.LowFloorCount),
		fmt.Sprintf("%d", resp.Analysis.MidFloorCount),
		fmt.Sprintf("%d", resp.Analysis.HighFloorCount),
	}
	if err := writer.Write(record); err != nil {
		return err
	}

	// 销售状态详情
	writer.Write([]string{}) // 空行
	writer.Write([]string{"销售状态详情"})
	writer.Write([]string{"状态名称", "套数"})

	for _, status := range resp.Sale.StatusDetails {
		statusRecord := []string{
			status.StatusName,
			fmt.Sprintf("%d", status.Count),
		}
		if err := writer.Write(statusRecord); err != nil {
			return err
		}
	}

	return nil
}
