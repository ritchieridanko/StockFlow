import React, { useEffect, useState } from "react";
import toast from "react-hot-toast";

import { COLORS } from "../constants";
import { useFetchAllStockMovements, useUpdateStock } from "../hooks";
import type { Product, ProductListItem, StockMovementType } from "../types";
import { Dropdown, TextField } from "./ui";

interface StockMovementsModalProps {
  isVisible: boolean;
  product: ProductListItem | null;
  onUpdate: () => void;
  onClickOutside: () => void;
}

const StockMovementsModal: React.FC<StockMovementsModalProps> = ({
  isVisible,
  product,
  onUpdate,
  onClickOutside,
}: StockMovementsModalProps) => {
  const [updateFormVisible, setUpdateFormVisible] = useState<boolean>(false);

  const { stockMovements, loading: isFetchingStockMovements } =
    useFetchAllStockMovements(product?.id ?? 0);
  const {
    form: updateStockForm,
    setForm: setUpdateStockForm,
    clearForm: clearUpdateStockForm,
    updateStock,
  } = useUpdateStock();

  useEffect(() => {
    if (product) {
      setUpdateFormVisible(false);
      clearUpdateStockForm();
    }
  }, [product]);

  const handleCreate = async () => {
    if (!product) return;
    if (updateStockForm.type !== "in" && updateStockForm.type !== "out") {
      toast.error("Invalid movement type", {
        id: "smm-err-toast",
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

    const res: Product | undefined = await updateStock(product.id);
    if (res) {
      onClickOutside();
      onUpdate();
    }
  };

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
            {isFetchingStockMovements ? (
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
        {!updateFormVisible && (
          <div className="w-full flex flex-row justify-end items-center">
            <button
              className="px-4 py-2 bg-blue-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-200"
              onClick={() => setUpdateFormVisible(true)}
            >
              <p className="font-sans font-medium text-blue-500 text-base">
                Update Stock
              </p>
            </button>
          </div>
        )}
        {updateFormVisible && (
          <div className="mt-10 w-full flex flex-col justify-center gap-4">
            <p className="font-sans font-medium text-gray-700 text-xl">
              Create a New Stock Movement
            </p>
            <div className="flex flex-col justify-center items-center gap-4">
              <Dropdown
                selectedOption={updateStockForm.type}
                defaultOption="Select a movement type"
                options={["In", "Out"]}
                onSelect={(v) =>
                  setUpdateStockForm((form) => ({
                    ...form,
                    type: v as StockMovementType & "",
                  }))
                }
              />
              <TextField
                type="number"
                placeholder="Enter the movement quantity"
                value={String(updateStockForm.quantity)}
                onChange={(v) =>
                  setUpdateStockForm((form) => ({
                    ...form,
                    quantity: Number(v),
                  }))
                }
                min={1}
                step={1}
              />
              <TextField
                type="text"
                placeholder="Enter the movement note"
                value={updateStockForm.note ?? ""}
                onChange={(v) =>
                  setUpdateStockForm((form) => ({ ...form, note: v }))
                }
              />
            </div>
            <div className="flex flex-row justify-end items-center gap-4">
              <button
                className="px-4 py-2 bg-gray-200 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-gray-300"
                onClick={() => setUpdateFormVisible(false)}
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
                  Create Stock Movement
                </p>
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default StockMovementsModal;
