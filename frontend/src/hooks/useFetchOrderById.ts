import { useCallback, useEffect, useState } from "react";
import toast from "react-hot-toast";

import { orderApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, Order } from "../types";
import { isAPIErrResponse } from "../utils";

export const useFetchOrderById = (id: number) => {
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchOrderById = useCallback(async (id: number) => {
    setLoading(true);

    try {
      const data: APIResponse<Order> = await orderApi.getById(id);
      setOrder(data.data);
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
        id: "foid-err-toast",
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
  }, []);

  useEffect(() => {
    if (id) fetchOrderById(id);
  }, [fetchOrderById, id]);

  return { order, loading };
};
