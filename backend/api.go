package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/search", withCORS(handleSearch))                   // 搜索楼盘列表
	mux.HandleFunc("/api/buildings", withCORS(handleBuildings))             // 获取楼栋列表
	mux.HandleFunc("/api/analyze", withCORS(handleAnalyze))                 // 执行分析
	mux.HandleFunc("/api/compare", withCORS(handleCompare))                 // 多楼盘对比
	mux.HandleFunc("/api/history", withCORS(handleHistory))                 // 历史记录（GET/DELETE）
	mux.HandleFunc("/api/favorites", withCORS(handleFavorites))             // 收藏楼盘（GET/POST/DELETE）
	mux.HandleFunc("/api/favorites/update", withCORS(handleUpdateFavorite)) // 更新收藏（PUT）
	mux.HandleFunc("/api/export/csv", withCORS(handleExportCSV))            // 导出 CSV
	mux.HandleFunc("/api/cache/clear", withCORS(handleCacheClear))          // 清除缓存
}

// withCORS 跨域中间件，允许 Vue 开发服务器访问
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

// writeJSON 写入 JSON 响应
func writeJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// writeOK 写入成功响应
func writeOK(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, APIResponse{Code: 0, Message: "success", Data: data})
}

// writeErr 写入错误响应
func writeErr(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusOK, APIResponse{Code: 1, Message: msg, Data: nil})
}

// ---- Handler 实现 ----

// handleSearch GET /api/search?keyword=乐宸&pageIndex=1&pageSize=12&zone=
func handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, "仅支持 GET 请求")
		return
	}
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		writeErr(w, "keyword 参数不能为空")
		return
	}
	pageIndex, _ := strconv.Atoi(r.URL.Query().Get("pageIndex"))
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize <= 0 {
		pageSize = 12
	}
	zone := r.URL.Query().Get("zone")

	resp, err := postJSON(DefaultConfig.ProjectListURL, ProjectSearchRequest{
		Project:   keyword,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Total:     0,
		Zone:      zone,
	})
	if err != nil {
		writeErr(w, fmt.Sprintf("搜索楼盘失败: %v", err))
		return
	}

	var searchResult ProjectSearchResponse
	if err = json.Unmarshal(resp, &searchResult); err != nil {
		writeErr(w, fmt.Sprintf("解析搜索结果失败: %v", err))
		return
	}
	if searchResult.Status != 200 {
		writeErr(w, fmt.Sprintf("接口返回异常: %s", searchResult.Msg))
		return
	}

	writeOK(w, ProjectListResponse{
		Total:    searchResult.Data.Total,
		Projects: searchResult.Data.List,
	})
}

// handleBuildings GET /api/buildings?ysProjectId=xxx&preSellId=xxx
func handleBuildings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, "仅支持 GET 请求")
		return
	}
	ysProjectId, err := strconv.Atoi(r.URL.Query().Get("ysProjectId"))
	if err != nil || ysProjectId == 0 {
		writeErr(w, "ysProjectId 参数无效")
		return
	}
	preSellId, err := strconv.Atoi(r.URL.Query().Get("preSellId"))
	if err != nil || preSellId == 0 {
		writeErr(w, "preSellId 参数无效")
		return
	}

	resp, err := postForm(DefaultConfig.BuildingNameURL, map[string]string{
		"ysProjectId": strconv.Itoa(ysProjectId),
		"preSellId":   strconv.Itoa(preSellId),
	})
	if err != nil {
		writeErr(w, fmt.Sprintf("获取楼栋列表失败: %v", err))
		return
	}

	var buildingResp BuildingNameResponse
	if err = json.Unmarshal(resp, &buildingResp); err != nil {
		writeErr(w, fmt.Sprintf("解析楼栋列表失败: %v", err))
		return
	}
	if buildingResp.Status != 200 {
		writeErr(w, fmt.Sprintf("楼栋接口返回异常: %s", buildingResp.Msg))
		return
	}

	writeOK(w, BuildingListResponse{Buildings: buildingResp.Data})
}

// handleAnalyze POST /api/analyze
func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, "仅支持 POST 请求")
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, fmt.Sprintf("请求体解析失败: %v", err))
		return
	}
	if req.Keyword == "" && req.YsProjectId == 0 {
		writeErr(w, "keyword 或 ysProjectId 不能同时为空")
		return
	}

	result, err := doAnalyze(req)
	if err != nil {
		writeErr(w, err.Error())
		return
	}
	writeOK(w, result)
}

// handleCompare POST /api/compare
func handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, "仅支持 POST 请求")
		return
	}

	var req CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, fmt.Sprintf("请求体解析失败: %v", err))
		return
	}
	if len(req.Items) == 0 {
		writeErr(w, "对比列表不能为空")
		return
	}
	if len(req.Items) > 5 {
		writeErr(w, "最多支持 5 个楼盘同时对比")
		return
	}

	compareItems := make([]CompareItem, 0, len(req.Items))
	for _, item := range req.Items {
		result, err := doAnalyze(item)
		if err != nil {
			writeErr(w, fmt.Sprintf("分析「%s %s」失败: %v", item.Keyword, item.BuildingName, err))
			return
		}
		compareItems = append(compareItems, CompareItem{
			ProjectName:  result.Params.ProjectName,
			BuildingName: item.BuildingName,
			HouseType:    item.HouseType,
			Zone:         item.Zone,
			Analysis:     result.Analysis,
			Sale:         result.Sale,
		})
	}

	writeOK(w, CompareResponse{
		Items:     compareItems,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
}

// handleHistory GET/DELETE /api/history
func handleHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /api/history?projectName=xxx  可按楼盘名过滤
		records, err := LoadHistory()
		if err != nil {
			writeErr(w, fmt.Sprintf("加载历史记录失败: %v", err))
			return
		}
		// 可选过滤
		filterProject := r.URL.Query().Get("projectName")
		if filterProject != "" {
			filtered := make([]HistoryRecord, 0)
			for _, rec := range records {
				if strings.Contains(rec.ProjectName, filterProject) {
					filtered = append(filtered, rec)
				}
			}
			records = filtered
		}
		writeOK(w, records)

	case http.MethodDelete:
		// DELETE /api/history?id=xxx
		id := r.URL.Query().Get("id")
		if id == "" {
			writeErr(w, "id 参数不能为空")
			return
		}
		if err := DeleteHistory(id); err != nil {
			writeErr(w, fmt.Sprintf("删除历史记录失败: %v", err))
			return
		}
		writeOK(w, nil)

	default:
		writeErr(w, "仅支持 GET / DELETE 请求")
	}
}

// handleFavorites GET/POST/DELETE /api/favorites
func handleFavorites(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := LoadFavorites()
		if err != nil {
			writeErr(w, fmt.Sprintf("加载收藏失败: %v", err))
			return
		}
		writeOK(w, items)
	case http.MethodPost:
		var item FavoriteItem
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			writeErr(w, fmt.Sprintf("收藏请求解析失败: %v", err))
			return
		}
		if item.ProjectName == "" {
			writeErr(w, "projectName 不能为空")
			return
		}
		if item.ID == "" {
			item.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
		}
		if item.CreatedAt == "" {
			item.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		}
		// 默认开启推送
		if !item.EnablePush {
			item.EnablePush = true
		}
		if !item.PriceAlert {
			item.PriceAlert = true
		}
		if !item.SaleAlert {
			item.SaleAlert = true
		}

		if err := AddFavorite(item); err != nil {
			writeErr(w, fmt.Sprintf("保存收藏失败: %v", err))
			return
		}

		// 发送新收藏通知
		go func() {
			message := BuildNewFavoriteMessage(item.ProjectName, item.BuildingName, item.HouseType, item.Zone)
			if err := SendWechatMarkdownMessage(message); err != nil {
				fmt.Printf("发送新收藏通知失败: %v\n", err)
			} else {
				fmt.Printf("已发送新收藏通知: %s\n", item.ProjectName)
			}
		}()

		writeOK(w, nil)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeErr(w, "id 参数不能为空")
			return
		}
		if err := DeleteFavorite(id); err != nil {
			writeErr(w, fmt.Sprintf("删除收藏失败: %v", err))
			return
		}
		writeOK(w, nil)
	default:
		writeErr(w, "仅支持 GET / POST / DELETE 请求")
	}
}

// handleExportCSV GET /api/export/csv?keyword=xxx&buildingName=xxx&houseType=xxx&zone=
func handleExportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErr(w, "仅支持 GET 请求")
		return
	}
	req := AnalyzeRequest{
		Keyword:      r.URL.Query().Get("keyword"),
		BuildingName: r.URL.Query().Get("buildingName"),
		HouseType:    r.URL.Query().Get("houseType"),
		Zone:         r.URL.Query().Get("zone"),
		YsProjectId:  parseIntDefault(r.URL.Query().Get("ysProjectId"), 0),
		PreSellId:    parseIntDefault(r.URL.Query().Get("preSellId"), 0),
		FybId:        parseIntDefault(r.URL.Query().Get("fybId"), 0),
	}
	if req.Keyword == "" && req.YsProjectId == 0 {
		writeErr(w, "keyword 或 ysProjectId 不能为空")
		return
	}
	result, err := doAnalyze(req)
	if err != nil {
		writeErr(w, err.Error())
		return
	}
	if err := exportHousesCSV(w, result); err != nil {
		writeErr(w, fmt.Sprintf("导出 CSV 失败: %v", err))
	}
}

// handleUpdateFavorite PUT /api/favorites/:id
func handleUpdateFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErr(w, "仅支持 PUT 请求")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeErr(w, "id 参数不能为空")
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		writeErr(w, fmt.Sprintf("更新请求解析失败: %v", err))
		return
	}

	favorites, err := LoadFavorites()
	if err != nil {
		writeErr(w, fmt.Sprintf("加载收藏失败: %v", err))
		return
	}

	updated := false
	for i, fav := range favorites {
		if fav.ID == id {
			if enablePush, ok := updateData["enablePush"].(bool); ok {
				favorites[i].EnablePush = enablePush
			}
			if priceAlert, ok := updateData["priceAlert"].(bool); ok {
				favorites[i].PriceAlert = priceAlert
			}
			if saleAlert, ok := updateData["saleAlert"].(bool); ok {
				favorites[i].SaleAlert = saleAlert
			}
			updated = true
			break
		}
	}

	if !updated {
		writeErr(w, "未找到对应的收藏记录")
		return
	}

	if err := SaveFavorites(favorites); err != nil {
		writeErr(w, fmt.Sprintf("更新收藏失败: %v", err))
		return
	}

	writeOK(w, nil)
}

// handleCacheClear POST /api/cache/clear
func handleCacheClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, "仅支持 POST 请求")
		return
	}
	analysisCacheMu.Lock()
	analysisCache = make(map[string]*CacheEntry)
	analysisCacheMu.Unlock()
	writeOK(w, nil)
}

// ---- 核心分析逻辑（供 handleAnalyze / handleCompare 复用）----

// doAnalyze 执行完整的爬取+分析流程，带缓存
func doAnalyze(req AnalyzeRequest) (*AnalyzeResponse, error) {
	var params *ResolvedProjectParams
	var err error

	// 1. 解析楼盘参数
	if req.FybId != 0 && req.YsProjectId != 0 && req.PreSellId != 0 {
		// 前端直接传入了完整 ID，跳过搜索步骤
		params = &ResolvedProjectParams{
			YsProjectId: req.YsProjectId,
			PreSellId:   req.PreSellId,
			FybId:       req.FybId,
			ProjectName: req.Keyword,
		}
	} else {
		// 通过关键字自动解析
		// 创建临时配置，避免修改全局默认值
		tempConfig := RequestConfig
		if req.BuildingName != "" {
			tempConfig.BuildingName = req.BuildingName
		}
		// 只有当明确指定户型时才使用，否则分析全部户型
		if req.HouseType != "" {
			tempConfig.HouseType = req.HouseType
		} else {
			tempConfig.HouseType = "" // 空字符串表示不限户型
		}

		// 保存原始配置并临时替换
		oldConfig := RequestConfig
		RequestConfig = tempConfig
		params, err = ResolveProjectParamsWithZoneAndConfig(req.Keyword, req.Zone, tempConfig)
		// 恢复原始配置
		RequestConfig = oldConfig

		if err != nil {
			return nil, fmt.Errorf("解析楼盘参数失败: %v", err)
		}
	}

	// 2. 检查缓存（30 分钟有效）
	buildingName := req.BuildingName
	if buildingName == "" {
		buildingName = RequestConfig.BuildingName
	}
	houseType := req.HouseType
	// 如果前端没有指定户型，就分析全部户型（空字符串）
	// 不要使用默认值，因为默认值可能是特定户型
	if houseType == "" {
		houseType = "" // 空字符串表示不限户型
	}

	if entry, ok := GetCache(params.ProjectName, buildingName, houseType); ok {
		if IsCacheValid(entry, 30*time.Minute) {
			fmt.Printf("[Cache HIT] %s %s %s\n", params.ProjectName, buildingName, houseType)
			resp := entry.Response
			return &resp, nil
		}
	}

	// 3. 爬取房源列表
	houseInfoResp, err := postJSON(DefaultConfig.HouseInfoURL, HouseInfoRequest{
		Buildingbranch: "",
		Floor:          "",
		FybId:          fmt.Sprintf("%d", params.FybId),
		Housenb:        "",
		Status:         "-1",
		Type:           houseType,
		YsProjectId:    params.YsProjectId,
		PreSellId:      params.PreSellId,
	})
	if err != nil {
		return nil, fmt.Errorf("爬取房源列表失败: %v", err)
	}

	var houseInfoResult HouseInfoResponse
	if err = json.Unmarshal(houseInfoResp, &houseInfoResult); err != nil {
		return nil, fmt.Errorf("房源列表解析失败: %v", err)
	}

	// 4. 汇总原始房源列表
	allHouses := make([]HouseItem, 0)
	for _, g := range houseInfoResult.Data {
		allHouses = append(allHouses, g.List...)
	}

	// 5. 分析
	analysis := Analyze(houseInfoResult.Data)
	sale := AnalyzeSaleStatus(houseInfoResult.Data)

	now := time.Now().Format("2006-01-02 15:04:05")
	result := AnalyzeResponse{
		Analysis:    analysis,
		Sale:        sale,
		Houses:      allHouses,
		FloorGroups: houseInfoResult.Data,
		Params:      *params,
		Alerts:      buildAlerts(params.ProjectName, buildingName, houseType, analysis, sale),
		UpdatedAt:   now,
	}

	// 6. 写入缓存
	SetCache(params.ProjectName, buildingName, houseType, result)

	// 7. 持久化历史记录
	record := HistoryRecord{
		ID:           fmt.Sprintf("%d", time.Now().UnixMilli()),
		ProjectName:  params.ProjectName,
		BuildingName: buildingName,
		HouseType:    houseType,
		Zone:         req.Zone,
		Analysis:     analysis,
		Sale:         sale,
		CreatedAt:    now,
	}
	if err = SaveHistory(record); err != nil {
		fmt.Printf("保存历史记录失败（不影响主流程）: %v\n", err)
	}

	return &result, nil
}

// buildAlerts 根据最近历史生成简单预警
func buildAlerts(projectName, buildingName, houseType string, analysis AnalysisResult, sale SaleSummary) []PriceAlert {
	records, err := LoadHistory()
	if err != nil {
		return nil
	}
	filtered := make([]HistoryRecord, 0)
	for _, r := range records {
		if r.ProjectName == projectName && r.BuildingName == buildingName && r.HouseType == houseType {
			filtered = append(filtered, r)
		}
	}
	if len(filtered) == 0 {
		return []PriceAlert{{
			Level:     "info",
			Title:     "首次记录",
			Message:   "当前楼盘暂无可对比历史记录，已建立首条分析基线。",
			Direction: "flat",
			ChangePct: 0,
		}}
	}
	last := filtered[len(filtered)-1]
	alerts := make([]PriceAlert, 0)
	if last.Analysis.MedianUnitPrice > 0 {
		change := (analysis.MedianUnitPrice - last.Analysis.MedianUnitPrice) / last.Analysis.MedianUnitPrice * 100
		if change >= 3 {
			alerts = append(alerts, PriceAlert{Level: "danger", Title: "价格上涨预警", Message: fmt.Sprintf("中位单价较上次上涨 %.2f%%", change), Direction: "up", ChangePct: change})
		} else if change <= -3 {
			alerts = append(alerts, PriceAlert{Level: "success", Title: "价格回落提醒", Message: fmt.Sprintf("中位单价较上次下降 %.2f%%", -change), Direction: "down", ChangePct: change})
		}
	}
	if last.Sale.ForSaleRate > 0 {
		rateChange := sale.ForSaleRate - last.Sale.ForSaleRate
		if rateChange <= -5 {
			alerts = append(alerts, PriceAlert{Level: "warning", Title: "去化加速提醒", Message: fmt.Sprintf("在售比例较上次下降 %.2f%%，说明去化速度可能加快", -rateChange), Direction: "down", ChangePct: rateChange})
		}
	}
	if len(alerts) == 0 {
		alerts = append(alerts, PriceAlert{Level: "info", Title: "走势平稳", Message: "相比上次记录，价格与去化变化较小。", Direction: "flat", ChangePct: 0})
	}
	return alerts
}
