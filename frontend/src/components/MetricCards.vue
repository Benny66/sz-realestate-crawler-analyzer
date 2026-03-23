<template>
  <div class="metric-cards">
    <!-- 价格指标 -->
    <el-card class="metric-card price-card" shadow="hover">
      <div class="card-title">
        <el-icon color="#409EFF"><Money /></el-icon>
        价格指标
      </div>
      <div class="metric-grid">
        <div class="metric-item">
          <div class="metric-label">最低单价</div>
          <div class="metric-value primary">{{ fmt(analysis.min_unit_price) }}<span class="unit">元/㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">最高单价</div>
          <div class="metric-value danger">{{ fmt(analysis.max_unit_price) }}<span class="unit">元/㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">单价中位数</div>
          <div class="metric-value warning">{{ fmt(analysis.median_unit_price) }}<span class="unit">元/㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">楼层价格梯度</div>
          <div class="metric-value">{{ analysis.floor_price_premium.toFixed(2) }}<span class="unit">%</span></div>
        </div>
        <div class="metric-item wide">
          <div class="metric-label">总价区间</div>
          <div class="metric-value">
            {{ fmtW(analysis.min_total_price) }} ~ {{ fmtW(analysis.max_total_price) }}
            <span class="unit">万</span>
          </div>
        </div>
        <div class="metric-item wide">
          <div class="metric-label">总价均值</div>
          <div class="metric-value success">{{ fmtW(analysis.avg_total_price) }}<span class="unit">万</span></div>
        </div>
      </div>
    </el-card>

    <!-- 面积指标 -->
    <el-card class="metric-card area-card" shadow="hover">
      <div class="card-title">
        <el-icon color="#67C23A"><Grid /></el-icon>
        面积指标
      </div>
      <div class="metric-grid">
        <div class="metric-item">
          <div class="metric-label">产权建筑面积</div>
          <div class="metric-value primary">{{ analysis.building_area.toFixed(2) }}<span class="unit">㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">套内面积</div>
          <div class="metric-value">{{ analysis.inside_area.toFixed(2) }}<span class="unit">㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">赠送面积</div>
          <div class="metric-value success">+{{ analysis.expand_area.toFixed(2) }}<span class="unit">㎡</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">实际使用面积</div>
          <div class="metric-value warning">{{ analysis.actual_use_area.toFixed(2) }}<span class="unit">㎡</span></div>
        </div>
        <div class="metric-item wide">
          <div class="metric-label">得房率</div>
          <div class="housing-rate-bar">
            <el-progress
              :percentage="analysis.housing_rate"
              :color="housingRateColor"
              :stroke-width="14"
              striped
              striped-flow
            />
          </div>
        </div>
        <div class="metric-item wide">
          <div class="metric-label">实际使用成本</div>
          <div class="metric-value">{{ fmt(analysis.cost_per_actual_area) }}<span class="unit">元/㎡</span></div>
        </div>
      </div>
    </el-card>

    <!-- 供应指标 -->
    <el-card class="metric-card supply-card" shadow="hover">
      <div class="card-title">
        <el-icon color="#E6A23C"><OfficeBuilding /></el-icon>
        供应指标
      </div>
      <div class="metric-grid">
        <div class="metric-item">
          <div class="metric-label">在售套数</div>
          <div class="metric-value primary">{{ analysis.total_count }}<span class="unit">套</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">在售楼层</div>
          <div class="metric-value">{{ analysis.min_floor }} ~ {{ analysis.max_floor }}<span class="unit">层</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">低楼层(2-9)</div>
          <div class="metric-value">{{ analysis.low_floor_count }}<span class="unit">套</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">中楼层(10-19)</div>
          <div class="metric-value">{{ analysis.mid_floor_count }}<span class="unit">套</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">高楼层(20+)</div>
          <div class="metric-value success">{{ analysis.high_floor_count }}<span class="unit">套</span></div>
        </div>
        <div class="metric-item">
          <div class="metric-label">单层最多/最少</div>
          <div class="metric-value">{{ analysis.max_per_floor }} / {{ analysis.min_per_floor }}<span class="unit">套</span></div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Money, Grid, OfficeBuilding } from '@element-plus/icons-vue'
import type { AnalysisResult } from '@/types'

const props = defineProps<{ analysis: AnalysisResult }>()

function fmt(val: number): string {
  return val > 0 ? val.toLocaleString('zh-CN', { maximumFractionDigits: 0 }) : '-'
}
function fmtW(val: number): string {
  return val > 0 ? (val / 10000).toFixed(2) : '-'
}

const housingRateColor = computed(() => {
  const r = props.analysis.housing_rate
  if (r >= 80) return '#67C23A'
  if (r >= 70) return '#E6A23C'
  return '#F56C6C'
})
</script>

<style scoped>
.metric-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}
.metric-card {
  min-height: 200px;
}
.card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 16px;
}
.metric-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 8px;
}
.metric-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.metric-item.wide {
  grid-column: span 2;
}
.metric-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.metric-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--el-text-color-primary);
  line-height: 1.2;
}
.metric-value.primary { color: var(--el-color-primary); }
.metric-value.success { color: var(--el-color-success); }
.metric-value.warning { color: var(--el-color-warning); }
.metric-value.danger  { color: var(--el-color-danger); }
.unit {
  font-size: 12px;
  font-weight: 400;
  color: var(--el-text-color-secondary);
  margin-left: 2px;
}
.housing-rate-bar {
  padding-top: 4px;
}
@media (max-width: 1200px) {
  .metric-cards {
    grid-template-columns: 1fr;
  }
}
</style>
