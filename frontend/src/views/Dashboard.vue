<template>
  <div class="dashboard">
    <!-- 顶部导航栏 -->
    <el-header class="app-header">
      <div class="header-left">
        <el-icon size="24" color="#fff"><House /></el-icon>
        <span class="app-title">深圳房产信息智能分析平台</span>
        <el-tag type="info" size="small" effect="dark">深圳住建局数据</el-tag>
      </div>
      <div class="header-right">
        <el-tooltip content="清除缓存，强制重新爬取">
          <el-button size="small" @click="onClearCache" :loading="cacheClearing">
            <el-icon><RefreshRight /></el-icon> 清除缓存
          </el-button>
        </el-tooltip>
        <el-button
          size="small"
          :type="showHistory ? 'primary' : 'default'"
          @click="showHistory = !showHistory"
        >
          <el-icon><Clock /></el-icon> 历史记录
        </el-button>
        <el-tag v-if="lastUpdated" type="success" size="small">
          更新: {{ lastUpdated }}
        </el-tag>
      </div>
    </el-header>

    <el-container class="main-container">
      <!-- 左侧搜索面板 -->
      <el-aside width="280px" class="side-panel">
        <SearchPanel
          :loading="analyzing"
          :compare-loading="comparing"
          :compare-list="compareList"
          @analyze="onAnalyze"
          @add-to-compare="onAddToCompare"
          @remove-compare="onRemoveCompare"
          @clear-compare="compareList = []"
          @start-compare="onStartCompare"
        />
      </el-aside>

      <!-- 主内容区 -->
      <el-main class="main-content">
        <!-- 空状态 -->
        <div v-if="!result && !analyzing" class="empty-state">
          <el-empty description="请在左侧搜索楼盘并点击「开始分析」">
            <template #image>
              <el-icon size="80" color="#dcdfe6"><House /></el-icon>
            </template>
            <el-button type="primary" @click="onQuickSearch">
              快速查询「乐宸花园 1栋 三房」
            </el-button>
          </el-empty>
        </div>

        <!-- 加载中 -->
        <div v-else-if="analyzing" class="loading-state">
          <el-skeleton :rows="8" animated />
        </div>

        <!-- 分析结果 -->
        <template v-else-if="result">
          <!-- 楼盘信息标题 -->
          <div class="result-header">
            <div class="result-title">
              <el-icon><OfficeBuilding /></el-icon>
              {{ result.params.projectName }}
              <el-tag size="small">{{ currentReq?.buildingName }}</el-tag>
              <el-tag size="small" type="success">{{ currentReq?.houseType || '全部户型' }}</el-tag>
            </div>
            <div class="result-meta">
              数据更新时间: {{ result.updatedAt }}
            </div>
          </div>

          <!-- 核心指标卡片 -->
          <MetricCards :analysis="result.analysis" style="margin-bottom: 16px" />

          <!-- 图表区域 -->
          <el-row :gutter="16" style="margin-bottom: 16px">
            <el-col :span="10">
              <SaleStatusChart :sale="result.sale" />
            </el-col>
            <el-col :span="14">
              <PriceScatter :houses="result.houses" :avg-price="result.analysis.median_unit_price" />
            </el-col>
          </el-row>

          <el-row :gutter="16" style="margin-bottom: 16px">
            <el-col :span="24">
              <FloorHeatmap :floor-groups="result.floorGroups" />
            </el-col>
          </el-row>

          <!-- 房源明细表 -->
          <HouseTable
            :houses="result.houses"
            :avg-price="result.analysis.median_unit_price"
            style="margin-bottom: 16px"
          />
        </template>

        <!-- 历史记录面板 -->
        <div v-if="showHistory" class="history-overlay">
          <HistoryPanel />
        </div>

        <!-- 多楼盘对比面板 -->
        <el-dialog
          v-model="showCompare"
          title=""
          width="90%"
          top="4vh"
          :close-on-click-modal="false"
          destroy-on-close
        >
          <CompareChart
            v-if="compareResult"
            :items="compareResult.items"
            @close="showCompare = false"
          />
          <div v-else-if="comparing" style="text-align: center; padding: 40px">
            <el-icon class="is-loading" size="40"><Loading /></el-icon>
            <div style="margin-top: 12px; color: #909399">正在对比分析中，请稍候...</div>
          </div>
        </el-dialog>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  House,
  OfficeBuilding,
  Clock,
  RefreshRight,
  Loading,
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import SearchPanel from '@/components/SearchPanel.vue'
import MetricCards from '@/components/MetricCards.vue'
import SaleStatusChart from '@/components/SaleStatusChart.vue'
import FloorHeatmap from '@/components/FloorHeatmap.vue'
import PriceScatter from '@/components/PriceScatter.vue'
import HouseTable from '@/components/HouseTable.vue'
import HistoryPanel from '@/components/HistoryPanel.vue'
import CompareChart from '@/components/CompareChart.vue'
import { analyzeProject, compareProjects, clearCache } from '@/api'
import type { AnalyzeRequest, AnalyzeResponse, CompareResponse } from '@/types'

const result = ref<AnalyzeResponse | null>(null)
const analyzing = ref(false)
const lastUpdated = ref('')
const currentReq = ref<AnalyzeRequest | null>(null)
const showHistory = ref(false)
const cacheClearing = ref(false)

// 对比相关
const compareList = ref<AnalyzeRequest[]>([])
const comparing = ref(false)
const showCompare = ref(false)
const compareResult = ref<CompareResponse | null>(null)

async function onAnalyze(req: AnalyzeRequest) {
  analyzing.value = true
  currentReq.value = req
  try {
    result.value = await analyzeProject(req)
    lastUpdated.value = result.value.updatedAt
    ElMessage.success(`分析完成：${result.value.params.projectName}`)
  } catch {
    // 错误已在 api/index.ts 拦截器中处理
  } finally {
    analyzing.value = false
  }
}

function onQuickSearch() {
  onAnalyze({ keyword: '乐宸', buildingName: '1栋', houseType: '三房' })
}

function onAddToCompare(req: AnalyzeRequest) {
  if (compareList.value.length >= 5) {
    ElMessage.warning('最多支持 5 个楼盘同时对比')
    return
  }
  const exists = compareList.value.some(
    (i) => i.keyword === req.keyword && i.buildingName === req.buildingName && i.houseType === req.houseType
  )
  if (exists) {
    ElMessage.warning('该楼盘已在对比列表中')
    return
  }
  compareList.value.push(req)
  ElMessage.success(`已加入对比：${req.keyword} ${req.buildingName}`)
}

function onRemoveCompare(idx: number) {
  compareList.value.splice(idx, 1)
}

async function onStartCompare() {
  if (compareList.value.length < 2) {
    ElMessage.warning('请至少添加 2 个楼盘进行对比')
    return
  }
  comparing.value = true
  showCompare.value = true
  compareResult.value = null
  try {
    compareResult.value = await compareProjects({ items: compareList.value })
  } catch {
    showCompare.value = false
  } finally {
    comparing.value = false
  }
}

async function onClearCache() {
  cacheClearing.value = true
  try {
    await clearCache()
    ElMessage.success('缓存已清除，下次查询将重新爬取')
  } finally {
    cacheClearing.value = false
  }
}
</script>

<style scoped>
.dashboard {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.app-header {
  background: linear-gradient(135deg, #1a6fc4 0%, #2196f3 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 56px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-title {
  color: #fff;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-container {
  flex: 1;
  overflow: hidden;
}

.side-panel {
  padding: 16px 12px 16px 16px;
  overflow-y: auto;
  background: #f5f7fa;
  border-right: 1px solid #e4e7ed;
}

.main-content {
  overflow-y: auto;
  padding: 16px;
  position: relative;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60vh;
}

.loading-state {
  padding: 20px;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.result-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 700;
  color: #303133;
}

.result-meta {
  font-size: 12px;
  color: #909399;
}

.history-overlay {
  margin-top: 16px;
}
</style>
