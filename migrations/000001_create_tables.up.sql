CREATE TABLE users (
    id UUID PRIMARY KEY default gen_random_uuid(),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    user_role VARCHAR(100) default 'user' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID references users(id),
    refresh_token text not null,
    expires_in TIMESTAMP not null,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
