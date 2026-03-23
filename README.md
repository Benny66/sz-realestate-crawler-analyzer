# 深圳房产信息智能分析平台

基于 **Go + Vue 3 + Element Plus + ECharts** 的深圳房地产信息平台智能分析系统。

该项目从深圳市住房和建设局公开平台抓取预售房源信息，对房源价格、楼层、销售状态、面积、性价比等指标进行量化分析，并通过前端可视化界面展示。

---

## 项目特性

- 支持 **楼盘关键字搜索**
- 支持 **楼栋联动选择**
- 支持 **户型筛选**
- 支持 **房源价格分析**
- 支持 **销售状态统计**
- 支持 **楼层热力图展示**
- 支持 **历史趋势记录**
- 支持 **多楼盘对比分析**
- 支持 **后端缓存与定时刷新**

---

## 项目结构

```text
.
├── backend/                        # Go 后端服务
│   ├── main.go                     # 启动 HTTP 服务、定时刷新任务
│   ├── api.go                      # REST API
│   ├── cache.go                    # 缓存与历史记录持久化
│   ├── config.go                   # 全局配置
│   ├── types.go                    # 数据结构定义
│   ├── http.go                     # 请求工具
│   ├── cookie.go                   # Cookie 自动获取与刷新
│   ├── resolver.go                 # 楼盘/楼栋参数解析
│   ├── analyzer.go                 # 数据分析逻辑
│   ├── go.mod
│   └── .gitignore
│
├── frontend/                       # Vue 3 前端
│   ├── src/
│   │   ├── api/                    # 前端接口封装
│   │   ├── components/             # 搜索、图表、表格、历史、对比组件
│   │   ├── types/                  # TypeScript 类型定义
│   │   ├── views/
│   │   │   └── Dashboard.vue       # 主页面
│   │   ├── App.vue
│   │   ├── main.ts
│   │   └── style.css
│   ├── package.json
│   ├── vite.config.ts
│   └── .gitignore
│
├── README.md                       # 根目录项目说明
└── lechenginfo_fullstack_memory.md # 根目录项目记忆文件
```

---

## 技术栈

### 后端
- Go 1.23.6
- net/http
- encoding/json
- 原生文件存储（history.json）

### 前端
- Vue 3
- TypeScript
- Vite
- Element Plus
- ECharts
- Axios

---

## 核心能力

### 1. 智能搜索
前端输入楼盘关键字后：
- 调用后端搜索楼盘接口
- 自动展示候选楼盘
- 选择楼盘后再拉取楼栋列表
- 选择楼栋与户型后执行分析

### 2. 数据分析
后端对抓取到的房源数据进行分析，包括：
- 最低/最高/中位挂牌单价
- 最低/最高/平均总价
- 建筑面积 / 套内面积 / 拓展面积
- 得房率
- 低中高楼层分布
- 销售状态统计
- 性价比衍生指标

### 3. 历史趋势
每次分析完成后，会将结果写入历史记录文件：
- 默认持久化到 `backend/history.json`
- 前端可展示价格趋势图
- 支持查看与删除历史记录

### 4. 多楼盘对比
支持最多 5 个楼盘同时对比，展示：
- 综合雷达图
- 价格柱状图
- 销售状态对比图
- 指标数据表

### 5. 缓存与定时刷新
- 后端分析结果会缓存 30 分钟
- 默认定时刷新默认楼盘数据
- 可手动清除缓存

---

## 后端 API

### 楼盘搜索
```http
GET /api/search?keyword=乐宸&pageIndex=1&pageSize=12&zone=
```

### 获取楼栋列表
```http
GET /api/buildings?ysProjectId=xxx&preSellId=xxx
```

### 执行分析
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

### 多楼盘对比
```http
POST /api/compare
Content-Type: application/json
```

### 获取历史记录
```http
GET /api/history
```

### 删除历史记录
```http
DELETE /api/history?id=记录ID
```

### 清除缓存
```http
POST /api/cache/clear
```

---

## 启动方式

## 1. 启动后端

```powershell
cd backend
go run .
```

后端默认监听：

```text
http://localhost:8080
```

---

## 2. 启动前端

```powershell
cd frontend
npm install
npm run dev
```

前端默认地址通常为：

```text
http://127.0.0.1:5173
```

如果端口被占用，Vite 会自动切换到其它端口，比如 `5174`。

---

## 默认演示配置

后端默认配置：
- 楼盘关键字：`乐宸`
- 楼栋：`1栋`
- 户型：`三房`

对应配置文件：

```text
backend/config.go
```

---

## 重要说明

### Cookie 机制
由于目标站点存在会话限制，后端会：
- 启动时自动访问入口页获取 Cookie
- 每 30 分钟自动刷新 Cookie
- 如果自动获取失败，则降级使用静态 Cookie

### 历史文件
历史分析记录默认保存在：

```text
backend/history.json
```

该文件建议加入 `.gitignore`，避免把运行时数据提交到 GitHub。

### Windows 环境说明
你当前是 Windows 环境，建议使用 PowerShell 启动：

```powershell
cd backend
go run .
```

```powershell
cd frontend
npm run dev
```

---

## 后续可扩展方向

- 支持区域筛选
- 支持更多户型维度分析
- 支持收藏楼盘
- 支持导出 Excel / CSV
- 支持房源价格变化预警
- 支持登录后保存个人分析记录
- 支持部署到云服务器

---

## 适用场景

- 深圳新房预售房源观察
- 楼盘横向对比
- 购房数据辅助分析
- 销售进度趋势跟踪
- 房产信息可视化展示

---

如果你后面还要继续扩展，我建议下一步可以继续做：

1. 登录页与用户体系
2. 导出报表
3. 收藏楼盘与订阅提醒
4. Docker 部署
5. Nginx 反向代理上线
