<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><Grid /></el-icon>
        <span>楼层分布 & 单价热力图</span>
      </div>
    </template>
    <v-chart :option="option" autoresize style="height: 380px" />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Grid } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart, ScatterChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  VisualMapComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { FloorGroup } from '@/types'

use([BarChart, ScatterChart, GridComponent, TooltipComponent, LegendComponent, VisualMapComponent, CanvasRenderer])

const props = defineProps<{ floorGroups: FloorGroup[] }>()

const option = computed(() => {
  // 按楼层排序
  const sorted = [...props.floorGroups].sort((a, b) => Number(a.floor) - Number(b.floor))

  const floors = sorted.map((g) => `${g.floor}层`)
  const forSaleCounts = sorted.map((g) => g.list.filter((h) => h.askpriceeachB > 0).length)
  const soldCounts = sorted.map((g) => g.list.filter((h) => h.askpriceeachB <= 0).length)
  const avgPrices = sorted.map((g) => {
    const onSale = g.list.filter((h) => h.askpriceeachB > 0)
    if (!onSale.length) return 0
    return Math.round(onSale.reduce((s, h) => s + h.askpriceeachB, 0) / onSale.length)
  })

  const maxPrice = Math.max(...avgPrices.filter((p) => p > 0))
  const minPrice = Math.min(...avgPrices.filter((p) => p > 0))

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter(params: any[]) {
        const floor = params[0]?.axisValue
        const lines = params.map((p: any) => `${p.marker}${p.seriesName}: ${p.value}`)
        const priceIdx = avgPrices[params[0]?.dataIndex]
        if (priceIdx > 0) lines.push(`<br/>均价: ${priceIdx.toLocaleString()} 元/㎡`)
        return `<b>${floor}</b><br/>${lines.join('<br/>')}`
      },
    },
    legend: { data: ['在售', '已售'], top: 0 },
    grid: { left: 60, right: 80, top: 40, bottom: 20 },
    xAxis: {
      type: 'value',
      name: '套数',
      nameTextStyle: { fontSize: 11 },
    },
    yAxis: {
      type: 'category',
      data: floors,
      axisLabel: { fontSize: 11 },
    },
    visualMap: {
      show: false,
      min: minPrice,
      max: maxPrice,
      inRange: { color: ['#91d5ff', '#096dd9'] },
      seriesIndex: 2,
    },
    series: [
      {
        name: '在售',
        type: 'bar',
        stack: 'total',
        data: forSaleCounts,
        itemStyle: { color: '#409EFF', borderRadius: [0, 3, 3, 0] },
        label: {
          show: true,
          position: 'insideRight',
          formatter: (p: any) => (p.value > 0 ? p.value : ''),
          fontSize: 11,
        },
      },
      {
        name: '已售',
        type: 'bar',
        stack: 'total',
        data: soldCounts,
        itemStyle: { color: '#dcdfe6', borderRadius: [0, 3, 3, 0] },
      },
      {
        name: '均价',
        type: 'scatter',
        yAxisIndex: 0,
        data: avgPrices.map((p, i) => (p > 0 ? [p / 1000, floors[i]] : null)).filter(Boolean),
        symbolSize: 10,
        tooltip: {
          formatter: (p: any) => `${p.name}<br/>均价: ${(p.value[0] * 1000).toLocaleString()} 元/㎡`,
        },
      },
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
</style>
