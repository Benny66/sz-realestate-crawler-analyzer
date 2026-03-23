<template>
  <el-card shadow="never">
    <template #header>
      <div class="card-header">
        <el-icon><List /></el-icon>
        <span>房源明细</span>
        <el-input
          v-model="search"
          placeholder="搜索房号/楼层/状态"
          clearable
          size="small"
          style="width: 200px; margin-left: auto"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="statusFilter" placeholder="状态筛选" clearable size="small" style="width: 130px">
          <el-option v-for="s in statusOptions" :key="s" :label="s" :value="s" />
        </el-select>
      </div>
    </template>

    <el-table
      :data="filteredHouses"
      stripe
      border
      size="small"
      :default-sort="{ prop: 'floor', order: 'descending' }"
      max-height="420"
    >
      <el-table-column prop="floor" label="楼层" width="70" sortable align="center">
        <template #default="{ row }">
          <el-tag :type="floorTagType(Number(row.floor))" size="small">{{ row.floor }}层</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="housenb" label="房号" width="80" align="center" />
      <el-table-column label="单价(元/㎡)" width="120" sortable :sort-method="sortByPrice" align="right">
        <template #default="{ row }">
          <span :class="priceClass(row.askpriceeachB)">
            {{ row.askpriceeachB > 0 ? row.askpriceeachB.toLocaleString() : '-' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="总价(万)" width="100" align="right" sortable :sort-method="sortByTotal">
        <template #default="{ row }">
          {{ row.askpricetotalB > 0 ? (row.askpricetotalB / 10000).toFixed(2) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="建筑面积(㎡)" width="110" align="right">
        <template #default="{ row }">{{ row.ysbuildingarea.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="套内面积(㎡)" width="110" align="right">
        <template #default="{ row }">{{ row.ysinsidearea.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="赠送面积(㎡)" width="110" align="right">
        <template #default="{ row }">
          <span class="expand-area">+{{ row.ysexpandarea.toFixed(2) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="备案单价" width="110" align="right">
        <template #default="{ row }">
          {{ row.recordedPricePerUnitInside > 0 ? row.recordedPricePerUnitInside.toLocaleString() : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="lastStatusName" label="状态" width="110" align="center">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.lastStatusName)" size="small">
            {{ row.lastStatusName }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>

    <div class="table-footer">
      共 {{ filteredHouses.length }} 条 / 总计 {{ houses.length }} 条
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { List, Search } from '@element-plus/icons-vue'
import type { HouseItem } from '@/types'

const props = defineProps<{ houses: HouseItem[]; avgPrice?: number }>()

const search = ref('')
const statusFilter = ref('')

const statusOptions = computed(() => {
  const set = new Set(props.houses.map((h) => h.lastStatusName))
  return [...set]
})

const filteredHouses = computed(() => {
  let list = props.houses
  if (statusFilter.value) {
    list = list.filter((h) => h.lastStatusName === statusFilter.value)
  }
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(
      (h) =>
        h.housenb.toLowerCase().includes(q) ||
        h.floor.includes(q) ||
        h.lastStatusName.includes(q)
    )
  }
  return list
})

function floorTagType(floor: number) {
  if (floor >= 20) return 'success'
  if (floor >= 10) return 'warning'
  return 'info'
}

function statusTagType(status: string) {
  if (status === '期房待售') return 'success'
  if (status.includes('合同') || status.includes('认购')) return 'danger'
  return 'info'
}

function priceClass(price: number) {
  if (!price) return ''
  const avg = props.avgPrice || 0
  if (!avg) return 'price-normal'
  if (price <= avg * 0.97) return 'price-low'
  if (price >= avg * 1.03) return 'price-high'
  return 'price-normal'
}

function sortByPrice(a: HouseItem, b: HouseItem) {
  return a.askpriceeachB - b.askpriceeachB
}
function sortByTotal(a: HouseItem, b: HouseItem) {
  return a.askpricetotalB - b.askpricetotalB
}
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}
.expand-area {
  color: var(--el-color-success);
}
.price-low  { color: var(--el-color-success); font-weight: 600; }
.price-high { color: var(--el-color-danger);  font-weight: 600; }
.price-normal { color: var(--el-text-color-primary); }
.table-footer {
  text-align: right;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 8px;
}
</style>
