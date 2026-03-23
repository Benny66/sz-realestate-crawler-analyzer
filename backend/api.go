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
	mux.HandleFunc("/api/search", withCORS(handleSearch))           // 搜索楼盘列表
	mux.HandleFunc("/api/buildings", withCORS(handleBuildings))     // 获取楼栋列表
	mux.HandleFunc("/api/analyze", withCORS(handleAnalyze))         // 执行分析
	mux.HandleFunc("/api/compare", withCORS(handleCompare))         // 多楼盘对比
	mux.HandleFunc("/api/history", withCORS(handleHistory))         // 历史记录（GET/DELETE）
	mux.HandleFunc("/api/cache/clear", withCORS(handleCacheClear))  // 清除缓存
}

// withCORS 跨域中间件，允许 Vue 开发服务器访问
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
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
		oldBuilding := RequestConfig.BuildingName
		oldType := RequestConfig.HouseType
		if req.BuildingName != "" {
			RequestConfig.BuildingName = req.BuildingName
		}
		if req.HouseType != "" {
			RequestConfig.HouseType = req.HouseType
		}
		params, err = ResolveProjectParams(req.Keyword)
		RequestConfig.BuildingName = oldBuilding
		RequestConfig.HouseType = oldType
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
	if houseType == "" {
		houseType = RequestConfig.HouseType
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
		Analysis:     analysis,
		Sale:         sale,
		CreatedAt:    now,
	}
	if err = SaveHistory(record); err != nil {
		fmt.Printf("保存历史记录失败（不影响主流程）: %v\n", err)
	}

	return &result, nil
}
