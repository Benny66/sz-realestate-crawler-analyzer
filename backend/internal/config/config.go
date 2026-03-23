package config

import (
	"fmt"
	"os"
	"time"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Wechat   WechatConfig   `yaml:"wechat"`
	Crawler  CrawlerConfig  `yaml:"crawler"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr string `yaml:"addr"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	HistoryFile    string `yaml:"history_file"`
	FavoritesFile  string `yaml:"favorites_file"`
	PushRecordsFile string `yaml:"push_records_file"`
}

// WechatConfig 微信配置
type WechatConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	ProjectListURL  string            `yaml:"project_list_url"`
	BuildingNameURL  string            `yaml:"building_name_url"`
	BuildingDictURL  string            `yaml:"building_dict_url"`
	HouseInfoURL    string            `yaml:"house_info_url"`
	Headers         map[string]string `yaml:"headers"`
	RequestConfig   RequestConfigType `yaml:"request_config"`
	Interval        time.Duration     `yaml:"interval"`
}

// RequestConfigType 爬取参数类型
type RequestConfigType struct {
	ProjectKeyword string `yaml:"project_keyword"`
	BuildingName   string `yaml:"building_name"`
	HouseType      string `yaml:"house_type"`
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		Database: DatabaseConfig{
			HistoryFile:    "history.json",
			FavoritesFile:  "favorites.json",
			PushRecordsFile: "push_records.json",
		},
		Wechat: WechatConfig{
			WebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f5e8ee74-efbd-4e02-9011-20ab8decb475",
		},
		Crawler: CrawlerConfig{
			ProjectListURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/ysf/publicity/getYsfYsPublicity",
			BuildingNameURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingNameListToPublicity",
			BuildingDictURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingDictToPublicity",
			HouseInfoURL:    "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getHouseInfoListToPublicity",
			Interval:        30 * time.Minute,
			Headers: map[string]string{
				"Accept":             "application/json, text/plain, */*",
				"Accept-Language":    "zh-CN,zh;q=0.9",
				"Cache-Control":      "no-cache",
				"Connection":         "keep-alive",
				"Origin":             "https://fdc.zjj.sz.gov.cn",
				"Pragma":             "no-cache",
				"Referer":            "https://fdc.zjj.sz.gov.cn/szfdcscjy/",
				"Sec-Fetch-Dest":     "empty",
				"Sec-Fetch-Mode":     "cors",
				"Sec-Fetch-Site":     "same-origin",
				"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
				"X-Requested-With":   "XMLHttpRequest",
				"sec-ch-ua":          `"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`,
				"sec-ch-ua-mobile":   "?0",
				"sec-ch-ua-platform": `"Windows"`,
			},
			RequestConfig: RequestConfigType{
				ProjectKeyword: "乐宸",
				BuildingName:   "1栋",
				HouseType:      "三房",
			},
		},
	}
}

// GetEnv 获取环境变量，如果不存在则使用默认值
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt 获取整数环境变量
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}