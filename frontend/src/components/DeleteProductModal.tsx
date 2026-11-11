import React from "react";

import type { ProductListItem } from "../types";

interface DeleteProductModalProps {
  isVisible: boolean;
  product: ProductListItem | null;
  onClickCancel: () => void;
  onClickDelete: () => void;
}

const DeleteProductModal: React.FC<DeleteProductModalProps> = ({
  isVisible,
  product,
  onClickCancel,
  onClickDelete,
}: DeleteProductModalProps) => {
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
          Confirm Deletion
        </p>
        <p className="font-sans font-normal text-gray-600 text-base">
          Are you sure you want to delete{" "}
          <span className="font-sans font-medium text-gray-700 text-base">
            {product?.name}
          </span>{" "}
          ({product?.sku})?
        </p>
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
            className="px-4 py-2 bg-red-500 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-red-600"
            onClick={onClickDelete}
          >
            <p className="font-sans font-medium text-white text-base">
              Delete Product
            </p>
          </button>
        </div>
      </div>
    </div>
  );
};

export default DeleteProductModal;
