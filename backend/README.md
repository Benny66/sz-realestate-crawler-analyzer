# backend - 深圳房产信息智能分析平台后端

该目录为项目的 Go 后端服务，负责：

- 爬取深圳房地产信息平台公开数据
- 自动解析楼盘与楼栋参数
- 分析房源价格、面积、楼层与销售状态
- 提供前端调用的 REST API
- 缓存分析结果
- 记录历史趋势数据

---

## 技术栈

- Go 1.23.6
- `net/http`
- `encoding/json`
- 原生文件存储（`history.json`）

---

## 目录说明

```text
backend/
├── main.go       # 启动 HTTP 服务、初始化 Cookie、定时刷新
├── api.go        # API 路由与请求处理
├── cache.go      # 内存缓存与历史记录持久化
├── config.go     # 默认配置与请求头配置
├── types.go      # 所有结构体定义
├── http.go       # POST 请求工具函数
├── cookie.go     # Cookie 自动获取与刷新
├── resolver.go   # 楼盘 / 楼栋参数自动解析
├── analyzer.go   # 数据分析逻辑
├── go.mod
└── .gitignore
```

---

## 启动方式

在 Windows PowerShell 中执行：

```powershell
cd backend
go run .
```

启动后默认监听：

```text
http://localhost:8080
```

---

## 运行流程

```text
main()
 ├─ initCookieClient()            # 自动获取 Cookie
 ├─ RegisterRoutes()              # 注册 API
 ├─ go scheduleAutoRefresh()      # 后台定时刷新默认楼盘数据
 └─ http.ListenAndServe(:8080)    # 启动 HTTP Server
```

---

## 默认分析配置

定义在：`config.go`

```go
ProjectKeyword = "乐宸"
BuildingName   = "1栋"
HouseType      = "三房"
Interval       = 3600
ServerAddr     = ":8080"
```

说明：
- `ProjectKeyword`：楼盘关键字
- `BuildingName`：楼栋名称
- `HouseType`：户型
- `Interval`：定时刷新间隔（秒）
- `ServerAddr`：服务监听端口

---

## API 列表

### 1. 搜索楼盘

```http
GET /api/search?keyword=乐宸&pageIndex=1&pageSize=12&zone=
```

返回：楼盘列表

---

### 2. 获取楼栋列表

```http
GET /api/buildings?ysProjectId=xxx&preSellId=xxx
```

返回：楼栋列表

---

### 3. 执行单楼盘分析

```http
POST /api/analyze
Content-Type: application/json
```

请求体示例：

```json
{
  "keyword": "乐宸",
  "buildingName": "1栋",
  "houseType": "三房"
}
```

返回内容包括：
- `analysis`：分析结果
- `sale`：销售状态汇总
- `houses`：原始房源明细
- `floorGroups`：楼层分组数据
- `params`：解析出的楼盘参数
- `updatedAt`：更新时间

---

### 4. 多楼盘对比

```http
POST /api/compare
Content-Type: application/json
```

请求体示例：

```json
{
  "items": [
    {
      "keyword": "乐宸",
      "buildingName": "1栋",
      "houseType": "三房"
    },
    {
      "keyword": "某某花园",
      "buildingName": "2栋",
      "houseType": "三房"
    }
  ]
}
```

最多支持 5 个楼盘同时对比。

---

### 5. 获取历史记录

```http
GET /api/history
```

可选参数：

```http
GET /api/history?projectName=乐宸
```

---

### 6. 删除历史记录

```http
DELETE /api/history?id=记录ID
```

---

### 7. 清除缓存

```http
POST /api/cache/clear
```

---

## 核心逻辑说明

### `doAnalyze()`
定义位置：`api.go`

这是当前项目最重要的后端入口函数，负责：

1. 解析楼盘与楼栋参数
2. 调用房源接口抓取数据
3. 调用 `Analyze()` 执行统计分析
4. 调用 `AnalyzeSaleStatus()` 统计销售状态
5. 写入缓存
6. 写入历史记录
7. 返回前端需要的完整响应

---

## 缓存机制

定义位置：`cache.go`

### 缓存 key 组成

```text
projectName + buildingName + houseType
```

### 默认规则

- 缓存有效期：30 分钟
- 命中缓存时直接返回结果
- 避免频繁请求目标站点

---

## 历史记录机制

默认文件：

```text
backend/history.json
```

说明：
- 每次分析成功后会写入一条记录
- 最多保留 500 条历史记录
- 供前端趋势图使用

注意：
- `history.json` 已加入 `.gitignore`
- 不建议上传到 GitHub

---

## Cookie 说明

目标站点需要有效 Cookie 才能稳定请求。

后端会：
- 启动时自动访问入口页获取 Cookie
- 每 30 分钟自动刷新 Cookie
- 若自动获取失败，使用 `config.go` 中静态 Cookie 降级

核心文件：`cookie.go`

---

## 开发建议

### 调试建议
如果接口异常，优先检查：

1. Cookie 是否失效
2. 目标接口是否变更
3. `keyword / buildingName / houseType` 是否正确
4. 目标站点是否限流

### 常见开发入口

- 新增接口：修改 `api.go`
- 调整默认参数：修改 `config.go`
- 扩展分析指标：修改 `analyzer.go`
- 修改数据结构：修改 `types.go`
- 修改缓存/历史逻辑：修改 `cache.go`

---

## Git 忽略说明

后端已忽略以下文件：

- `history.json`
- `*.exe`
- `*.out`
- `bin/`
- `.env*`
- IDE 临时文件

如需查看规则：

```text
backend/.gitignore
```

---

## 后续建议

后端后续适合继续扩展：

- 数据导出 CSV / Excel
- 价格波动预警
- 分页与限流
- 用户收藏与订阅
- Docker 化部署
- 数据库存储历史记录
