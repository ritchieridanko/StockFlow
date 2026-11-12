import type { StockMovementType } from "./stock_movement";

interface Metadata {
  request_id: string;
  timestamp: string;
}

interface Pagination {
  page: number;
  size: number;
  total: number;
}

export interface Order {
  id: number;
  customer: string;
  status: string;
  total: number;
  items: OrderItem[];
  created_at: string;
  updated_at: string;
  cancelled_at?: string;
}

export interface OrderItem {
  id: number;
  order_id: number;
  sku: string;
  quantity: number;
  unit_price: number;
  created_at: string;
  updated_at: string;
}

export interface Product {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductListItem {
  id: number;
  sku: string;
  name: string;
  price: number;
  stock: number;
  is_active: boolean;
  last_movement_at?: string;
  created_at: string;
  updated_at: string;
}

export interface StockMovement {
  id: number;
  product_id: number;
  type: StockMovementType;
  quantity: number;
  note?: string;
  created_at: string;
  updated_at: string;
}

export interface APIResponse<T> {
  status: number;
  message: string;
  data: T;
  metadata: Metadata;
  pagination?: Pagination;
}

export interface APIErrResponse {
  status: number;
  error: {
    code: string;
    message: string;
  };
  metadata: Metadata;
}

interface OrderItemPayload {
  sku: string;
  quantity: number;
}

export interface CreateOrderPayload {
  customer: string;
  items: OrderItemPayload[];
}

export interface CreateProductPayload {
  sku: string;
  name: string;
  price: number;
}

export interface GetAllProductsParams {
  search?: string;
  page: number;
  size: number;
}

export interface UpdateProductPayload {
  sku?: string;
  name?: string;
  price?: number;
}

export interface UpdateStockPayload {
  type: StockMovementType | "";
  quantity: number;
  note?: string;
}
