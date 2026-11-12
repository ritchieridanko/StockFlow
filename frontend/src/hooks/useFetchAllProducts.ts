import { useCallback, useEffect, useState } from "react";
import toast from "react-hot-toast";

import { productApi } from "../apis";
import { COLORS } from "../constants";
import type { APIResponse, Product, ProductListItem } from "../types";
import { isAPIErrResponse } from "../utils";

export const useFetchAllProducts = (
  page: number,
  size: number,
  fire: boolean,
) => {
  const [products, setProducts] = useState<ProductListItem[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);

  const fetchAllProducts = useCallback(async () => {
    setLoading(true);

    try {
      const data: APIResponse<ProductListItem[]> = await productApi.getAll({
        page: page,
        size: size,
      });

      setProducts(data.data);
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
        id: "fap-err-toast",
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
  }, [page, size]);

  useEffect(() => {
    if (fire) fetchAllProducts();
  }, [fetchAllProducts, fire]);

  const append = (product: Product) => {
    setProducts((prev) => [
      ...prev,
      {
        id: product.id,
        sku: product.sku,
        name: product.name,
        price: product.price,
        stock: product.stock,
        is_active: product.is_active,
        created_at: product.created_at,
        updated_at: product.updated_at,
      },
    ]);
  };

  const replace = (product: Product, last_movement_at?: string) => {
    setProducts((prev) =>
      prev.map((p) =>
        p.id === product.id
          ? {
              id: product.id,
              sku: product.sku,
              name: product.name,
              price: product.price,
              stock: product.stock,
              is_active: product.is_active,
              last_movement_at: last_movement_at,
              created_at: product.created_at,
              updated_at: product.updated_at,
            }
          : p,
      ),
    );
  };

  const remove = (id: number) => {
    setProducts((prev) => prev.filter((p) => p.id !== id));
  };

  return {
    products,
    total,
    loading,
    refetch: fetchAllProducts,
    append,
    replace,
    remove,
  };
};
