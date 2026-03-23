# lechenginfo 全栈项目记忆文档

> 深圳房产信息智能分析平台（全栈版）  
> 后端：Go Module `sz-realestate-crawler-analyzer`  
> 前端：Vue 3 + TypeScript + Element Plus + ECharts  
> 根目录结构：`backend/` + `frontend/`

---

## 一、项目概览

该项目已从原始命令行爬虫升级为：

- **后端 API 服务**：负责爬取、解析、分析、缓存、历史记录
- **前端可视化平台**：负责搜索、图表展示、表格展示、趋势分析、多楼盘对比

默认分析目标仍然是：
- 楼盘关键字：**乐宸**
- 楼栋：**1栋**
- 户型：**三房**

---

## 二、目录结构

```text
.
├── backend/
│   ├── main.go
│   ├── api.go
│   ├── cache.go
│   ├── analyzer.go
│   ├── resolver.go
│   ├── http.go
│   ├── cookie.go
│   ├── config.go
│   ├── types.go
│   ├── history.json
│   ├── go.mod
│   └── .gitignore
│
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   │   └── index.ts
│   │   ├── components/
│   │   │   ├── SearchPanel.vue
│   │   │   ├── MetricCards.vue
│   │   │   ├── SaleStatusChart.vue
│   │   │   ├── FloorHeatmap.vue
│   │   │   ├── PriceScatter.vue
│   │   │   ├── HouseTable.vue
│   │   │   ├── HistoryPanel.vue
│   │   │   └── CompareChart.vue
│   │   ├── types/
│   │   │   └── index.ts
│   │   ├── views/
│   │   │   └── Dashboard.vue
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── package.json
│   ├── vite.config.ts
│   └── .gitignore
│
├── README.md
└── lechenginfo_fullstack_memory.md
```

---

## 三、后端核心文件说明

### `backend/main.go`
负责：
- 启动 HTTP Server
- 初始化 Cookie
- 注册 API 路由
- 启动定时刷新任务

关键点：
- 已修复旧版 `scheduleTask()` 中 `return` 导致定时器不生效的问题
- 现在会真正进行周期性刷新

### `backend/api.go`
负责：
- 提供前端调用的 REST API
- 处理搜索、分析、历史记录、对比、清缓存

API 列表：
- `GET /api/search`
- `GET /api/buildings`
- `POST /api/analyze`
- `POST /api/compare`
- `GET /api/history`
- `DELETE /api/history`
- `POST /api/cache/clear`

### `backend/cache.go`
负责：
- 内存缓存分析结果
- 持久化历史记录到 `history.json`

默认逻辑：
- 缓存有效期：30 分钟
- 历史记录最多保留 500 条

### `backend/config.go`
负责：
- 服务监听地址
- 默认定时间隔
- 默认请求头
- 默认查询参数
- 静态 Cookie 备用值

新增字段：
- `ServerAddr = ":8080"`

### `backend/resolver.go`
负责：
- 根据关键字搜索楼盘
- 根据楼栋名称解析 `fybId`

核心函数：
- `ResolveProjectParams(keyword)`
- `ResolveFybId(ysProjectId, preSellId, buildingName)`

### `backend/analyzer.go`
负责：
- 价格分析
- 面积分析
- 楼层分析
- 销售状态统计
- 终端调试输出

核心函数：
- `Analyze(groups)`
- `AnalyzeSaleStatus(groups)`
- `PrintReport(r, s)`

### `backend/cookie.go`
负责：
- 自动获取 Cookie
- 自动刷新 Cookie
- Cookie 失效时降级为静态 Cookie

刷新周期：
- `30 分钟`

---

## 四、前端核心文件说明

### `frontend/src/views/Dashboard.vue`
主页面，负责：
- 页面布局
- 搜索分析
- 加载分析结果
- 展示图表与表格
- 打开历史面板
- 打开对比面板

### `frontend/src/components/SearchPanel.vue`
负责：
- 楼盘关键字自动补全
- 楼栋联动选择
- 户型选择
- 发起分析请求
- 将楼盘加入对比列表

### `frontend/src/components/MetricCards.vue`
负责展示三大指标卡：
- 价格指标
- 面积指标
- 供应指标

### `frontend/src/components/SaleStatusChart.vue`
负责：
- 销售状态饼图/环图
- 各状态明细表格

### `frontend/src/components/FloorHeatmap.vue`
负责：
- 楼层分布可视化
- 在售/已售套数堆叠条形图
- 楼层均价点图

### `frontend/src/components/PriceScatter.vue`
负责：
- 楼层 × 单价散点图
- 标记均价线
- 颜色区分低价/高价房源

### `frontend/src/components/HouseTable.vue`
负责：
- 房源明细表格展示
- 按状态筛选
- 按关键词搜索
- 单价、总价排序

### `frontend/src/components/HistoryPanel.vue`
负责：
- 加载历史记录
- 删除历史记录
- 展示最近趋势折线图

### `frontend/src/components/CompareChart.vue`
负责：
- 多楼盘雷达图对比
- 价格柱状图对比
- 销售状态对比
- 指标表格对比

### `frontend/src/api/index.ts`
负责：
- Axios 请求封装
- 统一错误提示
- 所有后端 API 调用方法

---

## 五、核心接口说明

### 1. 搜索楼盘
```http
GET /api/search?keyword=乐宸
```
返回楼盘列表，供前端自动补全。

### 2. 获取楼栋
```http
GET /api/buildings?ysProjectId=xxx&preSellId=xxx
```
返回楼栋列表，供前端楼栋下拉选择。

### 3. 分析单个楼盘
```http
POST /api/analyze
```
请求体：
```json
{
  "keyword": "乐宸",
  "buildingName": "1栋",
  "houseType": "三房"
}
```

返回内容包含：
- `analysis`
- `sale`
- `houses`
- `floorGroups`
- `params`
- `updatedAt`

### 4. 多楼盘对比
```http
POST /api/compare
```
请求体：
```json
{
  "items": [
    {"keyword": "乐宸", "buildingName": "1栋", "houseType": "三房"},
    {"keyword": "其它楼盘", "buildingName": "2栋", "houseType": "三房"}
  ]
}
```

### 5. 历史记录
```http
GET /api/history
DELETE /api/history?id=xxx
```

### 6. 清空缓存
```http
POST /api/cache/clear
```

---

## 六、数据流说明

### 单楼盘分析流程

```text
前端 SearchPanel.vue
  └─ 输入关键字
      └─ 调用 /api/search
          └─ 获取楼盘候选列表
              └─ 选择楼盘
                  └─ 调用 /api/buildings
                      └─ 获取楼栋列表
                          └─ 用户选择楼栋和户型
                              └─ 调用 /api/analyze
                                  └─ 后端 doAnalyze()
                                      ├─ ResolveProjectParams()
                                      ├─ ResolveFybId()
                                      ├─ postJSON(getHouseInfoListToPublicity)
                                      ├─ Analyze()
                                      ├─ AnalyzeSaleStatus()
                                      ├─ SetCache()
                                      └─ SaveHistory()
                                          └─ 前端展示图表与表格
```

### 多楼盘对比流程

```text
前端 SearchPanel.vue
  └─ 加入多个楼盘到 compareList
      └─ 调用 /api/compare
          └─ 后端逐个执行 doAnalyze()
              └─ 汇总 CompareResponse
                  └─ 前端 CompareChart.vue 展示
```

### 历史趋势流程

```text
每次 doAnalyze() 成功
  └─ SaveHistory(record)
      └─ 写入 backend/history.json
          └─ 前端 HistoryPanel.vue 调用 /api/history
              └─ 绘制趋势折线图
```

---

## 七、关键配置

### 后端默认配置
文件：`backend/config.go`

```go
ServerAddr      = ":8080"
Interval        = 3600
ProjectKeyword  = "乐宸"
BuildingName    = "1栋"
HouseType       = "三房"
```

### 前端代理配置
文件：`frontend/vite.config.ts`

目标：
- `/api` 代理到 `http://localhost:8080`

---

## 八、缓存与历史规则

### 缓存规则
- Key 组成：`projectName + buildingName + houseType`
- TTL：30 分钟
- 命中缓存时直接返回分析结果

### 历史记录规则
- 文件：`backend/history.json`
- 最多保留：500 条
- 每次成功分析后写入一条历史记录

---

## 九、Git 忽略建议

### backend 应忽略
- `history.json`
- `*.exe`
- `*.out`
- `bin/`
- 覆盖率文件
- `.env`

### frontend 应忽略
- `node_modules/`
- `dist/`
- `.vite/`
- `coverage/`
- `.env`

### 根目录应忽略
- `.vscode/`
- `.idea/`
- 系统临时文件

---

## 十、当前版本状态

当前项目已经完成：

- [x] Go 后端 API 化
- [x] Vue 3 前端展示
- [x] Element Plus UI
- [x] ECharts 图表展示
- [x] 历史趋势功能
- [x] 多楼盘对比功能
- [x] 缓存与定时刷新
- [x] 前后端目录拆分

---

## 十一、后续建议

下一步建议优先做：

1. 导出 Excel / CSV
2. 收藏楼盘功能
3. 价格波动预警
4. 多条件筛选（区域/用途/状态）
5. Docker 化部署
6. 上线 Nginx + HTTPS
7. 登录系统与个人收藏记录

---

## 十二、维护提示

如果后面继续开发，请记住：

- 后端核心分析入口：`doAnalyze()` in `backend/api.go`
- 历史记录文件：`backend/history.json`
- 前端主页面：`frontend/src/views/Dashboard.vue`
- 多楼盘对比组件：`frontend/src/components/CompareChart.vue`
- 趋势图组件：`frontend/src/components/HistoryPanel.vue`
- 楼盘智能搜索入口：`frontend/src/components/SearchPanel.vue`

如果以后继续扩展新功能，应优先在这个全栈记忆文档基础上更新，而不是沿用旧的命令行版本记忆。
