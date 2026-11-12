import { useCallback, useEffect, useState } from "react";
import toast from "react-hot-toast";

import { productApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, StockMovement } from "../types";
import { isAPIErrResponse } from "../utils";

export const useFetchAllStockMovements = (id: number) => {
  const [stockMovements, setStockMovements] = useState<StockMovement[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchAllStockMovements = useCallback(async (id: number) => {
    setLoading(true);

    try {
      const data: APIResponse<StockMovement[]> =
        await productApi.getAllStockMovements(id);
      setStockMovements(data.data);
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
        id: "fasm-err-toast",
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
    if (id > 0) fetchAllStockMovements(id);
  }, [fetchAllStockMovements, id]);

  const append = (sm: StockMovement) => {
    setStockMovements((prev) => [sm, ...prev]);
  };

  return { stockMovements, loading, append };
};
