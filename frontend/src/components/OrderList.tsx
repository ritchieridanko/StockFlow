import React from "react";

import type { Order } from "../types";
import { formatPrice } from "../utils";

interface OrderListProps {
  orders: Order[];
  isLoading: boolean;
  onClickRow: (value: Order) => void;
  onClickCancel: (value: Order) => void;
}

const OrderList: React.FC<OrderListProps> = ({
  orders,
  isLoading,
  onClickRow,
  onClickCancel,
}: OrderListProps) => {
  return (
    <table className="w-full bg-white border-collapse rounded-sm overflow-hidden shadow-md">
      <thead>
        <tr className="py-2 flex flex-row bg-gray-100 font-sans font-medium text-gray-700 text-base">
          <th className="flex-3">Customer</th>
          <th className="flex-2">Total</th>
          <th className="flex-2">Status</th>
          <th className="flex-2">Created At</th>
          <th className="flex-2">Cancelled At</th>
          <th className="flex-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        {isLoading ? (
          [...Array(10)].map((_, index) => (
            <tr key={index} className="p-3 flex flex-row gap-4 animate-pulse">
              <td className="flex-3">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
            </tr>
          ))
        ) : orders.length === 0 ? (
          <tr className="py-10 flex justify-center items-center font-sans font-normal text-gray-500 text-base italic">
            <td colSpan={6}>No orders found</td>
          </tr>
        ) : (
          orders.map((o, index) => (
            <tr
              key={index}
              className="py-2 flex bg-white font-sans font-normal text-gray-700 text-base text-center cursor-pointer hover:bg-blue-100 hover:text-blue-500 active:bg-blue-200 transition-all duration-200 ease-in-out"
              onClick={() => onClickRow(o)}
            >
              <td className="flex-3">{o.customer}</td>
              <td className="flex-2">{formatPrice(o.total)}</td>
              <td className="flex-2">{o.status}</td>
              <td className="flex-2">
                {new Date(o.created_at).toLocaleString()}
              </td>
              <td className="flex-2">
                {o.cancelled_at
                  ? new Date(o.cancelled_at).toLocaleString()
                  : ""}
              </td>
              <td className="flex-2">
                <div className="flex flex-row justify-center items-center gap-2">
                  <button
                    className="px-2 bg-red-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-red-200"
                    onClick={(e) => {
                      e.stopPropagation();
                      onClickCancel(o);
                    }}
                  >
                    <p className="font-sans font-medium text-red-500 text-sm">
                      Cancel
                    </p>
                  </button>
                </div>
              </td>
            </tr>
          ))
        )}
      </tbody>
    </table>
  );
};

export default OrderList;
