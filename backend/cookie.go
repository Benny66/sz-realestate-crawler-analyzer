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
	cookieExpireAt   time.Time
)

const (
	cookieRefreshInterval = 30 * time.Minute
	cookieInitURL         = "https://fdc.zjj.sz.gov.cn/szfdcscjy/"
)

func initCookieClient() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("创建 CookieJar 失败: %v", err)
	}
	client := &http.Client{Jar: jar, Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", cookieInitURL, nil)
	if err != nil {
		return fmt.Errorf("创建初始化请求失败: %v", err)
	}
	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("访问入口页失败: %v", err)
	}
	defer resp.Body.Close()

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
	return nil
}

func ensureCookieValid() error {
	cookieMu.RLock()
	expired := time.Now().After(cookieExpireAt)
	cookieMu.RUnlock()
	if expired || globalHTTPClient == nil {
		return initCookieClient()
	}
	return nil
}

func GetHTTPClient() (*http.Client, error) {
	if err := ensureCookieValid(); err != nil {
		return nil, err
	}
	cookieMu.RLock()
	defer cookieMu.RUnlock()
	return globalHTTPClient, nil
}

func GetCookieString() (string, error) {
	if err := ensureCookieValid(); err != nil {
		return Cookies, nil
	}
	cookieMu.RLock()
	defer cookieMu.RUnlock()

	req, _ := http.NewRequest("GET", cookieInitURL, nil)
	cookies := globalCookieJar.Cookies(req.URL)
	parts := make([]string, 0, len(cookies))
	for _, c := range cookies {
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; "), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}