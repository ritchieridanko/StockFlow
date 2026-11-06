CREATE TABLE products (
  id SERIAL PRIMARY KEY,

  sku TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  price INT NOT NULL CHECK (price >= 0),
  stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index records by sku
CREATE INDEX idx_products_sku ON products USING GIN (sku gin_trgm_ops);

-- Index records by name
CREATE INDEX idx_products_name ON products USING GIN (name gin_trgm_ops);
