
-- =============================================
-- ENUM TYPE
-- =============================================

CREATE TYPE approval_status AS ENUM ('PENDING','APPROVED','REJECTED');


-- =============================================
-- CORE TABLES (Users and Auth)
-- =============================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- User Profile
    email TEXT NULL UNIQUE,
    name TEXT NOT NULL,
    profile_image_url TEXT,
    phone TEXT,
    role TEXT NOT NULL CHECK (role IN ('USER','ADMIN')) DEFAULT 'USER',

    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    last_login_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_auth_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    auth_provider TEXT NOT NULL CHECK (auth_provider IN ('LOCAL', 'GOOGLE')),
    
    oauth_sub TEXT,
    password_hash TEXT, 
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT user_provider_unique UNIQUE (user_id,auth_provider),
    CONSTRAINT oauth_sub_unique UNIQUE (oauth_sub)
);

CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    address TEXT NOT NULL,
    city TEXT,
    state TEXT,
    pincode TEXT,

    is_primary BOOLEAN NOT NULL DEFAULT FALSE  
); 


-- =============================================
-- REFRESH TOKENS
-- =============================================

CREATE TABLE refresh_token (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT NOT NULL,
    ip_address TEXT NOT NULL,

    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_token_user_id ON refresh_token(user_id);
CREATE INDEX idx_refresh_token_expires_at ON refresh_token(expires_at);