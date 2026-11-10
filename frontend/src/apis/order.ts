import { client } from "./client";
import type { APIResponse, CreateOrderPayload, Order } from "../types";

export const orderApi = {
  create: async (payload: CreateOrderPayload): Promise<APIResponse<Order>> => {
    const res = await client.post<APIResponse<Order>>("/orders", payload);
    return res.data;
  },

  getById: async (id: number): Promise<APIResponse<Order>> => {
    const res = await client.get<APIResponse<Order>>(`/orders/${id}`);
    return res.data;
  },

  cancel: async (id: number): Promise<APIResponse<Order>> => {
    const res = await client.post<APIResponse<Order>>(`/orders/${id}/cancel`);
    return res.data;
  },
};
