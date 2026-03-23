package main

// AppConfig 应用配置
type AppConfig struct {
	ServerAddr      string // HTTP 服务监听地址（如 ":8080"）
	Interval        int    // 定时爬取间隔（秒）
	ProjectListURL  string // 楼盘列表搜索接口
	BuildingNameURL string // 楼栋名称列表接口
	BuildingDictURL string // 楼栋字典接口
	HouseInfoURL    string // 房源列表接口
}

// DefaultConfig 默认配置
var DefaultConfig = AppConfig{
	ServerAddr:      ":8080",
	Interval:        3600,
	ProjectListURL:  "https://fdc.zjj.sz.gov.cn/szfdcscjy/ysf/publicity/getYsfYsPublicity",
	BuildingNameURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingNameListToPublicity",
	BuildingDictURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingDictToPublicity",
	HouseInfoURL:    "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getHouseInfoListToPublicity",
}

// RequestConfig 默认爬取参数（定时任务使用，API 请求时会动态覆盖）
var RequestConfig = struct {
	ProjectKeyword string
	BuildingName   string
	HouseType      string
}{
	ProjectKeyword: "乐宸",
	BuildingName:   "1栋",
	HouseType:      "三房",
}

// DefaultHeaders 请求头配置（模拟真实浏览器）
var DefaultHeaders = map[string]string{
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
}

// Cookies 静态备用 Cookie（自动获取失败时降级使用）
var Cookies = "WSESSIONID-SZFDC-SCJY=hTUDywiQZTMBGFfVsA6Xt6kxyOKzlO_sbibll0kgLtc2sl5xQ3yk!1724894465; _trs_uv=mm8t1sju_850_6pv; Hm_lvt_ddaf92bcdd865fd907acdaba0285f9b1=1772509533; BIGipServerPool-fdc-SC-9001=2231723018.13091.0000; BIGipServerPool-fdc-SC-9004=2248500234.11299.0000; BIGipServerPool-fdc-SC-9002=2214945802.13347.0000"
