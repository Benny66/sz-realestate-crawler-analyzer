package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// postForm 发送POST请求（form表单格式）
func postForm(url string, data map[string]string) ([]byte, error) {
	// 构建form表单body
	parts := make([]string, 0, len(data))
	for k, v := range data {
		parts = append(parts, k+"="+v)
	}
	formBody := strings.Join(parts, "&")

	req, err := http.NewRequest("POST", url, strings.NewReader(formBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return doRequest(req)
}

// postJSON 发送POST请求（JSON格式）
func postJSON(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 打印请求参数
	// fmt.Printf("POST %s\n请求参数: %s\n", url, string(jsonData))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	return doRequest(req)
}

// setCommonHeaders 设置公共请求头（Cookie 自动获取）
func setCommonHeaders(req *http.Request) {
	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}
	// 优先使用自动获取的 Cookie，失败则降级到静态配置
	cookieStr, err := GetCookieString()
	if err != nil {
		cookieStr = Cookies // 降级到 config.go 中的静态 Cookie
	}
	req.Header.Set("Cookie", cookieStr)
}

// doRequest 执行HTTP请求并返回响应体
func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应内容: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
