import { useCallback, useEffect, useState } from "react";
import toast from "react-hot-toast";

import { orderApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, Order } from "../types";
import { isAPIErrResponse } from "../utils";

export const useFetchAllOrders = (
  page: number,
  size: number,
  search: string,
  fire: boolean,
) => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchAllOrders = useCallback(async () => {
    setLoading(true);

    try {
      const data: APIResponse<Order[]> = await orderApi.getAll({
        page: page,
        size: size,
        search: search,
      });

      setOrders(data.data);
      setTotal(data.pagination?.total ?? 0);
    } catch (err: unknown) {
      let message: string = "Something went wrong";

      if (isAPIErrResponse(err)) {
        message = err.error.message;
      } else if (err instanceof Error) {
        message = err.message;
      } else {
        console.error(err);
      }

      toast.error(message, {
        id: "fao-err-toast",
        duration: 5000,
        style: {
          background: COLORS.ERROR,
          color: COLORS.LIGHT,
        },
        iconTheme: {
          primary: COLORS.LIGHT,
          secondary: COLORS.ERROR,
        },
      });
    } finally {
      setLoading(false);
    }
  }, [page, size, search]);

  useEffect(() => {
    if (fire) fetchAllOrders();
  }, [fetchAllOrders, fire]);

  const append = (order: Order) => {
    setOrders((prev) => [order, ...prev]);
  };

  const replace = (order: Order) => {
    setOrders((prev) => prev.map((o) => (o.id === order.id ? order : o)));
  };

  return {
    orders,
    total,
    loading,
    refetch: fetchAllOrders,
    append,
    replace,
  };
};
