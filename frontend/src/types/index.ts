// ---- 基础数据结构 ----
  
export interface HouseItem {
  housenb: string
  floor: string
  ysbuildingarea: number
  ysexpandarea: number
  ysinsidearea: number
  askpricetotalB: number
  askpriceeachB: number
  recordedPricePerUnitInside: number
  lastStatusName: string
  buildingName: string
  useage: string
  barrierFree: string
}

export interface FloorGroup {
  floor: string
  list: HouseItem[]
}

export interface BuildingNameItem {
  label: string
  value: string
  key: string
  acive: boolean
}

export interface ProjectItem {
  id: string
  sypId: string
  zone: string
  strpreprojectid: string
  project: string
  name: string
  siteaddress: string
  passdate: string
  imagePath: string
}

// ---- 分析结果 ----

export interface AnalysisResult {
  min_unit_price: number
  max_unit_price: number
  median_unit_price: number
  min_recorded_unit_price: number
  max_recorded_unit_price: number
  min_total_price: number
  max_total_price: number
  avg_total_price: number
  total_price_span: number
  building_area: number
  expand_area: number
  inside_area: number
  actual_use_area: number
  housing_rate: number
  total_count: number
  min_floor: number
  max_floor: number
  low_floor_count: number
  mid_floor_count: number
  high_floor_count: number
  max_per_floor: number
  min_per_floor: number
  unit_price_per_expand_area: number
  cost_per_actual_area: number
  floor_price_premium: number
}

export interface SaleStatusStats {
  statusName: string
  count: number
}

export interface SaleSummary {
  total_count: number
  for_sale_count: number
  sold_count: number
  status_details: SaleStatusStats[]
  sold_rate: number
  for_sale_rate: number
}

export interface ResolvedProjectParams {
  ysProjectId: number
  preSellId: number
  fybId: number
  projectName: string
  buildingName: string  // 实际使用的楼栋名称
  autoSelected: boolean  // 是否自动选择了楼栋
}

// ---- API 请求/响应 ----

export interface AnalyzeRequest {
  keyword: string
  buildingName: string
  houseType: string
  zone?: string
  ysProjectId?: number
  preSellId?: number
  fybId?: number
}

export interface PriceAlert {
  level: 'info' | 'warning' | 'danger' | 'success'
  title: string
  message: string
  direction: 'up' | 'down' | 'flat'
  changePct: number
}

export interface AnalyzeResponse {
  analysis: AnalysisResult
  sale: SaleSummary
  houses: HouseItem[]
  floorGroups: FloorGroup[]
  params: ResolvedProjectParams
  alerts: PriceAlert[]
  updatedAt: string
}

export interface ProjectListResponse {
  total: number
  projects: ProjectItem[]
}

export interface BuildingListResponse {
  buildings: BuildingNameItem[]
}

export interface APIResponse<T = unknown> {
  code: number
  message: string
  data: T
}

// ---- 历史记录 ----

export interface HistoryRecord {
  id: string
  projectName: string
  buildingName: string
  houseType: string
  zone?: string
  analysis: AnalysisResult
  sale: SaleSummary
  createdAt: string
}

// ---- 多楼盘对比 ----

export interface CompareItem {
  projectName: string
  buildingName: string
  houseType: string
  zone?: string
  analysis: AnalysisResult
  sale: SaleSummary
}

export interface CompareResponse {
  items: CompareItem[]
  updatedAt: string
}

export interface CompareRequest {
  items: AnalyzeRequest[]
}

export interface FavoriteItem {
  id?: string
  projectName: string
  buildingName: string
  houseType: string
  zone?: string
  enablePush?: boolean
  priceAlert?: boolean
  saleAlert?: boolean
  createdAt?: string
}
