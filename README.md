# 深圳房产信息爬取分析工具

## 项目结构

```
sz-realestate-crawler-analyzer/
├── main.go       # 程序入口、爬取主流程、定时调度
├── config.go     # 所有配置（URL、请求头、Cookie、业务参数）
├── types.go      # 所有结构体定义
├── http.go       # HTTP 请求工具函数
├── cookie.go     # Cookie 自动获取与刷新逻辑
├── analyzer.go   # 数据分析逻辑与报告打印
├── resolver.go   # 楼盘参数自动解析逻辑
└── README.md     # 项目说明文档
```

## 快速启动

### 环境要求
- Go 1.21 及以上版本

### 1. 克隆项目
```bash
git clone https://github.com/your-username/sz-realestate-crawler-analyzer.git
cd sz-realestate-crawler-analyzer
```

### 2. 修改配置（可选）
打开 `config.go`，按需修改以下参数：
```go
var RequestConfig = struct {
    ProjectKeyword string // 楼盘搜索关键字，默认"乐宸"
    BuildingName   string // 目标楼栋名称，默认"1栋"
    HouseType      string // 户型类型筛选，默认"三房"
}{
    ProjectKeyword: "乐宸",
    BuildingName:   "1栋",
    HouseType:      "三房",
}
```

### 3. 运行项目
```bash
go run .
```

### 4. 查看输出
程序启动后会自动执行以下流程：
1. 自动获取 Cookie（失败则降级使用 `config.go` 中的静态 Cookie）
2. 通过关键字搜索匹配目标楼盘，自动获取楼盘 ID 参数
3. 爬取房源数据并在终端打印分析报告

输出示例：
```
深圳房产信息爬取脚本启动！
爬取间隔: 1800 秒
==============================
正在自动获取 Cookie...
Cookie 自动获取成功，共 5 个

开始爬取房源信息 - 2025-01-01 10:00:00
正在搜索楼盘关键字: 「乐宸」...
匹配到楼盘: 「乐宸花园」 ysProjectId=36386 preSellId=141990
匹配到楼栋: 「1栋」 fybId=55674
成功获取 20 个楼层的房源数据

========== 房源数据分析报告 ==========
【一、价格类核心指标】
  ...
【五、销售状态统计】
  ...
======================================
```

### 5. Cookie 过期处理
若程序输出接口返回异常或数据为空，说明 Cookie 已过期，需手动更新 `config.go` 中的 `Cookies` 字段：
1. 打开浏览器访问 `https://fdc.zjj.sz.gov.cn/szfdcscjy/`
2. 按 F12 打开开发者工具 → Network → 任意请求 → 复制 `Cookie` 请求头的值
3. 替换 `config.go` 中的 `Cookies` 变量值

---

## 启动流程

```
启动 main.go
 └─ initCookieClient()                      # 自动获取 Cookie
     └─ GET https://fdc.zjj.sz.gov.cn/szfdcscjy/
         └─ 服务端下发 BIGipServerPool-* / WSESSIONID 等 Cookie
             └─ 存入 globalCookieJar
                 └─ 每次请求前 ensureCookieValid() 检查是否过期
                     └─ 过期则自动重新访问入口页刷新
 └─ scheduleTask()                          # 启动定时任务
     └─ crawlHouseInfo()                    # 立即执行一次爬取
         ├─ ResolveProjectParams()          # 自动解析楼盘参数
         │   ├─ postJSON()                  # 搜索楼盘（关键字"乐宸"）
         │   │   └─ getYsfYsPublicity       # 匹配 ysProjectId / preSellId
         │   └─ ResolveFybId()              # 获取目标楼栋 fybId
         │       └─ postForm()              # 匹配楼栋名称（"1栋"）
         │           └─ getBuildingNameListToPublicity
         ├─ postForm()                      # 请求楼栋字典接口
         │   └─ getBuildingDictToPublicity  # 获取可售楼层、单元分区等信息
         ├─ postJSON()                      # 请求房源列表接口
         │   └─ getHouseInfoListToPublicity # 获取房源详情及价格数据
         ├─ Analyze()                       # 数据量化分析（价格/面积/楼层）
         ├─ AnalyzeSaleStatus()             # 销售状态统计（在售/已售/比例）
         └─ PrintReport()                   # 打印分析报告
     └─ ticker 每 1800 秒循环执行 crawlHouseInfo()
```

## 分析指标说明

### 一、价格类核心指标
| 指标 | 说明 |
|------|------|
| 最低/最高挂牌单价 | 所有**在售**房源 `askpriceeachB` 字段的最小/最大值（元/㎡） |
| 挂牌单价中位数 | 所有**在售**房源单价排序后取中间值（元/㎡） |
| 备案单价区间 | `recordedPricePerUnitInside` 字段的最小/最大值，按套内面积计算，通常高于挂牌单价 |
| 最低/最高总价 | 所有**在售**房源 `askpricetotalB` 字段的最小/最大值（元） |
| 总价均值 | 所有**在售**房源总价求和后取平均（元） |
| 总价区间跨度 | 最高总价 - 最低总价，反映同楼栋价格差异（元） |

### 二、面积类核心指标
| 指标 | 说明 |
|------|------|
| 产权建筑面积 | `ysbuildingarea` 字段，即合同面积（㎡） |
| 赠送/拓展面积 | `ysexpandarea` 字段，不计入产权但可实际使用（㎡） |
| 套内面积 | `ysinsidearea` 字段，= 产权建筑面积 - 公摊面积（㎡） |
| 实际使用面积 | = 产权建筑面积 + 拓展面积（㎡） |
| 得房率 | = 套内面积 / 产权建筑面积 × 100（%） |

### 三、房源供应与楼层分布
> 仅统计 `lastStatusName` 为**期房待售**的在售房源

| 指标 | 说明 |
|------|------|
| 总在售套数 | 挂牌价大于 0 的房源数量之和 |
| 在售楼层范围 | 有在售房源的最低楼层 ~ 最高楼层 |
| 低楼层(2-9楼) | 该区间在售套数及占比 |
| 中楼层(10-19楼) | 该区间在售套数及占比 |
| 高楼层(20楼+) | 该区间在售套数及占比 |
| 单楼层最多/最少套数 | 反映在售房源的楼层分布均匀程度 |

### 四、性价比衍生指标
| 指标 | 计算公式 | 说明 |
|------|----------|------|
| 单价/拓展面积比 | 挂牌单价均值 / 拓展面积 | 每 1㎡ 拓展面积对应的单价成本，数值越低性价比越高 |
| 总价/实际使用面积比 | 总价均值 / 实际使用面积 | 反映实际使用面积的真实购房成本（元/㎡） |
| 楼层价格梯度 | (最高单价 - 最低单价) / 最低单价 × 100 | 高低楼层单价涨幅（%），涨幅越低高楼层性价比越高 |

### 五、销售状态统计
| 指标 | 说明 |
|------|------|
| 总套数 | 接口返回的全部房源数量（含在售、已售） |
| 在售套数 | `lastStatusName` 为"期房待售"的房源数量 |
| 已售套数 | "已签认购书"、"已录入合同"等非待售状态的房源数量 |
| 售出比例 | 已售套数 / 总套数 × 100（%） |
| 各状态明细 | 按套数降序列出每种状态的数量及占比 |

## Cookie 说明

Cookie 由 `cookie.go` 自动获取和管理，每 **30 分钟**自动刷新一次。

| Cookie 字段 | 来源 | 说明 |
|-------------|------|------|
| `WSESSIONID-SZFDC-SCJY` | 服务端 Session | 若需登录才能获取，则需手动更新 `config.go` 中的 `Cookies` |
| `BIGipServerPool-*` | 负载均衡 | 访问入口页时服务器自动下发 |
| `Hm_lvt_*` | 百度统计 | 访问入口页时自动下发 |
| `_trs_uv` | 访客追踪 | 访问入口页时自动下发 |

> ⚠️ 若自动获取 Cookie 失败，程序会自动降级使用 `config.go` 中配置的静态 Cookie。

## 配置说明

所有可调整参数均在 `config.go` 中：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Interval` | `1800` | 定时爬取间隔（秒） |
| `ProjectListURL` | - | 楼盘列表搜索接口地址 |
| `BuildingNameURL` | - | 楼栋名称列表接口地址 |
| `BuildingDictURL` | - | 楼栋字典接口地址 |
| `HouseInfoURL` | - | 房源列表接口地址 |
| `ProjectKeyword` | `乐宸` | 楼盘搜索关键字，切换楼盘时修改此处 |
| `BuildingName` | `1栋` | 目标楼栋名称，用于匹配 fybId |
| `HouseType` | `三房` | 户型类型筛选，空字符串表示不限户型 |