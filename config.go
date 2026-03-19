package main

// AppConfig 应用配置
type AppConfig struct {
	Interval        int    // 定时爬取间隔（秒）
	ProjectListURL  string // 楼盘列表搜索接口，用于通过关键字自动查找目标楼盘的 ID 参数
	BuildingNameURL string // 楼栋名称列表接口，用于通过楼栋名称匹配获取 fybId
	BuildingDictURL string // 楼栋字典接口地址，用于获取楼栋编号、可售楼层、单元分区等基础信息
	HouseInfoURL    string // 房源列表接口地址，用于获取指定楼盘的房源详情及价格数据
}

// DefaultConfig 默认配置
var DefaultConfig = AppConfig{
	Interval:        3600,
	ProjectListURL:  "https://fdc.zjj.sz.gov.cn/szfdcscjy/ysf/publicity/getYsfYsPublicity",
	BuildingNameURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingNameListToPublicity",
	BuildingDictURL: "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getBuildingDictToPublicity",
	HouseInfoURL:    "https://fdc.zjj.sz.gov.cn/szfdcscjy/projectPublish/getHouseInfoListToPublicity",
}

// RequestConfig 爬取请求参数配置（只需配置楼盘关键字、目标楼栋和户型，其余参数自动获取）
var RequestConfig = struct {
	ProjectKeyword string // 楼盘搜索关键字（如"乐宸"），用于从接口自动匹配目标楼盘
	BuildingName   string // 目标楼栋名称（如"1栋"），用于从楼栋列表中匹配 fybId
	HouseType      string // 户型类型筛选条件（如"三房"、"两房"），空字符串表示不限户型
}{
	ProjectKeyword: "盛境瑞府",
	BuildingName:   "2栋",
	HouseType:      "三房",
}

// DefaultHeaders 请求头配置（模拟真实浏览器行为，避免被服务端拦截）
var DefaultHeaders = map[string]string{
	"Accept":             "application/json, text/plain, */*",                                                                               // 接受的响应数据类型
	"Accept-Language":    "zh-CN,zh;q=0.9",                                                                                                  // 语言偏好
	"Cache-Control":      "no-cache",                                                                                                        // 禁用缓存，确保每次获取最新数据
	"Connection":         "keep-alive",                                                                                                      // 保持长连接
	"Origin":             "https://fdc.zjj.sz.gov.cn",                                                                                       // 请求来源域名
	"Pragma":             "no-cache",                                                                                                        // HTTP/1.0 兼容的禁用缓存指令
	"Referer":            "https://fdc.zjj.sz.gov.cn/szfdcscjy/",                                                                            // 请求来源页面，服务端用于校验合法性
	"Sec-Fetch-Dest":     "empty",                                                                                                           // 请求目标类型（fetch请求）
	"Sec-Fetch-Mode":     "cors",                                                                                                            // 请求模式（跨域资源共享）
	"Sec-Fetch-Site":     "same-origin",                                                                                                     // 请求与页面同源
	"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36", // 浏览器标识
	"X-Requested-With":   "XMLHttpRequest",                                                                                                  // 标识为 Ajax 异步请求
	"sec-ch-ua":          `"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`,                                               // 浏览器品牌及版本信息
	"sec-ch-ua-mobile":   "?0",                                                                                                              // 是否移动端（否）
	"sec-ch-ua-platform": `"Windows"`,                                                                                                       // 操作系统平台
}

// Cookies 会话 Cookie 配置
// 注意：Cookie 有有效期，自动获取失败时需手动从浏览器复制更新此处
// 各字段说明：
//   - WSESSIONID-SZFDC-SCJY : 服务端 Session 标识，用于维持会话状态，过期后接口将返回未授权错误
//   - _trs_uv                : 访客唯一标识，由访客追踪系统生成
//   - Hm_lvt_*               : 百度统计 Cookie，记录访问时间戳
//   - BIGipServerPool-*      : F5 负载均衡 Cookie，用于将请求路由到固定后端节点
var Cookies = "WSESSIONID-SZFDC-SCJY=hTUDywiQZTMBGFfVsA6Xt6kxyOKzlO_sbibll0kgLtc2sl5xQ3yk!1724894465; _trs_uv=mm8t1sju_850_6pv; Hm_lvt_ddaf92bcdd865fd907acdaba0285f9b1=1772509533; BIGipServerPool-fdc-SC-9001=2231723018.13091.0000; BIGipServerPool-fdc-SC-9004=2248500234.11299.0000; BIGipServerPool-fdc-SC-9002=2214945802.13347.0000"
