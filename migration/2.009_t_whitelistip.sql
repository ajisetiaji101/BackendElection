CREATE TABLE ip_whitelist (
    id SERIAL PRIMARY KEY,          -- Auto-incrementing unique identifier
    ip_address VARCHAR(255) NOT NULL,       -- IP address (supports both IPv4 and IPv6)
    description TEXT,               -- Optional description or reason for whitelisting
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the entry was created
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP -- Timestamp for the last update
);