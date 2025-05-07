CREATE TABLE IF NOT EXISTS otas (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL,
    version_name VARCHAR(255) NOT NULL,
    version_code INTEGER NOT NULL,
    url TEXT NOT NULL,
    release_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_otas_app_id ON otas(app_id);
CREATE INDEX idx_otas_version_code ON otas(version_code); 