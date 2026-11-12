import React from "react";
import toast from "react-hot-toast";

import { COLORS } from "../constants";
import { useCreateProduct } from "../hooks";
import type { Product } from "../types";
import { TextField } from "./ui";

interface CreateProductModalProps {
  isVisible: boolean;
  onClickCreate: (value: Product) => void;
  onClickCancel: () => void;
}

const CreateProductModal: React.FC<CreateProductModalProps> = ({
  isVisible,
  onClickCreate,
  onClickCancel,
}: CreateProductModalProps) => {
  const { form, setForm, clearForm, createProduct } = useCreateProduct();

  const handleCreate = async () => {
    if (form.name === "" || form.sku === "") {
      toast.error("All fields are required", {
        id: "cpm-err-toast",
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

    const res: Product | undefined = await createProduct();
    if (res) {
      onClickCreate(res);
      clearForm();
    }
  };

  if (!isVisible) return null;

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
          Create a New Product
        </p>
        <div className="flex flex-col justify-center items-center gap-4">
          <TextField
            type="text"
            placeholder="Enter the product name"
            value={form.name}
            onChange={(v) => setForm((form) => ({ ...form, name: v }))}
          />
          <TextField
            type="text"
            placeholder="Enter the product sku"
            value={form.sku}
            onChange={(v) => setForm((form) => ({ ...form, sku: v }))}
          />
          <TextField
            type="number"
            placeholder="Enter the product price"
            value={String(form.price)}
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
            className="px-4 py-2 bg-blue-500 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-600"
            onClick={handleCreate}
          >
            <p className="font-sans font-medium text-white text-base">
              Create Product
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default CreateProductModal;
