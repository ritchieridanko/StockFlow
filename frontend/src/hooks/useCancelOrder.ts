import { useCallback, useState } from "react";
import toast from "react-hot-toast";

import { orderApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, Order } from "../types";
import { isAPIErrResponse } from "../utils";

export const useCancelOrder = () => {
  const [loading, setLoading] = useState<boolean>(false);

  const cancelOrder = useCallback(async (id: number) => {
    setLoading(true);

    try {
      const data: APIResponse<Order> = await orderApi.cancel(id);

      toast.success(data.message, {
        id: "clo-toast",
        duration: 5000,
        style: {
          background: COLORS.SUCCESS,
          color: COLORS.LIGHT,
        },
        iconTheme: {
          primary: COLORS.LIGHT,
          secondary: COLORS.SUCCESS,
        },
      });

      return data.data;
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
        id: "clo-err-toast",
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

  return { loading, cancelOrder };
};
