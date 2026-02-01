CREATE TABLE device_messages (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    device_id CHAR(36),
    device_type_code CHAR(100), 
    protocol CHAR(100), 
    message TEXT, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

