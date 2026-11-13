import React, { useEffect, useState } from "react";
import toast from "react-hot-toast";

import { Minus, Plus } from "../assets/icons";
import { COLORS } from "../constants";
import type { ProductListItem } from "../types";

interface OrderProductModalProps {
  isVisible: boolean;
  product: ProductListItem | null;
  onClickConfirm: (product: ProductListItem, quantity: number) => void;
  onClickCancel: () => void;
}

const OrderProductModal: React.FC<OrderProductModalProps> = ({
  isVisible,
  product,
  onClickConfirm,
  onClickCancel,
}: OrderProductModalProps) => {
  const [quantity, setQuantity] = useState<number>(0);

  useEffect(() => {
    if (product) setQuantity(0);
  }, [product]);

  const handleConfirm = async () => {
    if (!product) return;
    if (quantity <= 0 || quantity > product.stock) {
      toast.error("Invalid quantity", {
        id: "opm-err-toast",
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

    onClickConfirm(product, quantity);
  };

  if (!isVisible) return null;

  const isSubClickable: boolean = !!product && quantity > 0;
  const isAddClickable: boolean = !!product && quantity < product.stock;

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
          Add to Order List
        </p>
        <div className="flex flex-col justify-center items-start">
          <p className="font-sans font-normal text-gray-600 text-base">
            Product:{" "}
            <span className="font-sans font-medium text-gray-700 text-base">
              {product?.name}
            </span>
          </p>
          <p className="font-sans font-normal text-gray-600 text-base">
            SKU:{" "}
            <span className="font-sans font-medium text-gray-700 text-base">
              {product?.sku}
            </span>
          </p>
          <p className="font-sans font-normal text-gray-600 text-base">
            Stock:{" "}
            <span className="font-sans font-medium text-gray-700 text-base">
              {product?.stock}
            </span>
          </p>
        </div>
        <div className="my-2 flex flex-row justify-center items-center gap-4">
          <button
            className={`
              px-2 py-1 rounded-sm transition-all duration-200 ease-in-out
              ${isSubClickable ? "bg-blue-100 cursor-pointer hover:bg-blue-200" : "bg-gray-200"}
            `}
            onClick={() => {
              if (isSubClickable) setQuantity((prev) => prev - 1);
            }}
          >
            <span
              className={isSubClickable ? "text-blue-500" : "text-gray-700"}
            >
              <Minus size="size-5" strokeWidth={2} />
            </span>
          </button>
          <p
            className={`
              px-5 py-2 flex justify-center items-center rounded-sm font-sans font-medium text-base select-none
              ${quantity === 0 ? "bg-gray-200 text-gray-700" : "bg-blue-100 text-blue-500"}
            `}
          >
            {quantity}
          </p>
          <button
            className={`
              px-2 py-1 rounded-sm transition-all duration-200 ease-in-out
              ${isAddClickable ? "bg-blue-100 cursor-pointer hover:bg-blue-200" : "bg-gray-200"}
            `}
            onClick={() => {
              if (isAddClickable) setQuantity((prev) => prev + 1);
            }}
          >
            <span
              className={isAddClickable ? "text-blue-500" : "text-gray-700"}
            >
              <Plus size="size-5" strokeWidth={2} />
            </span>
          </button>
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
            className="px-4 py-2 bg-green-500 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-green-600"
            onClick={handleConfirm}
          >
            <p className="font-sans font-medium text-white text-base">
              Confirm
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default OrderProductModal;
