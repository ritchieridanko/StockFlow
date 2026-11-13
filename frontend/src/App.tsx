import React, { useState } from "react";
import { Toaster } from "react-hot-toast";

import {
  CancelOrderModal,
  CreateOrderModal,
  CreateProductModal,
  DeleteProductModal,
  OrderList,
  OrderProductModal,
  Paginator,
  ProductList,
  StockMovementsModal,
  UpdateProductModal,
} from "./components";
import { Toggler } from "./components/ui";
import { useFetchAllOrders, useFetchAllProducts } from "./hooks";
import type { Order, OrderProduct, ProductListItem, Tab } from "./types";

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
    append: appendOrder,
    replace: replaceOrder,
  } = useFetchAllOrders(page, size, tab === "Orders");

  const [selectedProduct, setSelectedProduct] =
    useState<ProductListItem | null>(null);
  const [createProductModalVisible, setCreateProductModalVisible] =
    useState<boolean>(false);
  const [orderProductModalVisible, setOrderProductModalVisible] =
    useState<boolean>(false);
  const [updateProductModalVisible, setUpdateProductModalVisible] =
    useState<boolean>(false);
  const [deleteProductModalVisible, setDeleteProductModalVisible] =
    useState<boolean>(false);
  const [stockMovementsModalVisible, setStockMovementsModalVisible] =
    useState<boolean>(false);

  const [orderProducts, setOrderProducts] = useState<OrderProduct[]>([]);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [createOrderModalVisible, setCreateOrderModalVisible] =
    useState<boolean>(false);
  const [cancelOrderModalVisible, setCancelOrderModalVisible] =
    useState<boolean>(false);

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
              {tab === "Products" ? (
                <button
                  className="px-4 py-2 bg-blue-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-200"
                  onClick={() => setCreateProductModalVisible(true)}
                >
                  <p className="font-sans font-medium text-blue-500 text-base">
                    New Product
                  </p>
                </button>
              ) : (
                <button
                  className="px-4 py-2 bg-blue-100 rounded-sm cursor-pointer transition-all duration-200 ease-in-out hover:bg-blue-200"
                  onClick={() => setCreateOrderModalVisible(true)}
                >
                  <p className="font-sans font-medium text-blue-500 text-base">
                    New Order
                  </p>
                </button>
              )}
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
                  onClickOrder={(v) => {
                    setSelectedProduct(v);
                    setOrderProductModalVisible(true);
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
                {!isFetchingProducts && products.length > 0 && (
                  <Paginator
                    currentPage={page}
                    totalPages={Math.ceil(totalProducts / size)}
                    onPageChange={(v) => setPage((prev) => prev + v)}
                  />
                )}
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
                    setCancelOrderModalVisible(true);
                  }}
                />
                {!isFetchingOrders && orders.length > 0 && (
                  <Paginator
                    currentPage={page}
                    totalPages={Math.ceil(totalOrders / size)}
                    onPageChange={(v) => setPage((prev) => prev + v)}
                  />
                )}
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

      <OrderProductModal
        isVisible={orderProductModalVisible}
        product={selectedProduct}
        onClickConfirm={(product, quantity) => {
          setOrderProducts((prev) => [
            ...prev,
            {
              name: product.name,
              sku: product.sku,
              quantity: quantity,
            },
          ]);
          setOrderProductModalVisible(false);
          setSelectedProduct(null);
        }}
        onClickCancel={() => {
          setOrderProductModalVisible(false);
          setSelectedProduct(null);
        }}
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

      <CreateOrderModal
        isVisible={createOrderModalVisible}
        orderProducts={orderProducts}
        onClickCreate={(v) => {
          appendOrder(v);
          setOrderProducts([]);
          setCreateOrderModalVisible(false);
        }}
        onClickDelete={(v) => {
          setOrderProducts((prev) =>
            prev.filter((product) => product.sku !== v.sku),
          );
        }}
        onClickCancel={() => setCreateOrderModalVisible(false)}
      />

      <CancelOrderModal
        isVisible={cancelOrderModalVisible}
        order={selectedOrder}
        onClickConfirm={(v) => {
          replaceOrder(v);
          setCancelOrderModalVisible(false);
          setSelectedOrder(null);
        }}
        onClickNo={() => {
          setCancelOrderModalVisible(false);
          setSelectedOrder(null);
        }}
      />

      <Toaster position="bottom-right" />
    </>
  );
};

export default App;
