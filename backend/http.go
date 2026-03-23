package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// postForm 发送 POST 请求（form 表单格式）
func postForm(url string, data map[string]string) ([]byte, error) {
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

// postJSON 发送 POST 请求（JSON 格式）
func postJSON(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %v", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return doRequest(req)
}

// setCommonHeaders 设置公共请求头
func setCommonHeaders(req *http.Request) {
	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}
	cookieStr, err := GetCookieString()
	if err != nil {
		cookieStr = Cookies
	}
	req.Header.Set("Cookie", cookieStr)
}

// doRequest 执行 HTTP 请求并返回响应体
func doRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}
	return body, nil
}