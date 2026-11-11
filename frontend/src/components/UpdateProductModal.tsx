import React, { useEffect } from "react";
import toast from "react-hot-toast";

import { COLORS } from "../constants";
import { useUpdateProduct } from "../hooks";
import type { Product, ProductListItem } from "../types";
import { TextField } from "./ui";

interface UpdateProductModalProps {
  isVisible: boolean;
  product: ProductListItem | null;
  onClickUpdate: (value: Product) => void;
  onClickCancel: () => void;
}

const UpdateProductModal: React.FC<UpdateProductModalProps> = ({
  isVisible,
  product,
  onClickUpdate,
  onClickCancel,
}: UpdateProductModalProps) => {
  if (!isVisible) return null;

  const { form, setForm, clearForm, updateProduct } = useUpdateProduct();

  useEffect(() => {
    if (product) {
      setForm({
        name: product.name,
        sku: product.sku,
        price: product.price,
      });
    }
  }, [product]);

  const handleUpdate = async () => {
    if (!product) return;
    if (form.name === "" || form.sku === "") {
      toast.error("All fields are required", {
        id: "upm-err-toast",
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
      return;
    }

    const res: Product | undefined = await updateProduct(product.id);
    if (res) {
      onClickUpdate(res);
      clearForm();
      return;
    }
  };

  return (
    <div
      className="fixed z-99 inset-0 flex justify-center items-center bg-black/50 backdrop-blur-xs cursor-pointer"
      onClick={onClickCancel}
    >
      <div
        className="p-5 w-full max-w-md flex flex-col justify-center gap-4 bg-white rounded-lg cursor-default shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="font-sans font-medium text-gray-700 text-xl">
          Update Product
        </p>
        <div className="flex flex-col justify-center items-center gap-4">
          <TextField
            type="text"
            placeholder="Enter the product name"
            value={form.name ?? ""}
            onChange={(v) => setForm((form) => ({ ...form, name: v }))}
          />
          <TextField
            type="text"
            placeholder="Enter the product sku"
            value={form.sku ?? ""}
            onChange={(v) => setForm((form) => ({ ...form, sku: v }))}
          />
          <TextField
            type="number"
            placeholder="Enter the product price"
            value={String(form.price ?? 0)}
            onChange={(v) => setForm((form) => ({ ...form, price: Number(v) }))}
            min={0}
            step={1}
          />
        </div>
        <div className="flex flex-row justify-end items-center gap-4">
          <button
            className="px-4 py-2 bg-gray-200 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-gray-300"
            onClick={onClickCancel}
          >
            <p className="font-sans font-medium text-gray-700 text-base">
              Cancel
            </p>
          </button>
          <button
            className="px-4 py-2 bg-orange-500 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-orange-600"
            onClick={handleUpdate}
          >
            <p className="font-sans font-medium text-white text-base">
              Update Product
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default UpdateProductModal;
