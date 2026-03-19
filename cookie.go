package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

var (
	globalCookieJar  *cookiejar.Jar
	globalHTTPClient *http.Client
	cookieMu         sync.RWMutex
	cookieExpireAt   time.Time // Cookie 过期时间（保守估计）
)

const (
	// Cookie 刷新间隔，保守设置为 30 分钟
	cookieRefreshInterval = 30 * time.Minute
	// 初始化 Cookie 的入口页面
	cookieInitURL = "https://fdc.zjj.sz.gov.cn/szfdcscjy/"
)

// initCookieClient 初始化带 CookieJar 的 HTTP 客户端，并访问入口页获取 Cookie
func initCookieClient() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("创建 CookieJar 失败: %v", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 15 * time.Second,
	}

	// 访问入口页，触发服务端下发 Cookie
	req, err := http.NewRequest("GET", cookieInitURL, nil)
	if err != nil {
		return fmt.Errorf("创建初始化请求失败: %v", err)
	}
	// 使用配置中的请求头，模拟真实浏览器
	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("访问入口页失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查是否拿到了关键 Cookie
	cookies := jar.Cookies(req.URL)
	if len(cookies) == 0 {
		return fmt.Errorf("未能从入口页获取到任何 Cookie")
	}

	cookieMu.Lock()
	globalCookieJar = jar
	globalHTTPClient = client
	cookieExpireAt = time.Now().Add(cookieRefreshInterval)
	cookieMu.Unlock()

	fmt.Printf("Cookie 自动获取成功，共 %d 个，有效期至: %s\n",
		len(cookies), cookieExpireAt.Format("2006-01-02 15:04:05"))
	printCookies(cookies)

	return nil
}

// ensureCookieValid 确保 Cookie 有效，过期则自动刷新
func ensureCookieValid() error {
	cookieMu.RLock()
	expired := time.Now().After(cookieExpireAt)
	cookieMu.RUnlock()

	if expired || globalHTTPClient == nil {
		fmt.Println("Cookie 已过期或未初始化，正在自动刷新...")
		return initCookieClient()
	}
	return nil
}

// GetHTTPClient 获取带有效 Cookie 的 HTTP 客户端
func GetHTTPClient() (*http.Client, error) {
	if err := ensureCookieValid(); err != nil {
		return nil, err
	}
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return globalHTTPClient, nil
}

// GetCookieString 获取当前 Cookie 字符串（用于手动设置 Header 的场景）
func GetCookieString() (string, error) {
	if err := ensureCookieValid(); err != nil {
		// 自动获取失败时，降级使用 config.go 中的静态 Cookie
		fmt.Printf("自动获取 Cookie 失败，降级使用静态 Cookie: %v\n", err)
		return Cookies, nil
	}

	cookieMu.RLock()
	defer cookieMu.RUnlock()

	// 解析入口页 URL 获取对应的 Cookie
	req, _ := http.NewRequest("GET", cookieInitURL, nil)
	cookies := globalCookieJar.Cookies(req.URL)

	parts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; "), nil
}

// printCookies 打印获取到的 Cookie 列表（调试用）
func printCookies(cookies []*http.Cookie) {
	fmt.Println("  获取到的 Cookie:")
	for _, c := range cookies {
		expires := "Session"
		if !c.Expires.IsZero() {
			expires = c.Expires.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("    %-40s 过期: %s\n", c.Name+"="+c.Value[:min(len(c.Value), 20)]+"...", expires)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
