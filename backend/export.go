package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// exportHousesCSV 导出房源 CSV
func exportHousesCSV(w http.ResponseWriter, resp *AnalyzeResponse) error {
	filename := fmt.Sprintf("houses_%s_%s_%s.csv",
		resp.Params.ProjectName,
		resp.Params.PreSellId,
		time.Now().Format("20060102150405"),
	)

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(http.StatusOK)

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

// parseIntDefault 安全转换 int
func parseIntDefault(v string, d int) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		return d
	}
	return n
}
