INSERT INTO vehicles (id, external_id, plate, name, created_at, updated_at)
VALUES
  ('01HV3K2M8BXYZ00000VEHICLE1', 'EXT-VH-001', 'ABC-123', 'Truck Alpha',   now(), now()),
  ('01HV3K2M8BXYZ00000VEHICLE2', 'EXT-VH-002', 'DEF-456', 'Van Beta',      now(), now()),
  ('01HV3K2M8BXYZ00000VEHICLE3', 'EXT-VH-003', 'GHI-789', 'Sedan Gamma',   now(), now()),
  ('01HV3K2M8BXYZ00000VEHICLE4', 'EXT-VH-004', 'JKL-012', 'SUV Delta',     now(), now()),
  ('01HV3K2M8BXYZ00000VEHICLE5', 'EXT-VH-005', 'MNO-345', 'Pickup Epsilon', now(), now())
ON CONFLICT DO NOTHING;
