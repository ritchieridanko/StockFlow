import { useCallback, useState } from "react";
import toast from "react-hot-toast";

import { productApi } from "../apis";
import { COLORS } from "../constants";
import { isAPIErrResponse } from "../utils";

export const useDeleteProduct = () => {
  const [loading, setLoading] = useState<boolean>(false);

  const deleteProduct = useCallback(async (id: number): Promise<boolean> => {
    setLoading(true);

    try {
      await productApi.delete(id);

      toast.success("Product deleted successfully", {
        id: "dp-toast",
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

      return true;
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
        id: "dp-err-toast",
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

      return false;
    } finally {
      setLoading(false);
    }
  }, []);

  return { loading, deleteProduct };
};
