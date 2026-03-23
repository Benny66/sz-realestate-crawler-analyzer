import axios from 'axios'
import type {  
  APIResponse,
  AnalyzeRequest,
  AnalyzeResponse,
  ProjectListResponse,
  BuildingListResponse,
  HistoryRecord,
  CompareRequest,
  CompareResponse,
  FavoriteItem,
} from '@/types'
import { ElMessage } from 'element-plus'

const http = axios.create({
  baseURL: '/api',
  timeout: 60000,
})

// 响应拦截器：统一处理业务错误
http.interceptors.response.use(
  (res) => {
    const data: APIResponse = res.data
    if (data.code !== 0) {
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message))
    }
    return res
  },
  (err) => {
    ElMessage.error(err.message || '网络错误，请检查后端服务是否启动')
    return Promise.reject(err)
  }
)

// 搜索楼盘列表
export function searchProjects(
  keyword: string,
  pageIndex = 1,
  pageSize = 12,
  zone = ''
): Promise<ProjectListResponse> {
  return http
    .get<APIResponse<ProjectListResponse>>('/search', {
      params: { keyword, pageIndex, pageSize, zone },
    })
    .then((r) => r.data.data)
}

// 获取楼栋列表
export function getBuildings(
  ysProjectId: number,
  preSellId: number
): Promise<BuildingListResponse> {
  return http
    .get<APIResponse<BuildingListResponse>>('/buildings', {
      params: { ysProjectId, preSellId },
    })
    .then((r) => r.data.data)
}

// 执行分析
export function analyzeProject(req: AnalyzeRequest): Promise<AnalyzeResponse> {
  return http
    .post<APIResponse<AnalyzeResponse>>('/analyze', req)
    .then((r) => r.data.data)
}

// 多楼盘对比
export function compareProjects(req: CompareRequest): Promise<CompareResponse> {
  return http
    .post<APIResponse<CompareResponse>>('/compare', req)
    .then((r) => r.data.data)
}

// 获取历史记录
export function getHistory(projectName = ''): Promise<HistoryRecord[]> {
  return http
    .get<APIResponse<HistoryRecord[]>>('/history', {
      params: projectName ? { projectName } : {},
    })
    .then((r) => r.data.data)
}

// 删除历史记录
export function deleteHistory(id: string): Promise<void> {
  return http
    .delete<APIResponse<null>>('/history', { params: { id } })
    .then(() => undefined)
}

// 获取收藏列表
export function getFavorites(): Promise<FavoriteItem[]> {
  return http
    .get<APIResponse<FavoriteItem[]>>('/favorites')
    .then((r) => r.data.data)
}

// 新增收藏
export function addFavorite(item: FavoriteItem): Promise<void> {
  return http.post<APIResponse<null>>('/favorites', item).then(() => undefined)
}

// 删除收藏
export function deleteFavorite(id: string): Promise<void> {
  return http.delete<APIResponse<null>>('/favorites', { params: { id } }).then(() => undefined)
}

// 更新收藏推送设置
export function updateFavorite(id: string, data: { enablePush?: boolean; priceAlert?: boolean; saleAlert?: boolean }): Promise<void> {
  return http.put<APIResponse<null>>(`/favorites/update?id=${id}`, data).then(() => undefined)
}

// 导出 CSV（直接打开下载链接）
export function buildExportCSVUrl(req: AnalyzeRequest): string {
  const params = new URLSearchParams()
  if (req.keyword) params.set('keyword', req.keyword)
  if (req.buildingName) params.set('buildingName', req.buildingName)
  if (req.houseType) params.set('houseType', req.houseType)
  if (req.zone) params.set('zone', req.zone)
  if (req.ysProjectId) params.set('ysProjectId', String(req.ysProjectId))
  if (req.preSellId) params.set('preSellId', String(req.preSellId))
  if (req.fybId) params.set('fybId', String(req.fybId))
  return `/api/export/csv?${params.toString()}`
}

// 清除缓存
export function clearCache(): Promise<void> {
  return http.post<APIResponse<null>>('/cache/clear').then(() => undefined)
}
