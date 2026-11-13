import React from "react";

import { useFetchOrderById } from "../hooks";
import type { Order } from "../types";
import { formatPrice } from "../utils";

interface OrderItemsModalProps {
  isVisible: boolean;
  order: Order | null;
  onClickOutside: () => void;
}

const OrderItemsModal: React.FC<OrderItemsModalProps> = ({
  isVisible,
  order,
  onClickOutside,
}: OrderItemsModalProps) => {
  const { order: orderWithItems, loading: isFetchingOrderById } =
    useFetchOrderById(order?.id ?? 0);

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
          Order Items
        </p>
        <div className="flex flex-row justify-between items-center">
          <div className="flex flex-col justify-center items-start">
            <p className="font-sans font-normal text-gray-600 text-base">
              Customer:{" "}
              <span className="font-sans font-medium text-gray-700 text-base">
                {orderWithItems?.customer || order?.customer}
              </span>
            </p>
            <p className="font-sans font-normal text-gray-600 text-base">
              Total:{" "}
              <span className="font-sans font-medium text-gray-700 text-base">
                {formatPrice(orderWithItems?.total || order?.total || 0)}
              </span>
            </p>
          </div>
          {orderWithItems?.cancelled_at && (
            <p className="px-4 py-2 bg-red-500 rounded-sm font-sans font-medium text-white text-base">
              CANCELLED
            </p>
          )}
        </div>
        <table className="w-full bg-white border-collapse rounded-sm overflow-hidden shadow-md">
          <thead>
            <tr className="py-2 flex flex-row bg-gray-100 font-sans font-medium text-gray-700 text-base">
              <th className="flex-3">SKU</th>
              <th className="flex-2">Quantity</th>
              <th className="flex-3">Unit Price</th>
              <th className="flex-3">Created At</th>
            </tr>
          </thead>
          <tbody>
            {isFetchingOrderById ? (
              [...Array(10)].map((_, index) => (
                <tr
                  key={index}
                  className="p-3 flex flex-row gap-4 animate-pulse"
                >
                  <td className="flex-3">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-2">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-3">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                  <td className="flex-3">
                    <div className="h-5 bg-gray-200 rounded-sm"></div>
                  </td>
                </tr>
              ))
            ) : !orderWithItems || orderWithItems.items.length === 0 ? (
              <tr className="py-10 flex justify-center items-center font-sans font-normal text-gray-500 text-base italic">
                <td colSpan={4}>No order items found</td>
              </tr>
            ) : (
              orderWithItems.items.map((o, index) => (
                <tr
                  key={index}
                  className="py-2 flex bg-white font-sans font-normal text-gray-700 text-base text-center cursor-pointer hover:bg-blue-100 hover:text-blue-500 active:bg-blue-200 transition-all duration-200 ease-in-out"
                >
                  <td className="flex-3">{o.sku}</td>
                  <td className="flex-2">{o.quantity}</td>
                  <td className="flex-3">{formatPrice(o.unit_price)}</td>
                  <td className="flex-3">
                    {new Date(o.created_at).toLocaleString()}
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

export default OrderItemsModal;
