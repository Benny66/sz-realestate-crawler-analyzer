<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><PieChart /></el-icon>
        <span>销售状态分布</span>
        <el-tag type="info" size="small">共 {{ sale.total_count }} 套</el-tag>
      </div>
    </template>
    <div class="chart-wrap">
      <v-chart :option="option" autoresize style="height: 280px" />
    </div>
    <!-- 状态明细表 -->
    <el-table :data="sale.status_details" size="small" style="margin-top: 8px">
      <el-table-column prop="statusName" label="状态" />
      <el-table-column prop="count" label="套数" width="70" align="right" />
      <el-table-column label="占比" width="100" align="right">
        <template #default="{ row }">
          {{ ((row.count / sale.total_count) * 100).toFixed(1) }}%
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { PieChart } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { PieChart as EPieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { SaleSummary } from '@/types'

use([EPieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const props = defineProps<{ sale: SaleSummary }>()

const COLORS = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9B59B6']

const option = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} 套 ({d}%)',
  },
  legend: {
    orient: 'vertical',
    right: 10,
    top: 'center',
    textStyle: { fontSize: 12 },
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['38%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      emphasis: {
        label: { show: true, fontSize: 14, fontWeight: 'bold' },
      },
      data: (props.sale.status_details || []).map((s, i) => ({
        name: s.statusName,
        value: s.count,
        itemStyle: { color: COLORS[i % COLORS.length] },
      })),
    },
  ],
}))
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.chart-wrap {
  width: 100%;
}
</style>
