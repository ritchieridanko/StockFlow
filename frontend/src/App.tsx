import React, { useState } from "react";
import { Toaster } from "react-hot-toast";

import {
  CreateProductModal,
  DeleteProductModal,
  Paginator,
  ProductList,
} from "./components";
import { useDeleteProduct, useFetchAllProducts } from "./hooks";
import type { ProductListItem } from "./types";

const App: React.FC = () => {
  const [page, setPage] = useState<number>(1);
  const [size, setSize] = useState<number>(10);

  const {
    products,
    total: totalProducts,
    loading: isFetchingProducts,
    append: appendProduct,
    remove: removeProduct,
  } = useFetchAllProducts(page, size);
  const { deleteProduct } = useDeleteProduct();

  const [selectedProduct, setSelectedProduct] =
    useState<ProductListItem | null>(null);
  const [createProductModalVisible, setCreateProductModalVisible] =
    useState<boolean>(false);
  const [deleteProductModalVisible, setDeleteProductModalVisible] =
    useState<boolean>(false);

  const handleDelete = async () => {
    if (!selectedProduct) return;

    const res: boolean = await deleteProduct(selectedProduct.id);
    if (res) {
      removeProduct(selectedProduct.id);
    }
  };

  return (
    <>
      <div className="p-10 h-screen w-full bg-gray-50">
        <div className="flex flex-col justify-start items-center gap-10">
          <h2 className="font-sans font-medium text-gray-700 text-3xl text-center tracking-wider text-shadow-md select-none">
            StockFlow
          </h2>

          <div className="w-full flex flex-col justify-start items-center gap-4">
            <div className="w-full flex flex-row justify-end items-center gap-4">
              <button
                className="px-4 py-2 bg-blue-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-200"
                onClick={() => setCreateProductModalVisible(true)}
              >
                <p className="font-sans font-medium text-blue-500 text-base">
                  New Product
                </p>
              </button>
            </div>

            <ProductList
              products={products}
              isLoading={isFetchingProducts}
              onClickUpdate={() => {}}
              onClickDelete={(v) => {
                setSelectedProduct(v);
                setDeleteProductModalVisible(true);
              }}
            />

            <Paginator
              currentPage={page}
              totalPages={Math.ceil(totalProducts / size)}
              onPageChange={(v) => setPage((prev) => prev + v)}
            />
          </div>
        </div>
      </div>

      <CreateProductModal
        isVisible={createProductModalVisible}
        onClickCreate={(v) => {
          appendProduct(v);
          setCreateProductModalVisible(false);
        }}
        onClickCancel={() => setCreateProductModalVisible(false)}
      />

      <DeleteProductModal
        isVisible={deleteProductModalVisible}
        product={selectedProduct}
        onClickDelete={() => {
          handleDelete();
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
        onClickCancel={() => {
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
      />

      <Toaster position="bottom-right" />
    </>
  );
};

export default App;
