INSERT INTO
  products (sku, name, price, stock)
VALUES
  ('SEED-001', 'ProBook 14-inch Laptop 16GB RAM', 12500000, 15),
  ('SEED-002', 'UltraWide 29-inch IPS Monitor', 3200000, 25),
  ('SEED-003', 'Mechanical RGB Gaming Keyboard', 850000, 50),
  ('SEED-004', 'Wireless Ergonomic Mouse', 350000, 8),
  ('SEED-005', 'Active Noise Cancelling Headphones', 1450000, 30),
  ('SEED-006', 'Standing Electric Desk 120cm', 4500000, 5),
  ('SEED-007', 'Ergonomic Mesh Office Chair', 2100000, 12),
  ('SEED-008', 'USB-C 7-in-1 Multiport Hub', 450000, 40),
  ('SEED-009', 'Vertical Laptop Dual Stand', 250000, 6),
  ('SEED-010', 'Braided USB-C to USB-C Cable 2m', 120000, 100),
  ('SEED-011', 'Wi-Fi 6 Mesh Router System', 1850000, 18),
  ('SEED-012', '8-Port Gigabit Desktop Switch', 380000, 22)
ON CONFLICT (sku) DO NOTHING;
