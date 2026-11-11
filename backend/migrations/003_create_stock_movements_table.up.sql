CREATE TABLE stock_movements (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL,

  type TEXT NOT NULL CHECK (type IN ('in', 'out')),
  quantity INT NOT NULL CHECK (quantity > 0),
  note TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

-- Index records by product id
CREATE INDEX idx_stock_movements_product_id ON stock_movements (product_id);
