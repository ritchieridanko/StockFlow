import { useCallback, useState } from "react";
import toast from "react-hot-toast";

import { productApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, Product, UpdateStockPayload } from "../types";
import { isAPIErrResponse } from "../utils";

export const useUpdateStock = () => {
  const [form, setForm] = useState<UpdateStockPayload>({
    type: "in",
    quantity: 0,
    note: "",
  });
  const [loading, setLoading] = useState<boolean>(false);

  const updateStock = useCallback(
    async (id: number) => {
      setLoading(true);

      try {
        const data: APIResponse<Product> = await productApi.updateStock(
          id,
          form,
        );

        toast.success(data.message, {
          id: "us-toast",
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
          id: "us-err-toast",
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
    },
    [form],
  );

  const clearForm = () => {
    setForm({ type: "in", quantity: 0, note: "" });
  };

  return { form, setForm, clearForm, loading, updateStock };
};
