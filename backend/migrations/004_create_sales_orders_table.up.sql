CREATE TABLE sales_orders (
  id SERIAL PRIMARY KEY,

  customer TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'CREATED' CHECK (status IN ('CREATED','CANCELLED')),
  total INT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  cancelled_at TIMESTAMPTZ
);

-- Index records by customer
CREATE INDEX idx_sales_orders_customer ON sales_orders USING GIN (customer gin_trgm_ops);
