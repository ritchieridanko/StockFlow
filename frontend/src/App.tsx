import React, { useState } from "react";
import { Toaster } from "react-hot-toast";

import {
  CreateProductModal,
  DeleteProductModal,
  OrderList,
  Paginator,
  ProductList,
  StockMovementsModal,
  UpdateProductModal,
} from "./components";
import { Toggler } from "./components/ui";
import { useFetchAllOrders, useFetchAllProducts } from "./hooks";
import type { Order, ProductListItem, Tab } from "./types";

const App: React.FC = () => {
  const [page, setPage] = useState<number>(1);
  const [size, setSize] = useState<number>(10);
  const [tab, setTab] = useState<Tab>("Products");

  const {
    products,
    total: totalProducts,
    loading: isFetchingProducts,
    refetch: refetchProducts,
    append: appendProduct,
    replace: replaceProduct,
    remove: removeProduct,
  } = useFetchAllProducts(page, size, tab === "Products");
  const {
    orders,
    total: totalOrders,
    loading: isFetchingOrders,
  } = useFetchAllOrders(page, size, tab === "Orders");

  const [selectedProduct, setSelectedProduct] =
    useState<ProductListItem | null>(null);
  const [createProductModalVisible, setCreateProductModalVisible] =
    useState<boolean>(false);
  const [updateProductModalVisible, setUpdateProductModalVisible] =
    useState<boolean>(false);
  const [deleteProductModalVisible, setDeleteProductModalVisible] =
    useState<boolean>(false);
  const [stockMovementsModalVisible, setStockMovementsModalVisible] =
    useState<boolean>(false);

  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);

  return (
    <>
      <div className="p-10 h-screen w-full bg-gray-50">
        <div className="flex flex-col justify-start items-center gap-10">
          <h2 className="font-sans font-medium text-gray-700 text-3xl text-center tracking-wider text-shadow-md select-none">
            StockFlow
          </h2>

          <div className="w-full flex flex-col justify-start items-center gap-4">
            <div className="w-full flex flex-row justify-between items-center gap-4">
              <Toggler
                state1={{
                  label: "Products",
                  onClick: () => {
                    setPage(1);
                    setTab("Products");
                  },
                }}
                state2={{
                  label: "Orders",
                  onClick: () => {
                    setPage(1);
                    setTab("Orders");
                  },
                }}
                currentState={tab}
              />
              <button
                className="px-4 py-2 bg-blue-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-200"
                onClick={() => setCreateProductModalVisible(true)}
              >
                <p className="font-sans font-medium text-blue-500 text-base">
                  New Product
                </p>
              </button>
            </div>

            {tab === "Products" ? (
              <>
                <ProductList
                  products={products}
                  isLoading={isFetchingProducts}
                  onClickRow={(v) => {
                    setSelectedProduct(v);
                    setStockMovementsModalVisible(true);
                  }}
                  onClickUpdate={(v) => {
                    setSelectedProduct(v);
                    setUpdateProductModalVisible(true);
                  }}
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
              </>
            ) : (
              <>
                <OrderList
                  orders={orders}
                  isLoading={isFetchingOrders}
                  onClickRow={(v) => {
                    setSelectedOrder(v);
                  }}
                  onClickCancel={(v) => {
                    setSelectedOrder(v);
                  }}
                />
                <Paginator
                  currentPage={page}
                  totalPages={Math.ceil(totalOrders / size)}
                  onPageChange={(v) => setPage((prev) => prev + v)}
                />
              </>
            )}
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

      <UpdateProductModal
        isVisible={updateProductModalVisible}
        product={selectedProduct}
        onClickUpdate={(v) => {
          replaceProduct(v, selectedProduct?.last_movement_at);
          setUpdateProductModalVisible(false);
          setSelectedProduct(null);
        }}
        onClickCancel={() => {
          setUpdateProductModalVisible(false);
          setSelectedProduct(null);
        }}
      />

      <DeleteProductModal
        isVisible={deleteProductModalVisible}
        product={selectedProduct}
        onClickDelete={(v) => {
          removeProduct(v);
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
        onClickCancel={() => {
          setDeleteProductModalVisible(false);
          setSelectedProduct(null);
        }}
      />

      <StockMovementsModal
        isVisible={stockMovementsModalVisible}
        product={selectedProduct}
        onUpdate={() => {
          refetchProducts();
        }}
        onClickOutside={() => {
          setStockMovementsModalVisible(false);
          setSelectedProduct(null);
        }}
      />

      <Toaster position="bottom-right" />
    </>
  );
};

export default App;
