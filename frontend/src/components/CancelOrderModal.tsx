import React from "react";

import { useCancelOrder } from "../hooks";
import type { Order } from "../types";
import { formatPrice } from "../utils";

interface CancelOrderModalProps {
  isVisible: boolean;
  order: Order | null;
  onClickConfirm: (value: Order) => void;
  onClickNo: () => void;
}

const CancelOrderModal: React.FC<CancelOrderModalProps> = ({
  isVisible,
  order,
  onClickConfirm,
  onClickNo,
}: CancelOrderModalProps) => {
  const { cancelOrder } = useCancelOrder();

  const handleCancel = async () => {
    if (!order) return;

    const res: Order | undefined = await cancelOrder(order.id);
    if (res) {
      onClickConfirm(res);
    }
  };

  if (!isVisible) return null;

  return (
    <div
      className="fixed z-99 inset-0 flex justify-center items-center bg-black/50 backdrop-blur-xs cursor-pointer"
      onClick={onClickNo}
    >
      <div
        className="p-5 w-full max-w-md flex flex-col justify-center gap-4 bg-white rounded-lg cursor-default shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="font-sans font-medium text-gray-700 text-xl">
          Order Cancellation
        </p>
        <div className="flex flex-col justify-center items-start">
          <p className="font-sans font-normal text-gray-600 text-base">
            Are you sure you want to cancel Order #{order?.id}?
          </p>
          <p className="font-sans font-normal text-gray-600 text-base">
            Customer:{" "}
            <span className="font-sans font-medium text-gray-700 text-base">
              {order?.customer}
            </span>
          </p>
          <p className="font-sans font-normal text-gray-600 text-base">
            Total:{" "}
            <span className="font-sans font-medium text-gray-700 text-base">
              {formatPrice(order?.total ?? 0)}
            </span>
          </p>
        </div>
        <div className="flex flex-row justify-end items-center gap-4">
          <button
            className="px-4 py-2 bg-gray-200 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-gray-300"
            onClick={onClickNo}
          >
            <p className="font-sans font-medium text-gray-700 text-base">No</p>
          </button>
          <button
            className="px-4 py-2 bg-red-500 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-red-600"
            onClick={handleCancel}
          >
            <p className="font-sans font-medium text-white text-base">
              Cancel Order
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default CancelOrderModal;
