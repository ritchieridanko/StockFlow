CREATE TABLE stock_movements (
  id SERIAL PRIMARY KEY,
  product_id INT NOT NULL REFERENCES products (id),

  type TEXT NOT NULL CHECK (type IN ('in', 'out')),
  qty INT NOT NULL CHECK (qty > 0),
  note TEXT,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index records by product id
CREATE INDEX idx_stock_movements_product_id ON stock_movements (product_id);
