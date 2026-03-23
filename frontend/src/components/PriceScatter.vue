<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><TrendCharts /></el-icon>
        <span>楼层 × 单价散点图</span>
        <el-tag size="small" type="success">仅在售房源</el-tag>
      </div>
    </template>
    <v-chart :option="option" autoresize style="height: 320px" />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { TrendCharts } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { ScatterChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  MarkLineComponent,
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { HouseItem } from '@/types'

use([ScatterChart, GridComponent, TooltipComponent, MarkLineComponent, CanvasRenderer])

const props = defineProps<{ houses: HouseItem[] }>()

const option = computed(() => {
  const onSale = props.houses.filter((h) => h.askpriceeachB > 0)
  const data = onSale.map((h) => [Number(h.floor), h.askpriceeachB, h.housenb, h.askpricetotalB])

  const prices = onSale.map((h) => h.askpriceeachB)
  const avgPrice = prices.length ? prices.reduce((a, b) => a + b, 0) / prices.length : 0

  return {
    tooltip: {
      trigger: 'item',
      formatter: (p: any) => {
        const [floor, price, room, total] = p.value
        return `
          <b>${room}</b><br/>
          楼层: ${floor}层<br/>
          单价: ${price.toLocaleString()} 元/㎡<br/>
          总价: ${(total / 10000).toFixed(2)} 万
        `
      },
    },
    grid: { left: 70, right: 20, top: 20, bottom: 40 },
    xAxis: {
      type: 'value',
      name: '楼层',
      nameLocation: 'end',
      nameTextStyle: { fontSize: 11 },
      minInterval: 1,
    },
    yAxis: {
      type: 'value',
      name: '单价(元/㎡)',
      nameTextStyle: { fontSize: 11 },
      axisLabel: {
        formatter: (v: number) => (v / 10000).toFixed(1) + 'w',
      },
    },
    series: [
      {
        type: 'scatter',
        data,
        symbolSize: 10,
        itemStyle: {
          color: (p: any) => {
            const price = p.value[1]
            if (price <= avgPrice * 0.97) return '#67C23A'
            if (price >= avgPrice * 1.03) return '#F56C6C'
            return '#409EFF'
          },
          opacity: 0.8,
        },
        markLine: {
          silent: true,
          lineStyle: { color: '#E6A23C', type: 'dashed' },
          data: [
            {
              yAxis: avgPrice,
              label: {
                formatter: `均价 ${Math.round(avgPrice).toLocaleString()}`,
                fontSize: 11,
              },
            },
          ],
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
