package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"sz-realestate-crawler-analyzer/internal/config"
	"sz-realestate-crawler-analyzer/internal/model"
	"sz-realestate-crawler-analyzer/pkg/http"
)

// CrawlerService 爬虫服务
type CrawlerService struct {
	config *config.Config
	client *http.Client
}

// NewCrawlerService 创建爬虫服务
func NewCrawlerService(cfg *config.Config) *CrawlerService {
	client := http.NewClient("", cfg.Crawler.Headers)
	return &CrawlerService{
		config: cfg,
		client: client,
	}
}

// SearchProjects 搜索楼盘
func (s *CrawlerService) SearchProjects(keyword, zone string, pageIndex, pageSize int) (*model.ProjectListResponse, error) {
	searchReq := map[string]interface{}{
		"project":   keyword,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
		"total":     0,
		"zone":      zone,
	}
	
	resp, err := s.client.PostJSON(s.config.Crawler.ProjectListURL, searchReq)
	if err != nil {
		return nil, fmt.Errorf("搜索楼盘失败: %v", err)
	}
	
	var searchResult struct {
		Data struct {
			Total    int               `json:"total"`
			List     []model.ProjectItem `json:"list"`
		} `json:"data"`
		Msg    string `json:"msg"`
		Status int    `json:"status"`
	}
	
	if err = json.Unmarshal(resp, &searchResult); err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %v", err)
	}
	
	if searchResult.Status != 200 {
		return nil, fmt.Errorf("接口返回异常: %s", searchResult.Msg)
	}
	
	return &model.ProjectListResponse{
		Total:    searchResult.Data.Total,
		Projects: searchResult.Data.List,
	}, nil
}

// GetBuildings 获取楼栋列表
func (s *CrawlerService) GetBuildings(ysProjectId, preSellId int) (*model.BuildingListResponse, error) {
	resp, err := s.client.PostForm(s.config.Crawler.BuildingNameURL, map[string]string{
		"ysProjectId": strconv.Itoa(ysProjectId),
		"preSellId":   strconv.Itoa(preSellId),
	})
	if err != nil {
		return nil, fmt.Errorf("获取楼栋列表失败: %v", err)
	}
	
	var buildingResp struct {
		Data   []model.BuildingNameItem `json:"data"`
		Msg    string                   `json:"msg"`
		Status int                      `json:"status"`
	}
	
	if err = json.Unmarshal(resp, &buildingResp); err != nil {
		return nil, fmt.Errorf("解析楼栋列表失败: %v", err)
	}
	
	if buildingResp.Status != 200 {
		return nil, fmt.Errorf("楼栋接口返回异常: %s", buildingResp.Msg)
	}
	
	return &model.BuildingListResponse{
		Buildings: buildingResp.Data,
	}, nil
}

// FetchHouses 获取房源列表
func (s *CrawlerService) FetchHouses(params *model.ResolvedProjectParams, houseType string) ([]model.FloorGroup, error) {
	houseReq := map[string]interface{}{
		"buildingbranch": "",
		"floor":          "",
		"fybId":          fmt.Sprintf("%d", params.FybId),
		"housenb":        "",
		"status":         "-1",
		"type":           houseType,
		"ysProjectId":    params.YsProjectId,
		"preSellId":      params.PreSellId,
	}
	
	resp, err := s.client.PostJSON(s.config.Crawler.HouseInfoURL, houseReq)
	if err != nil {
		return nil, fmt.Errorf("爬取房源列表失败: %v", err)
	}
	
	var houseInfoResult struct {
		Data   []model.FloorGroup `json:"data"`
		Msg    string             `json:"msg"`
		Status int                `json:"status"`
	}
	
	if err = json.Unmarshal(resp, &houseInfoResult); err != nil {
		return nil, fmt.Errorf("房源列表解析失败: %v", err)
	}
	
	if houseInfoResult.Status != 200 {
		return nil, fmt.Errorf("房源接口返回异常: %s", houseInfoResult.Msg)
	}
	
	return houseInfoResult.Data, nil
}

// ResolveProjectParams 解析楼盘参数
func (s *CrawlerService) ResolveProjectParams(keyword, zone, buildingName string) (*model.ResolvedProjectParams, error) {
	fmt.Printf("正在搜索楼盘关键字: 「%s」 区域: 「%s」...\n", keyword, zone)
	
	searchResult, err := s.SearchProjects(keyword, zone, 1, 12)
	if err != nil {
		return nil, err
	}
	
	for _, project := range searchResult.Projects {
		if strings.Contains(project.ProjectName, keyword) {
			ysProjectId, err := strconv.Atoi(project.SypId)
			if err != nil {
				return nil, fmt.Errorf("ysProjectId 转换失败: %v", err)
			}
			preSellId, err := strconv.Atoi(project.Id)
			if err != nil {
				return nil, fmt.Errorf("preSellId 转换失败: %v", err)
			}
			
			params := &model.ResolvedProjectParams{
				YsProjectId: ysProjectId,
				PreSellId:   preSellId,
				ProjectName: project.ProjectName,
			}
			
			fmt.Printf("匹配到楼盘: 「%s」 ysProjectId=%d preSellId=%d\n",
				params.ProjectName, params.YsProjectId, params.PreSellId)
			
			// 解析楼栋ID
			fybId, actualBuildingName, err := s.ResolveFybId(ysProjectId, preSellId, buildingName)
			if err != nil {
				return nil, fmt.Errorf("获取 fybId 失败: %v", err)
			}
			
			params.FybId = fybId
			params.BuildingName = actualBuildingName
			params.AutoSelected = (buildingName == "" && actualBuildingName != "")
			
			return params, nil
		}
	}
	
	if zone != "" {
		return nil, fmt.Errorf("未找到包含关键字「%s」且区域为「%s」的楼盘", keyword, zone)
	}
	return nil, fmt.Errorf("未找到包含关键字「%s」的楼盘", keyword)
}

// ResolveFybId 解析楼栋ID
func (s *CrawlerService) ResolveFybId(ysProjectId, preSellId int, buildingName string) (int, string, error) {
	fmt.Printf("正在获取楼栋列表，匹配目标楼栋: 「%s」...\n", buildingName)
	
	buildingResult, err := s.GetBuildings(ysProjectId, preSellId)
	if err != nil {
		return 0, "", err
	}
	
	// 如果没有楼栋数据
	if len(buildingResult.Buildings) == 0 {
		return 0, "", fmt.Errorf("该楼盘没有可用的楼栋数据")
	}
	
	// 如果 buildingName 为空，使用第一个楼栋
	if buildingName == "" {
		building := buildingResult.Buildings[0]
		fybId, err := strconv.Atoi(building.Key)
		if err != nil {
			return 0, "", fmt.Errorf("fybId 转换失败: %v", err)
		}
		fmt.Printf("自动选择第一个楼栋: 「%s」 fybId=%d\n", building.Label, fybId)
		return fybId, building.Label, nil
	}
	
	// 按名称匹配楼栋
	for _, building := range buildingResult.Buildings {
		if building.Label == buildingName {
			fybId, err := strconv.Atoi(building.Key)
			if err != nil {
				return 0, "", fmt.Errorf("fybId 转换失败: %v", err)
			}
			fmt.Printf("匹配到楼栋: 「%s」 fybId=%d\n", building.Label, fybId)
			return fybId, building.Label, nil
		}
	}
	
	// 如果指定的楼栋不存在，返回错误和可用的楼栋列表
	availableBuildings := make([]string, len(buildingResult.Buildings))
	for i, b := range buildingResult.Buildings {
		availableBuildings[i] = b.Label
	}
	return 0, "", fmt.Errorf("未找到楼栋「%s」，可用楼栋: %v", buildingName, availableBuildings)
}