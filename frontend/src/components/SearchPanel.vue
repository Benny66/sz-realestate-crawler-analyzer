<template>
  <el-card class="search-panel" shadow="never">
    <template #header>
      <div class="panel-header">
        <el-icon><Search /></el-icon>
        <span>智能搜索</span>
      </div>
    </template>

    <el-form :model="form" label-width="80px" @submit.prevent>
      <!-- 楼盘关键字搜索 -->
      <el-form-item label="楼盘名称">
        <el-autocomplete
          v-model="form.keyword"
          :fetch-suggestions="fetchSuggestions"
          placeholder="输入楼盘关键字，如：乐宸"
          clearable
          style="width: 100%"
          @select="onProjectSelect"
          @clear="resetProject"
        >
          <template #default="{ item }">
            <div class="suggestion-item">
              <span class="project-name">{{ item.project }}</span>
              <el-tag size="small" type="info">{{ item.zone }}</el-tag>
            </div>
          </template>
        </el-autocomplete>
      </el-form-item>

      <!-- 区域选择 -->
      <el-form-item label="区域">
        <el-select v-model="form.zone" placeholder="不限区域" clearable style="width: 100%">
          <el-option label="不限" value="" />
          <el-option label="福田" value="福田" />
          <el-option label="罗湖" value="罗湖" />
          <el-option label="南山" value="南山" />
          <el-option label="宝安" value="宝安" />
          <el-option label="龙岗" value="龙岗" />
          <el-option label="龙华" value="龙华" />
          <el-option label="坪山" value="坪山" />
          <el-option label="光明" value="光明" />
          <el-option label="盐田" value="盐田" />
          <el-option label="大鹏" value="大鹏" />
        </el-select>
      </el-form-item>

      <!-- 楼栋选择 -->
      <el-form-item label="楼栋">
        <el-select
          v-model="form.buildingName"
          placeholder="请先搜索楼盘"
          :disabled="!buildings.length"
          clearable
          style="width: 100%"
        >
          <el-option
            v-for="b in buildings"
            :key="b.key"
            :label="b.label"
            :value="b.label"
          />
        </el-select>
      </el-form-item>

      <!-- 户型选择 -->
      <el-form-item label="户型">
        <el-select v-model="form.houseType" placeholder="不限户型" clearable style="width: 100%">
          <el-option label="不限" value="" />
          <el-option label="一房" value="一房" />
          <el-option label="两房" value="两房" />
          <el-option label="三房" value="三房" />
          <el-option label="四房" value="四房" />
          <el-option label="五房" value="五房" />
        </el-select>
      </el-form-item>

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button
          type="primary"
          :loading="loading"
          :disabled="!form.keyword"
          style="width: 100%"
          @click="onAnalyze"
        >
          <el-icon><DataAnalysis /></el-icon>
          开始分析
        </el-button>
      </el-form-item>

      <el-form-item>
        <el-button style="width: 100%" @click="onAddToCompare" :disabled="!form.keyword">
          <el-icon><Plus /></el-icon>
          加入对比
        </el-button>
      </el-form-item>
    </el-form>

    <!-- 对比列表预览 -->
    <div v-if="compareList.length" class="compare-preview">
      <div class="compare-title">
        <span>对比列表 ({{ compareList.length }}/5)</span>
        <el-button link type="danger" @click="$emit('clearCompare')">清空</el-button>
      </div>
      <el-tag
        v-for="(item, idx) in compareList"
        :key="idx"
        closable
        class="compare-tag"
        @close="$emit('removeCompare', idx)"
      >
        {{ item.keyword }} {{ item.buildingName }} {{ item.houseType }}
      </el-tag>
      <el-button
        type="warning"
        size="small"
        style="width: 100%; margin-top: 8px"
        :loading="compareLoading"
        @click="$emit('startCompare')"
      >
        <el-icon><ScaleToOriginal /></el-icon>
        开始对比
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Search, DataAnalysis, Plus, ScaleToOriginal } from '@element-plus/icons-vue'
import { searchProjects, getBuildings } from '@/api'
import type { ProjectItem, BuildingNameItem, AnalyzeRequest } from '@/types'

const props = defineProps<{
  loading: boolean
  compareLoading: boolean
  compareList: AnalyzeRequest[]
}>()

const emit = defineEmits<{
  (e: 'analyze', req: AnalyzeRequest): void
  (e: 'addToCompare', req: AnalyzeRequest): void
  (e: 'removeCompare', idx: number): void
  (e: 'clearCompare'): void
  (e: 'startCompare'): void
}>()

const form = reactive({
  keyword: '',
  zone: '',
  buildingName: '',
  houseType: '',
})

const buildings = ref<BuildingNameItem[]>([])
const selectedProject = ref<ProjectItem | null>(null)

// 自动补全搜索建议
async function fetchSuggestions(
  query: string,
  cb: (results: (ProjectItem & { value: string })[]) => void
) {
  if (!query || query.length < 1) {
    cb([])
    return
  }
  try {
    const res = await searchProjects(query, 1, 12, form.zone)
    cb(res.projects.map((p) => ({ ...p, value: p.project })))
  } catch {
    cb([])
  }
}

// 选中楼盘后自动加载楼栋列表
async function onProjectSelect(item: ProjectItem) {
  selectedProject.value = item
  form.keyword = item.project
  form.buildingName = ''
  buildings.value = []
  try {
    const res = await getBuildings(Number(item.sypId), Number(item.id))
    buildings.value = res.buildings
    if (res.buildings.length === 1) {
      form.buildingName = res.buildings[0].label
    }
  } catch {
    // 忽略
  }
}

function resetProject() {
  selectedProject.value = null
  buildings.value = []
  form.buildingName = ''
}

function buildRequest(): AnalyzeRequest {
  return {
    keyword: form.keyword,
    buildingName: form.buildingName,
    houseType: form.houseType,
    zone: form.zone,
    ysProjectId: selectedProject.value ? Number(selectedProject.value.sypId) : 0,
    preSellId: selectedProject.value ? Number(selectedProject.value.id) : 0,
  }
}

function onAnalyze() {
  if (!form.keyword) return
  emit('analyze', buildRequest())
}

// 当分析完成后，如果后端自动选择了楼栋，更新前端显示
defineExpose({
  updateBuildingName(name: string) {
    form.buildingName = name
  }
})

function onAddToCompare() {
  if (!form.keyword) return
  emit('addToCompare', buildRequest())
}
</script>

<style scoped>
.search-panel {
  height: 100%;
}
.panel-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}
.suggestion-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.project-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.compare-preview {
  border-top: 1px solid var(--el-border-color-lighter);
  padding-top: 12px;
  margin-top: 4px;
}
.compare-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  margin-bottom: 8px;
}
.compare-tag {
  margin: 3px;
}
</style>
