<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><Clock /></el-icon>
        <span>历史记录 & 价格趋势</span>
        <el-button link type="primary" size="small" @click="loadHistory" style="margin-left: auto">
          <el-icon><Refresh /></el-icon> 刷新
        </el-button>
      </div>
    </template>

    <!-- 趋势图 -->
    <div v-if="trendData.length > 1">
      <div class="section-title">价格趋势（近 {{ trendData.length }} 次）</div>
      <v-chart :option="trendOption" autoresize style="height: 220px; margin-bottom: 16px" />
    </div>

    <!-- 历史列表 -->
    <el-table :data="records" size="small" stripe max-height="320" v-loading="loading">
      <el-table-column prop="createdAt" label="时间" width="150" />
      <el-table-column prop="projectName" label="楼盘" />
      <el-table-column prop="buildingName" label="楼栋" width="70" />
      <el-table-column prop="houseType" label="户型" width="70" />
      <el-table-column label="中位单价" width="110" align="right">
        <template #default="{ row }">
          {{ row.analysis.median_unit_price.toLocaleString() }} 元/㎡
        </template>
      </el-table-column>
      <el-table-column label="在售/总套" width="90" align="center">
        <template #default="{ row }">
          {{ row.sale.for_sale_count }}/{{ row.sale.total_count }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="80" align="center">
        <template #default="{ row }">
          <el-popconfirm title="确认删除此记录？" @confirm="onDelete(row.id)">
            <template #reference>
              <el-button link type="danger" size="small">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Clock, Refresh } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { getHistory, deleteHistory } from '@/api'
import type { HistoryRecord } from '@/types'
import { ElMessage } from 'element-plus'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const records = ref<HistoryRecord[]>([])
const loading = ref(false)

async function loadHistory() {
  loading.value = true
  try {
    records.value = await getHistory()
  } finally {
    loading.value = false
  }
}

async function onDelete(id: string) {
  await deleteHistory(id)
  ElMessage.success('已删除')
  await loadHistory()
}

onMounted(loadHistory)

// 趋势数据：取同一楼盘+楼栋+户型的历史记录
const trendData = computed(() => {
  // 按楼盘分组，取最多的那组
  const groups: Record<string, HistoryRecord[]> = {}
  for (const r of records.value) {
    const key = `${r.projectName}_${r.buildingName}_${r.houseType}`
    if (!groups[key]) groups[key] = []
    groups[key].push(r)
  }
  const best = Object.values(groups).sort((a, b) => b.length - a.length)[0] || []
  return best.slice(-20) // 最近 20 条
})

const trendOption = computed(() => {
  const data = trendData.value
  const times = data.map((r) => r.createdAt.slice(5, 16))
  const minPrices = data.map((r) => r.analysis.min_unit_price)
  const medianPrices = data.map((r) => r.analysis.median_unit_price)
  const maxPrices = data.map((r) => r.analysis.max_unit_price)
  const forSaleRates = data.map((r) => r.sale.for_sale_rate)

  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['最低单价', '中位单价', '最高单价', '在售率%'], top: 0, textStyle: { fontSize: 11 } },
    grid: { left: 60, right: 50, top: 30, bottom: 30 },
    xAxis: { type: 'category', data: times, axisLabel: { fontSize: 10, rotate: 30 } },
    yAxis: [
      {
        type: 'value',
        name: '元/㎡',
        nameTextStyle: { fontSize: 10 },
        axisLabel: { formatter: (v: number) => (v / 10000).toFixed(1) + 'w', fontSize: 10 },
      },
      {
        type: 'value',
        name: '在售率%',
        nameTextStyle: { fontSize: 10 },
        max: 100,
        axisLabel: { formatter: '{value}%', fontSize: 10 },
      },
    ],
    series: [
      { name: '最低单价', type: 'line', data: minPrices, smooth: true, itemStyle: { color: '#67C23A' } },
      { name: '中位单价', type: 'line', data: medianPrices, smooth: true, lineStyle: { width: 2 }, itemStyle: { color: '#409EFF' } },
      { name: '最高单价', type: 'line', data: maxPrices, smooth: true, itemStyle: { color: '#F56C6C' } },
      { name: '在售率%', type: 'line', data: forSaleRates, smooth: true, yAxisIndex: 1, lineStyle: { type: 'dashed' }, itemStyle: { color: '#E6A23C' } },
    ],
  }
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.section-title {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-bottom: 4px;
}
</style>
