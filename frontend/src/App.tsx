import React, { useState } from "react";
import { Toaster } from "react-hot-toast";

import { DeleteProductModal, Paginator, ProductList } from "./components";
import { useDeleteProduct, useFetchAllProducts } from "./hooks";
import type { ProductListItem } from "./types";

const App: React.FC = () => {
  const [page, setPage] = useState<number>(1);
  const [size, setSize] = useState<number>(10);

  const {
    products,
    total: totalProducts,
    loading: isFetchingProducts,
    remove: removeProduct,
  } = useFetchAllProducts(page, size);
  const { deleteProduct } = useDeleteProduct();

  const [selectedProduct, setSelectedProduct] =
    useState<ProductListItem | null>(null);
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

      <DeleteProductModal
        isVisible={deleteProductModalVisible}
        product={selectedProduct}
        onClickCancel={() => {
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
        onClickDelete={() => {
          handleDelete();
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
      />

      <Toaster position="bottom-right" />
    </>
  );
};

export default App;
