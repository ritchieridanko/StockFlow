import React from "react";
import toast from "react-hot-toast";

import { Trash } from "../assets/icons";
import { COLORS } from "../constants";
import { useCreateOrder } from "../hooks";
import type { CreateOrderPayload, Order, OrderProduct } from "../types";
import { TextField } from "./ui";

interface CreateOrderModalProps {
  isVisible: boolean;
  orderProducts: OrderProduct[];
  onClickCreate: (value: Order) => void;
  onClickDelete: (value: OrderProduct) => void;
  onClickCancel: () => void;
}

const CreateOrderModal: React.FC<CreateOrderModalProps> = ({
  isVisible,
  orderProducts,
  onClickCreate,
  onClickDelete,
  onClickCancel,
}: CreateOrderModalProps) => {
  const { form, setForm, clearForm, createOrder } = useCreateOrder();

  const handleCreate = async () => {
    const payload: CreateOrderPayload = {
      customer: form.customer,
      items: orderProducts.map(({ sku, quantity }) => ({ sku, quantity })),
    };
    if (payload.customer === "" || payload.items.length === 0) {
      toast.error("All fields are required", {
        id: "com-err-toast",
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

    const res: Order | undefined = await createOrder(payload);
    if (res) {
      onClickCreate(res);
    }

    clearForm();
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
          Create a New Order
        </p>
        <div className="flex flex-col justify-center items-center gap-4">
          <TextField
            type="text"
            placeholder="Enter the customer name"
            value={form.customer}
            onChange={(v) => setForm((form) => ({ ...form, customer: v }))}
          />
          <div className="w-full flex flex-col justify-center items-center gap-2">
            <p className="font-sans font-medium text-gray-700 text-lg">
              Order Items
            </p>
            {orderProducts.length === 0 ? (
              <p className="py-5 font-sans font-normal text-gray-500 text-base italic">
                No order items found
              </p>
            ) : (
              <div className="w-full flex flex-col justify-center items-center gap-2">
                {orderProducts.map((product, index) => (
                  <div
                    key={index}
                    className="px-3 py-2 w-full flex flex-row justify-between items-center bg-blue-100 rounded-sm"
                  >
                    <div className="flex flex-5 flex-col justify-center items-start">
                      <p className="font-sans font-medium text-blue-500 text-sm">
                        {product.name}
                      </p>
                      <p className="font-sans font-normal text-blue-500 text-xs">
                        {product.sku}
                      </p>
                    </div>
                    <div className="flex flex-2 flex-row justify-between items-center">
                      <p className="font-sans font-medium text-blue-500 text-base">
                        Qty {product.quantity}
                      </p>
                      <button
                        className="px-2 py-1 bg-red-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-red-200"
                        onClick={(e) => {
                          e.stopPropagation();
                          onClickDelete(product);
                        }}
                      >
                        <span className="text-red-500">
                          <Trash size="size-4" strokeWidth={2} />
                        </span>
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
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
              Create Order
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default CreateOrderModal;
