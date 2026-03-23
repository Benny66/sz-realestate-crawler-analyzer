package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ResolveProjectParams 通过关键字自动搜索并解析目标楼盘的 ID 参数
func ResolveProjectParams(keyword string) (*ResolvedProjectParams, error) {
	fmt.Printf("正在搜索楼盘关键字: 「%s」...\n", keyword)

	searchResp, err := postJSON(DefaultConfig.ProjectListURL, ProjectSearchRequest{
		Project:   keyword,
		PageIndex: 1,
		PageSize:  12,
		Total:     0,
		Zone:      "",
	})
	if err != nil {
		return nil, fmt.Errorf("搜索楼盘失败: %v", err)
	}

	var searchResult ProjectSearchResponse
	if err = json.Unmarshal(searchResp, &searchResult); err != nil {
		return nil, fmt.Errorf("楼盘搜索结果解析失败: %v", err)
	}
	if searchResult.Status != 200 {
		return nil, fmt.Errorf("楼盘搜索接口返回异常: %s", searchResult.Msg)
	}

	for _, project := range searchResult.Data.List {
		if strings.Contains(project.ProjectName, keyword) {
			ysProjectId, err := strconv.Atoi(project.SypId)
			if err != nil {
				return nil, fmt.Errorf("ysProjectId 转换失败: %v", err)
			}
			preSellId, err := strconv.Atoi(project.Id)
			if err != nil {
				return nil, fmt.Errorf("preSellId 转换失败: %v", err)
			}

			params := &ResolvedProjectParams{
				YsProjectId: ysProjectId,
				PreSellId:   preSellId,
				ProjectName: project.ProjectName,
			}
			fmt.Printf("匹配到楼盘: 「%s」 ysProjectId=%d preSellId=%d\n",
				params.ProjectName, params.YsProjectId, params.PreSellId)

			fybId, err := ResolveFybId(ysProjectId, preSellId, RequestConfig.BuildingName)
			if err != nil {
				return nil, fmt.Errorf("获取 fybId 失败: %v", err)
			}
			params.FybId = fybId
			return params, nil
		}
	}

	return nil, fmt.Errorf("未找到包含关键字「%s」的楼盘", keyword)
}

// ResolveFybId 通过楼栋名称匹配获取对应的 fybId
func ResolveFybId(ysProjectId, preSellId int, buildingName string) (int, error) {
	fmt.Printf("正在获取楼栋列表，匹配目标楼栋: 「%s」...\n", buildingName)

	resp, err := postForm(DefaultConfig.BuildingNameURL, map[string]string{
		"ysProjectId": strconv.Itoa(ysProjectId),
		"preSellId":   strconv.Itoa(preSellId),
	})
	if err != nil {
		return 0, fmt.Errorf("获取楼栋列表失败: %v", err)
	}

	var buildingResp BuildingNameResponse
	if err = json.Unmarshal(resp, &buildingResp); err != nil {
		return 0, fmt.Errorf("楼栋列表解析失败: %v", err)
	}
	if buildingResp.Status != 200 {
		return 0, fmt.Errorf("楼栋列表接口返回异常: %s", buildingResp.Msg)
	}

	for _, building := range buildingResp.Data {
		if building.Label == buildingName {
			fybId, err := strconv.Atoi(building.Key)
			if err != nil {
				return 0, fmt.Errorf("fybId 转换失败: %v", err)
			}
			fmt.Printf("匹配到楼栋: 「%s」 fybId=%d\n", building.Label, fybId)
			return fybId, nil
		}
	}

	return 0, fmt.Errorf("未找到楼栋「%s」", buildingName)
}