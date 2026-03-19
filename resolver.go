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

	// 1. 请求楼盘搜索接口
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
	// 打印原始响应，排查数据结构或关键字匹配问题
	fmt.Printf("楼盘搜索原始响应:\n%s\n", string(searchResp))
	// 2. 解析响应
	var searchResult ProjectSearchResponse
	if err = json.Unmarshal(searchResp, &searchResult); err != nil {
		return nil, fmt.Errorf("楼盘搜索结果解析失败: %v", err)
	}
	if searchResult.Status != 200 {
		return nil, fmt.Errorf("楼盘搜索接口返回异常: %s", searchResult.Msg)
	}

	// 3. 从结果中精确匹配楼盘名称
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

			// 4. 继续获取 fybId
			fybId, err := ResolveFybId(ysProjectId, preSellId, RequestConfig.BuildingName)
			if err != nil {
				return nil, fmt.Errorf("获取 fybId 失败: %v", err)
			}
			params.FybId = fybId

			return params, nil
		}
	}

	return nil, fmt.Errorf("未找到包含关键字「%s」的楼盘，请检查关键字是否正确", keyword)
}

// ResolveFybId 通过楼栋名称匹配获取对应的 fybId
func ResolveFybId(ysProjectId, preSellId int, buildingName string) (int, error) {
	fmt.Printf("正在获取楼栋列表，匹配目标楼栋: 「%s」...\n", buildingName)

	// 1. 请求楼栋名称列表接口（form 表单格式）
	resp, err := postForm(DefaultConfig.BuildingNameURL, map[string]string{
		"ysProjectId": strconv.Itoa(ysProjectId),
		"preSellId":   strconv.Itoa(preSellId),
	})
	if err != nil {
		return 0, fmt.Errorf("获取楼栋列表失败: %v", err)
	}
	// 打印原始响应，排查数据结构或关键字匹配问题
	fmt.Printf("楼栋列表原始响应:\n%s\n", string(resp))

	// 2. 解析响应
	var buildingResp BuildingNameResponse
	if err = json.Unmarshal(resp, &buildingResp); err != nil {
		return 0, fmt.Errorf("楼栋列表解析失败: %v", err)
	}
	if buildingResp.Status != 200 {
		return 0, fmt.Errorf("楼栋列表接口返回异常: %s", buildingResp.Msg)
	}

	// 3. 匹配目标楼栋名称，提取 key 作为 fybId
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

	return 0, fmt.Errorf("未找到楼栋「%s」，请检查 config.go 中 BuildingName 配置是否正确", buildingName)
}
