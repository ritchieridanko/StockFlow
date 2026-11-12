import { useCallback, useState } from "react";
import toast from "react-hot-toast";

import { productApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, CreateProductPayload, Product } from "../types";
import { isAPIErrResponse } from "../utils";

export const useCreateProduct = () => {
  const [form, setForm] = useState<CreateProductPayload>({
    name: "",
    sku: "",
    price: 0,
  });
  const [loading, setLoading] = useState<boolean>(false);

  const createProduct = useCallback(async () => {
    setLoading(true);

    try {
      const data: APIResponse<Product> = await productApi.create(form);

      toast.success(data.message, {
        id: "cp-toast",
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
        id: "cp-err-toast",
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
  }, [form]);

  const clearForm = () => {
    setForm({ sku: "", name: "", price: 0 });
  };

  return { form, setForm, clearForm, loading, createProduct };
};
