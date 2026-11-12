import React from "react";

import { useFetchAllStockMovements } from "../hooks";
import type { ProductListItem } from "../types";

interface StockMovementsModalProps {
  isVisible: boolean;
  product: ProductListItem | null;
  onClickOutside: () => void;
}

const StockMovementsModal: React.FC<StockMovementsModalProps> = ({
  isVisible,
  product,
  onClickOutside,
}: StockMovementsModalProps) => {
  const { stockMovements, loading: isLoading } = useFetchAllStockMovements(
    product?.id ?? 0,
  );

  if (!isVisible) return null;

  return (
    <div
      className="fixed z-99 inset-0 flex justify-center items-center bg-black/50 backdrop-blur-xs cursor-pointer"
      onClick={onClickOutside}
    >
      <div
        className="p-5 w-full max-w-5xl flex flex-col justify-center gap-4 bg-white rounded-lg cursor-default shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="font-sans font-medium text-gray-700 text-xl">
          Stock Movements
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
        </div>
        <table className="w-full bg-white border-collapse rounded-sm overflow-hidden shadow-md">
          <thead>
            <tr className="py-2 flex flex-row bg-gray-100 font-sans font-medium text-gray-700 text-base">
              <th className="flex-2">Type</th>
              <th className="flex-2">Quantity</th>
              <th className="flex-4">Note</th>
              <th className="flex-2">Created At</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              [...Array(10)].map((_, index) => (
                <tr
                  key={index}
                  className="p-3 flex flex-row gap-4 animate-pulse"
                >
                  <td className="flex-2">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-2">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-4">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-2">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                </tr>
              ))
            ) : stockMovements.length === 0 ? (
              <tr className="py-10 flex justify-center items-center font-sans font-normal text-gray-500 text-base italic">
                <td colSpan={4}>No stock movements found</td>
              </tr>
            ) : (
              stockMovements.map((sm, index) => (
                <tr
                  key={index}
                  className="py-2 flex bg-white font-sans font-normal text-gray-700 text-base text-center cursor-pointer hover:bg-blue-100 hover:text-blue-500 active:bg-blue-200 transition-all duration-200 ease-in-out"
                >
                  <td className="flex-2">
                    {sm.type === "in" ? "Stock In" : "Stock Out"}
                  </td>
                  <td className="flex-2">{sm.quantity}</td>
                  <td className="flex-4">{sm.note ?? ""}</td>
                  <td className="flex-2">
                    {new Date(sm.created_at).toLocaleString()}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default StockMovementsModal;
