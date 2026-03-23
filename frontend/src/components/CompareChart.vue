<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><ScaleToOriginal /></el-icon>
        <span>多楼盘对比分析</span>
        <el-tag size="small" type="warning">{{ items.length }} 个楼盘</el-tag>
        <el-button link size="small" @click="$emit('close')" style="margin-left: auto">
          <el-icon><Close /></el-icon> 关闭
        </el-button>
      </div>
    </template>

    <el-tabs v-model="activeTab">
      <!-- 雷达图 -->
      <el-tab-pane label="综合雷达图" name="radar">
        <v-chart :option="radarOption" autoresize style="height: 380px" />
      </el-tab-pane>

      <!-- 价格对比柱状图 -->
      <el-tab-pane label="价格对比" name="price">
        <v-chart :option="priceOption" autoresize style="height: 320px" />
      </el-tab-pane>

      <!-- 销售状态对比 -->
      <el-tab-pane label="销售状态" name="sale">
        <v-chart :option="saleOption" autoresize style="height: 320px" />
      </el-tab-pane>

      <!-- 数据表格 -->
      <el-tab-pane label="数据明细" name="table">
        <el-table :data="tableData" border size="small">
          <el-table-column prop="label" label="指标" width="160" fixed />
          <el-table-column
            v-for="(item, idx) in items"
            :key="idx"
            :label="`${item.projectName} ${item.buildingName}`"
            align="right"
          >
            <template #default="{ row }">
              <span :class="row.best === idx ? 'best-value' : ''">
                {{ row.values[idx] }}
              </span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ScaleToOriginal, Close } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { RadarChart, BarChart } from 'echarts/charts'
import {
  RadarComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { CompareItem } from '@/types'

use([RadarChart, BarChart, RadarComponent, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{ items: CompareItem[] }>()
defineEmits(['close'])

const activeTab = ref('radar')
const COLORS = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#9B59B6']

function label(item: CompareItem) {
  return `${item.projectName}\n${item.buildingName}${item.houseType}`
}

// 雷达图：归一化各指标
const radarOption = computed(() => {
  const indicators = [
    { name: '中位单价', key: 'median_unit_price', inverse: true },
    { name: '得房率', key: 'housing_rate', inverse: false },
    { name: '赠送面积', key: 'expand_area', inverse: false },
    { name: '在售率', key: 'for_sale_rate', inverse: false },
    { name: '价格梯度', key: 'floor_price_premium', inverse: true },
    { name: '实际成本', key: 'cost_per_actual_area', inverse: true },
  ]

  const rawValues = indicators.map((ind) =>
    props.items.map((item) =>
      ind.key in item.analysis
        ? (item.analysis as any)[ind.key]
        : (item.sale as any)[ind.key] ?? 0
    )
  )

  // 归一化到 0-100
  const normalized = rawValues.map((vals, i) => {
    const min = Math.min(...vals)
    const max = Math.max(...vals)
    if (max === min) return vals.map(() => 50)
    return vals.map((v) => {
      const norm = ((v - min) / (max - min)) * 100
      return indicators[i].inverse ? 100 - norm : norm
    })
  })

  const seriesData = props.items.map((item, idx) => ({
    name: label(item),
    value: normalized.map((vals) => Math.round(vals[idx])),
    itemStyle: { color: COLORS[idx] },
    areaStyle: { opacity: 0.15 },
  }))

  return {
    tooltip: {},
    legend: { data: props.items.map(label), bottom: 0, textStyle: { fontSize: 11 } },
    radar: {
      indicator: indicators.map((ind) => ({ name: ind.name, max: 100 })),
      center: ['50%', '48%'],
      radius: '60%',
    },
    series: [{ type: 'radar', data: seriesData }],
  }
})

// 价格对比柱状图
const priceOption = computed(() => {
  const labels = props.items.map(label)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['最低单价', '中位单价', '最高单价'], top: 0 },
    grid: { left: 70, right: 20, top: 40, bottom: 60 },
    xAxis: { type: 'category', data: labels, axisLabel: { rotate: 15, fontSize: 11 } },
    yAxis: {
      type: 'value',
      name: '元/㎡',
      axisLabel: { formatter: (v: number) => (v / 10000).toFixed(1) + 'w' },
    },
    series: [
      { name: '最低单价', type: 'bar', data: props.items.map((i) => i.analysis.min_unit_price), itemStyle: { color: '#67C23A' } },
      { name: '中位单价', type: 'bar', data: props.items.map((i) => i.analysis.median_unit_price), itemStyle: { color: '#409EFF' } },
      { name: '最高单价', type: 'bar', data: props.items.map((i) => i.analysis.max_unit_price), itemStyle: { color: '#F56C6C' } },
    ],
  }
})

// 销售状态对比
const saleOption = computed(() => {
  const labels = props.items.map(label)
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['在售套数', '已售套数', '在售率%'], top: 0 },
    grid: { left: 60, right: 60, top: 40, bottom: 60 },
    xAxis: { type: 'category', data: labels, axisLabel: { rotate: 15, fontSize: 11 } },
    yAxis: [
      { type: 'value', name: '套数' },
      { type: 'value', name: '在售率%', max: 100, axisLabel: { formatter: '{value}%' } },
    ],
    series: [
      { name: '在售套数', type: 'bar', data: props.items.map((i) => i.sale.for_sale_count), itemStyle: { color: '#409EFF' } },
      { name: '已售套数', type: 'bar', data: props.items.map((i) => i.sale.sold_count), itemStyle: { color: '#dcdfe6' } },
      { name: '在售率%', type: 'line', yAxisIndex: 1, data: props.items.map((i) => i.sale.for_sale_rate), itemStyle: { color: '#E6A23C' } },
    ],
  }
})

// 数据明细表
const tableData = computed(() => {
  const rows = [
    { label: '中位单价(元/㎡)', values: props.items.map((i) => i.analysis.median_unit_price.toLocaleString()), bestFn: (vals: number[]) => vals.indexOf(Math.min(...vals)) },
    { label: '最低总价(万)', values: props.items.map((i) => (i.analysis.min_total_price / 10000).toFixed(2)), bestFn: (vals: number[]) => vals.indexOf(Math.min(...vals)) },
    { label: '得房率(%)', values: props.items.map((i) => i.analysis.housing_rate.toFixed(2) + '%'), bestFn: (vals: number[]) => vals.indexOf(Math.max(...vals)) },
    { label: '赠送面积(㎡)', values: props.items.map((i) => i.analysis.expand_area.toFixed(2)), bestFn: (vals: number[]) => vals.indexOf(Math.max(...vals)) },
    { label: '实际使用面积(㎡)', values: props.items.map((i) => i.analysis.actual_use_area.toFixed(2)), bestFn: (vals: number[]) => vals.indexOf(Math.max(...vals)) },
    { label: '在售套数', values: props.items.map((i) => String(i.sale.for_sale_count)), bestFn: (vals: number[]) => vals.indexOf(Math.max(...vals)) },
    { label: '在售率(%)', values: props.items.map((i) => i.sale.for_sale_rate.toFixed(2) + '%'), bestFn: (vals: number[]) => vals.indexOf(Math.max(...vals)) },
    { label: '楼层价格梯度(%)', values: props.items.map((i) => i.analysis.floor_price_premium.toFixed(2) + '%'), bestFn: (vals: number[]) => vals.indexOf(Math.min(...vals)) },
  ]
  return rows.map((row) => {
    const nums = row.values.map((v) => parseFloat(v.replace(/[^0-9.]/g, '')))
    return { ...row, best: row.bestFn(nums) }
  })
})
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.best-value {
  color: var(--el-color-success);
  font-weight: 700;
}
</style>
