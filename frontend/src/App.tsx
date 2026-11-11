import React, { useState } from "react";
import { Toaster } from "react-hot-toast";

import { useFetchAllProducts } from "./hooks";
import { formatPrice } from "./utils";

const App: React.FC = () => {
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(10);

  const {
    products,
    total: totalProducts,
    loading: isFetching,
  } = useFetchAllProducts(page, limit);

  return (
    <>
      <div className="p-10 h-screen w-full bg-white">
        <div className="flex flex-col justify-start items-center gap-10">
          <h2 className="font-sans font-medium text-gray-700 text-3xl text-center tracking-wider text-shadow-md">
            StockFlow
          </h2>

          <table className="w-full bg-white border-collapse rounded-sm overflow-hidden shadow-md">
            <thead>
              <tr className="py-2 flex bg-gray-200 font-sans font-medium text-base text-gray-700">
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
              {isFetching ? (
                [...Array(10)].map((_, index) => (
                  <tr key={index} className="p-3 flex animate-pulse">
                    <td className="flex-2 pr-4">
                      <div className="h-5 bg-gray-200 rounded-sm"></div>
                    </td>
                    <td className="flex-3 pr-4">
                      <div className="h-5 bg-gray-200 rounded-sm"></div>
                    </td>
                    <td className="flex-2 pr-4">
                      <div className="h-5 bg-gray-200 rounded-sm"></div>
                    </td>
                    <td className="flex-1 pr-4">
                      <div className="h-5 bg-gray-200 rounded-sm"></div>
                    </td>
                    <td className="flex-2 pr-4">
                      <div className="h-5 bg-gray-200 rounded-sm"></div>
                    </td>
                    <td className="flex-2 pr-4">
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
                    className="py-2 flex bg-white font-sans font-normal text-base text-gray-700 text-center cursor-pointer hover:bg-blue-200 hover:text-blue-500 active:bg-blue-500 active:text-white transition-all duration-200 ease-in-out"
                  >
                    <td className="flex-2">{p.sku}</td>
                    <td className="flex-3">{p.name}</td>
                    <td className="flex-2">{formatPrice(p.price)}</td>
                    <td className="flex-1">{p.stock}</td>
                    <td className="flex-2">
                      {p.is_active ? "Active" : "Inactive"}
                    </td>
                    <td className="flex-2">
                      {p.last_movement_at
                        ? new Date(p.last_movement_at).toLocaleString()
                        : ""}
                    </td>
                    <td className="flex-3">
                      <div className="flex flex-row justify-center items-center gap-2">
                        <button className="px-2 bg-orange-300 rounded-sm font-sans font-medium text-sm text-white cursor-pointer">
                          Update
                        </button>
                        <button className="px-2 bg-red-400 rounded-sm font-sans font-medium text-sm text-white cursor-pointer">
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>

          <div className="flex justify-center items-center gap-5">
            {page > 1 ? (
              <button
                className="cursor-pointer"
                onClick={() => setPage((prev) => prev - 1)}
              >
                <p className="font-sans font-medium text-blue-500 text-lg">
                  Prev
                </p>
              </button>
            ) : (
              <button>
                <p className="font-sans font-medium text-gray-500 text-lg">
                  Prev
                </p>
              </button>
            )}
            <p className="w-10 flex justify-center items-center bg-blue-200 rounded-sm font-sans font-medium text-blue-500 text-lg">
              {page}
            </p>
            {page * limit - totalProducts >= 0 ? (
              <button>
                <p className="font-sans font-medium text-gray-500 text-lg">
                  Next
                </p>
              </button>
            ) : (
              <button
                className="cursor-pointer"
                onClick={() => setPage((prev) => prev + 1)}
              >
                <p className="font-sans font-medium text-blue-500 text-lg">
                  Next
                </p>
              </button>
            )}
          </div>
        </div>
      </div>

      <Toaster position="bottom-right" />
    </>
  );
};

export default App;
