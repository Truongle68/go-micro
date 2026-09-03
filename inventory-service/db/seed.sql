-- =========================================================================
-- SEED DATA (For local development and testing)
-- =========================================================================

-- 1. Warehouses
INSERT INTO warehouses (id, code, name, region, address_line1, address_ward, address_district, address_city, lat, lng, is_active)
VALUES 
    (gen_random_uuid(), 'WH-HN-01', 'Hanoi Central Hub', 'North', '123 Pham Van Dong', 'Co Nhue', 'Bac Tu Liem', 'Ha Noi', 21.053, 105.782, TRUE),
    (gen_random_uuid(), 'WH-HCM-01', 'Saigon Logistics Park', 'South', '456 Nguyen Van Linh', 'Tan Phong', 'District 7', 'Ho Chi Minh', 10.732, 106.703, TRUE)
ON CONFLICT (code) DO NOTHING;

-- 2. Suppliers
INSERT INTO suppliers (id, version, code, name, email, phone, address_line1, address_city, is_active)
VALUES
    (gen_random_uuid(), 1, 'SUP-TECH-01', 'Global Tech Imports Ltd', 'contact@globaltech.com', '+84901234567', '789 Le Duan', 'Ha Noi', TRUE)
ON CONFLICT (code) DO NOTHING;
