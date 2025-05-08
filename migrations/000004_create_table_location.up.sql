CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    zone VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- rack/bin/area
    capacity DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_locations_code ON locations(code);
CREATE INDEX idx_locations_zone ON locations(zone);
CREATE INDEX idx_locations_type ON locations(type); 