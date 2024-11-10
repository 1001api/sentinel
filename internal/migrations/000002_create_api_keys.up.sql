CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL DEFAULT 'default',
    token VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expired_at TIMESTAMPTZ,

    FOREIGN KEY(user_id) REFERENCES users(id)
);
