import React from "react";

import type { ProductListItem } from "../types";
import { formatPrice } from "../utils";

interface ProductListProps {
  products: ProductListItem[];
  isLoading: boolean;
  onClickUpdate: (value: ProductListItem) => void;
  onClickDelete: (value: ProductListItem) => void;
}

const ProductList: React.FC<ProductListProps> = ({
  products,
  isLoading,
  onClickUpdate,
  onClickDelete,
}: ProductListProps) => {
  return (
    <table className="w-full bg-white border-collapse rounded-sm overflow-hidden shadow-md">
      <thead>
        <tr className="py-2 flex flex-row bg-gray-100 font-sans font-medium text-gray-700 text-base">
          <th className="flex-2">SKU</th>
          <th className="flex-3">Name</th>
          <th className="flex-2">Price</th>
          <th className="flex-1">Stock</th>
          <th className="flex-2">Status</th>
          <th className="flex-2">Last Movement</th>
          <th className="flex-3">Actions</th>
        </tr>
      </thead>
      <tbody>
        {isLoading ? (
          [...Array(10)].map((_, index) => (
            <tr key={index} className="p-3 flex flex-row gap-4 animate-pulse">
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-3">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-1">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-2">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
              <td className="flex-3">
                <div className="h-5 bg-gray-200 rounded-sm"></div>
              </td>
            </tr>
          ))
        ) : products.length === 0 ? (
          <tr className="py-10 flex justify-center items-center font-sans font-normal text-gray-500 text-base italic">
            <td colSpan={7}>No products found</td>
          </tr>
        ) : (
          products.map((p, index) => (
            <tr
              key={index}
              className="py-2 flex bg-white font-sans font-normal text-gray-700 text-base text-center cursor-pointer hover:bg-blue-100 hover:text-blue-500 active:bg-blue-200 transition-all duration-200 ease-in-out"
            >
              <td className="flex-2">{p.sku}</td>
              <td className="flex-3">{p.name}</td>
              <td className="flex-2">{formatPrice(p.price)}</td>
              <td className="flex-1">{p.stock}</td>
              <td className="flex-2">{p.is_active ? "Active" : "Inactive"}</td>
              <td className="flex-2">
                {p.last_movement_at
                  ? new Date(p.last_movement_at).toLocaleString()
                  : ""}
              </td>
              <td className="flex-3">
                <div className="flex flex-row justify-center items-center gap-2">
                  <button
                    className="px-2 bg-orange-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-orange-200"
                    onClick={() => onClickUpdate(p)}
                  >
                    <p className="font-sans font-medium text-orange-500 text-sm">
                      Update
                    </p>
                  </button>
                  <button
                    className="px-2 bg-red-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-red-200"
                    onClick={() => onClickDelete(p)}
                  >
                    <p className="font-sans font-medium text-red-500 text-sm">
                      Delete
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

export default ProductList;
