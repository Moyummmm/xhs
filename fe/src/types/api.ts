// API 响应通用类型
export interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data: T;
}

// 分页参数
export interface PaginationParams {
  page?: number;
  page_size?: number;
}

// 分页响应
export interface PaginationResponse {
  total: number;
  page: number;
  page_size: number;
  has_more: boolean;
}

// 分页列表响应
export interface PaginatedList<T> {
  list: T[];
  pagination: PaginationResponse;
}
