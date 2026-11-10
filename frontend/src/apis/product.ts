import { client } from "./client";
import type {
  APIResponse,
  CreateProductPayload,
  GetAllProductsParams,
  Product,
  ProductListItem,
  StockMovement,
  UpdateProductPayload,
  UpdateStockPayload,
} from "../types";

export const productApi = {
  create: async (
    payload: CreateProductPayload,
  ): Promise<APIResponse<Product>> => {
    const res = await client.post<APIResponse<Product>>("/products", payload);
    return res.data;
  },

  getAll: async (
    params: GetAllProductsParams,
  ): Promise<APIResponse<ProductListItem[]>> => {
    const res = await client.get<APIResponse<ProductListItem[]>>(
      `/products?page=${params.page}&limit=${params.limit}${params.search && `&search=${params.search}`}`,
    );
    return res.data;
  },

  getAllStockMovements: async (
    id: number,
  ): Promise<APIResponse<StockMovement[]>> => {
    const res = await client.get<APIResponse<StockMovement[]>>(
      `/products/${id}/movements`,
    );
    return res.data;
  },

  update: async (
    id: number,
    payload: UpdateProductPayload,
  ): Promise<APIResponse<Product>> => {
    const res = await client.put<APIResponse<Product>>(
      `/products/${id}`,
      payload,
    );
    return res.data;
  },

  updateStock: async (
    id: number,
    payload: UpdateStockPayload,
  ): Promise<APIResponse<Product>> => {
    const res = await client.post<APIResponse<Product>>(
      `/products/${id}/stock`,
      payload,
    );
    return res.data;
  },

  delete: async (id: number): Promise<void> => {
    await client.delete(`/products/${id}`);
  },
};
