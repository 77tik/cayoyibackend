// API 响应类型定义
export interface FluidCondition {
  id: number;
  name: string;
}

export interface FluidResult {
  // 根据实际接口返回补充字段
  effective_head: number;
  active_power: number;
  // 其他结果字段
}

export interface StructuralCondition {
  id: number;
  name: string;
  // 根据实际接口返回补充其他字段
}

export interface StructuralResult {
  // 根据实际接口返回补充字段
  effective_head: number;
  active_power: number;
  // 其他结果字段
}

export interface StrainMonitoring {
  // 根据实际接口返回补充字段
  timestamp: string;
  strain_value: number;
  // 其他监测数据字段
}
