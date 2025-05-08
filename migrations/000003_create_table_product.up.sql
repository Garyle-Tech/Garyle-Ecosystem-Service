CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(255) NOT NULL UNIQUE, --unique --ARS-1312321678 / 1232136278
    name VARCHAR(255) NOT NULL,
    description TEXT,
    unit VARCHAR(255) NOT NULL, --kg/pcs/box/lot
    weight DECIMAL(10, 2) NOT NULL, --24.5
    dimension VARCHAR(255) NOT NULL, --100x50x20
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_products_sku ON products(sku);
