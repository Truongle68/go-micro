CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS suppliers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version       INT NOT NULL DEFAULT 1,
    code          VARCHAR(50) NOT NULL UNIQUE,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL DEFAULT '',
    phone         VARCHAR(50) NOT NULL DEFAULT '',
    address_line1 VARCHAR(255) NOT NULL DEFAULT '',
    address_ward  VARCHAR(100) NOT NULL DEFAULT '',
    address_district VARCHAR(100) NOT NULL DEFAULT '',
    address_city  VARCHAR(100) NOT NULL DEFAULT '',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers(code);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers(is_active);

CREATE TABLE IF NOT EXISTS purchase_orders (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version       INT NOT NULL DEFAULT 1,
    code          VARCHAR(50) NOT NULL UNIQUE,
    supplier_id   UUID NOT NULL REFERENCES suppliers(id),
    supplier_name VARCHAR(255) NOT NULL DEFAULT '',
    warehouse_id  VARCHAR(255) NOT NULL,
    lines         JSONB NOT NULL DEFAULT '[]'::jsonb,
    status        VARCHAR(50) NOT NULL DEFAULT 'draft',
    created_by    VARCHAR(255) NOT NULL DEFAULT '',
    expected_at   TIMESTAMP WITH TIME ZONE,
    received_at   TIMESTAMP WITH TIME ZONE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_code ON purchase_orders(code);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_id ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_warehouse_id ON purchase_orders(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(status);

CREATE TABLE IF NOT EXISTS warehouses (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             VARCHAR(50) NOT NULL UNIQUE,
    name             VARCHAR(255) NOT NULL,
    region           VARCHAR(100) NOT NULL DEFAULT '',
    address_line1    VARCHAR(255) NOT NULL DEFAULT '',
    address_ward     VARCHAR(100) NOT NULL DEFAULT '',
    address_district VARCHAR(100) NOT NULL DEFAULT '',
    address_city     VARCHAR(100) NOT NULL DEFAULT '',
    lat              DOUBLE PRECISION NOT NULL DEFAULT 0,
    lng              DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_warehouses_code ON warehouses(code);
CREATE INDEX IF NOT EXISTS idx_warehouses_is_active ON warehouses(is_active);

CREATE TABLE IF NOT EXISTS stock_levels (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku               VARCHAR(100) NOT NULL,
    warehouse_id      UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    on_hand           INT NOT NULL DEFAULT 0,
    reserved          INT NOT NULL DEFAULT 0,
    reorder_threshold INT NOT NULL DEFAULT 0,
    reorder_quantity  INT NOT NULL DEFAULT 0,
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(sku, warehouse_id)
);

CREATE INDEX IF NOT EXISTS idx_stock_levels_sku ON stock_levels(sku);
CREATE INDEX IF NOT EXISTS idx_stock_levels_warehouse_id ON stock_levels(warehouse_id);

CREATE TABLE IF NOT EXISTS stock_reservations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku          VARCHAR(100) NOT NULL,
    warehouse_id UUID NOT NULL,
    quantity     INT NOT NULL,
    order_id     VARCHAR(255) NOT NULL,
    status       VARCHAR(50) NOT NULL DEFAULT 'pending',
    expires_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE,
    released_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_stock_reservations_order_id ON stock_reservations(order_id);
CREATE INDEX IF NOT EXISTS idx_stock_reservations_status ON stock_reservations(status);
CREATE INDEX IF NOT EXISTS idx_stock_reservations_expires_at ON stock_reservations(expires_at);

CREATE TABLE IF NOT EXISTS stock_movements (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku            VARCHAR(100) NOT NULL,
    warehouse_id   UUID NOT NULL,
    type           VARCHAR(50) NOT NULL,
    quantity       INT NOT NULL,
    reference_type VARCHAR(100) NOT NULL DEFAULT '',
    reference_id   VARCHAR(255) NOT NULL DEFAULT '',
    note           TEXT NOT NULL DEFAULT '',
    created_by     VARCHAR(255) NOT NULL DEFAULT '',
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stock_movements_sku ON stock_movements(sku);
CREATE INDEX IF NOT EXISTS idx_stock_movements_warehouse_id ON stock_movements(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference ON stock_movements(reference_type, reference_id);
