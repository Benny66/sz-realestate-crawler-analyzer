package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	globalHTTPClient *http.Client
	cookieMu         sync.RWMutex
	cookieExpireAt   time.Time
)

// Client HTTP客户端封装
type Client struct {
	client  *http.Client
	headers map[string]string
	baseURL string
}

// NewClient 创建新的HTTP客户端
func NewClient(baseURL string, headers map[string]string) *Client {
	if globalHTTPClient == nil {
		globalHTTPClient = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	return &Client{
		client:  globalHTTPClient,
		headers: headers,
		baseURL: baseURL,
	}
}

// ensureCookieValid 确保Cookie有效
func ensureCookieValid() error {
	cookieMu.RLock()
	expired := time.Now().After(cookieExpireAt)
	cookieMu.RUnlock()
	if expired || globalHTTPClient == nil {
		return initCookieClient()
	}
	return nil
}

// initCookieClient 初始化Cookie客户端
func initCookieClient() error {
	cookieMu.Lock()
	defer cookieMu.Unlock()

	// 这里可以添加自动获取Cookie的逻辑
	// 目前先设置一个未来的过期时间
	cookieExpireAt = time.Now().Add(30 * time.Minute)

	return nil
}

// PostJSON 发送JSON POST请求
func (c *Client) PostJSON(url string, data interface{}) ([]byte, error) {
	if err := ensureCookieValid(); err != nil {
		return nil, fmt.Errorf("cookie验证失败: %v", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// PostForm 发送表单POST请求
func (c *Client) PostForm(url string, data map[string]string) ([]byte, error) {
	if err := ensureCookieValid(); err != nil {
		return nil, fmt.Errorf("cookie验证失败: %v", err)
	}

	// 构建表单数据
	formData := ""
	for k, v := range data {
		if formData != "" {
			formData += "&"
		}
		formData += fmt.Sprintf("%s=%s", k, v)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// Get 发送GET请求
func (c *Client) Get(url string, params map[string]string) ([]byte, error) {
	if err := ensureCookieValid(); err != nil {
		return nil, fmt.Errorf("cookie验证失败: %v", err)
	}

	// 构建查询参数
	if len(params) > 0 {
		query := ""
		for k, v := range params {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", k, v)
		}
		url = url + "?" + query
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	c.setHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// setHeaders 设置请求头
func (c *Client) setHeaders(req *http.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

// cacheKey 生成缓存键
func cacheKey(projectName, buildingName, houseType string) string {
	return fmt.Sprintf("%s_%s_%s", projectName, buildingName, houseType)
}
