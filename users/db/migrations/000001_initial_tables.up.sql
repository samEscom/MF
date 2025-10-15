CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(14) PRIMARY KEY,
    name VARCHAR(255),
    last_name VARCHAR(255),
    age INT,
    email VARCHAR(255),
    is_active BOOLEAN,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone  
);
