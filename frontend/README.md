# frontend - 深圳房产信息智能分析平台前端

该目录为项目的 Vue 3 前端应用，负责：

- 搜索楼盘与楼栋
- 发起单楼盘分析请求
- 展示价格、面积、销售状态等可视化结果
- 展示历史趋势图
- 支持多楼盘对比分析

---

## 技术栈

- Vue 3
- TypeScript
- Vite
- Element Plus
- ECharts
- Axios

---

## 目录结构

```text
frontend/
├── src/
│   ├── api/
│   │   └── index.ts             # Axios 接口封装
│   ├── components/
│   │   ├── SearchPanel.vue      # 搜索面板
│   │   ├── MetricCards.vue      # 核心指标卡片
│   │   ├── SaleStatusChart.vue  # 销售状态图
│   │   ├── FloorHeatmap.vue     # 楼层热力图
│   │   ├── PriceScatter.vue     # 单价散点图
│   │   ├── HouseTable.vue       # 房源明细表格
│   │   ├── HistoryPanel.vue     # 历史记录与趋势图
│   │   └── CompareChart.vue     # 多楼盘对比图
│   ├── types/
│   │   └── index.ts             # 类型定义
│   ├── views/
│   │   └── Dashboard.vue        # 主页面
│   ├── App.vue
│   ├── main.ts
│   └── style.css
├── public/
├── package.json
├── vite.config.ts
└── .gitignore
```

---

## 安装依赖

```powershell
cd frontend
npm install
```

---

## 启动开发环境

```powershell
cd frontend
npm run dev
```

默认访问地址通常为：

```text
http://127.0.0.1:5173
```

如果端口占用，Vite 会自动切换到其它端口，例如：`5174`。

---

## 打包

```powershell
cd frontend
npm run build
```

> 说明：你当前本机 Node 版本较低时，某些依赖版本可能导致 `build` 有兼容性问题。开发模式通常可以正常运行，后续如需稳定打包，建议升级到 Node 18+，最好 Node 20 LTS。

---

## 与后端联调

前端通过 Vite 代理访问后端 API。

配置文件：

```text
frontend/vite.config.ts
```

代理规则：
- `/api` → `http://localhost:8080`

因此开发时请先启动后端：

```powershell
cd backend
go run .
```

再启动前端：

```powershell
cd frontend
npm run dev
```

---

## 核心页面说明

### `Dashboard.vue`
主页面，负责整合：
- 搜索
- 分析结果展示
- 图表展示
- 房源表格
- 历史记录面板
- 对比弹窗

### `SearchPanel.vue`
负责：
- 楼盘关键字输入
- 自动补全候选列表
- 楼栋联动选择
- 户型选择
- 发起分析
- 加入对比列表

### `MetricCards.vue`
展示三类核心指标：
- 价格指标
- 面积指标
- 供应指标

### `SaleStatusChart.vue`
展示：
- 销售状态环图
- 状态明细表

### `FloorHeatmap.vue`
展示：
- 各楼层在售 / 已售套数
- 楼层均价热力分布

### `PriceScatter.vue`
展示：
- 楼层与单价散点关系
- 均价参考线
- 高价/低价房源颜色区分

### `HouseTable.vue`
展示：
- 房源明细表格
- 状态筛选
- 房号/楼层搜索
- 排序功能

### `HistoryPanel.vue`
展示：
- 历史记录列表
- 单价趋势图
- 在售率趋势图
- 删除记录

### `CompareChart.vue`
展示多楼盘对比：
- 综合雷达图
- 价格柱状图
- 销售状态图
- 指标表格

---

## API 封装

文件：`src/api/index.ts`

封装的接口包括：
- `searchProjects()`
- `getBuildings()`
- `analyzeProject()`
- `compareProjects()`
- `getHistory()`
- `deleteHistory()`
- `clearCache()`

特点：
- 统一通过 Axios 请求
- 使用 Element Plus `ElMessage` 做错误提示
- 自动处理后端统一响应结构

---

## 界面功能

当前前端已支持：

- [x] 楼盘智能搜索
- [x] 楼栋下拉联动
- [x] 户型筛选
- [x] 单楼盘分析
- [x] 核心指标卡片展示
- [x] 销售状态饼图
- [x] 楼层热力图
- [x] 单价散点图
- [x] 房源表格检索
- [x] 历史趋势图
- [x] 多楼盘对比分析
- [x] 清缓存按钮

---

## 开发建议

### 推荐 Node 版本
建议使用：
- Node 18+
- 最佳为 Node 20 LTS

### 推荐 IDE
- VS Code
- 安装 Volar
- 安装 TypeScript Vue Plugin

### 常见开发入口

如果要改：
- 页面布局：`src/views/Dashboard.vue`
- 搜索交互：`src/components/SearchPanel.vue`
- 图表展示：`src/components/*.vue`
- 接口调用：`src/api/index.ts`
- 类型定义：`src/types/index.ts`
- 样式与全局风格：`src/style.css`

---

## Git 忽略说明

前端已忽略：
- `node_modules/`
- `dist/`
- `coverage/`
- `.vite/`
- `.env*`
- 编辑器临时文件

具体规则见：

```text
frontend/.gitignore
```

---

## 后续可扩展方向

前端后续适合继续做：

- 更丰富的筛选器（区域、状态、用途）
- 楼盘收藏
- 导出图片 / 导出报表
- 登录态与用户中心
- 深色模式
- 移动端适配
- 图表联动筛选
