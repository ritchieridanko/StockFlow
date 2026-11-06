CREATE TABLE sales_order_items (
  id SERIAL PRIMARY KEY,
  order_id INT NOT NULL REFERENCES sales_orders (id),

  sku TEXT NOT NULL REFERENCES products (sku),
  qty INT NOT NULL CHECK (qty > 0),
  unit_price INT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index records by order id
CREATE INDEX idx_sales_order_items_order_id ON sales_order_items (order_id);
