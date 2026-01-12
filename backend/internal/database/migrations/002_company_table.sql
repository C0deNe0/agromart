
-- =============================================
-- COMPANIES TABLE
-- =============================================

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    logo_url TEXT,

    business_email TEXT,
    business_phone TEXT,
    
    city TEXT,
    state TEXT,
    pincode TEXT,

    gst_number TEXT,
    pan_number TEXT,

    -- APPROVAL SYSTEM
    approval_status approval_status NOT NULL DEFAULT 'PENDING',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- ADMIN ACTIONS
    reviewed_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    -- FOLLOWER SYSTEM 
    follower_count INTEGER NOT NULL DEFAULT 0,
    product_visibility TEXT NOT NULL DEFAULT 'PUBLIC' 
        CHECK (product_visibility IN ('PUBLIC', 'FOLLOWERS_ONLY', 'PRIVATE')),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(owner_id,name),
    CHECK(length(name) >=3)
);

CREATE INDEX idx_companies_owner_id ON companies(owner_id);
CREATE INDEX idx_companies_approval_status ON companies(approval_status);
CREATE INDEX idx_companies_is_active ON companies(is_active);
CREATE INDEX idx_companies_submitted_at ON companies(submitted_at DESC);

COMMENT ON TABLE companies IS 'Companies require admin approval before they can create products';
COMMENT ON COLUMN companies.approval_status IS 'PENDING: Awaiting review | APPROVED: Can operate | REJECTED: Needs resubmission';
COMMENT ON COLUMN companies.rejection_reason IS 'Admin provides reason when rejecting';


-- =============================================
-- COMPANY RELATIONS & AUDIT
-- =============================================

CREATE TABLE company_followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    followed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(company_id, user_id)
);

CREATE INDEX idx_company_followers_company_id ON company_followers(company_id);
CREATE INDEX idx_company_followers_user_id ON company_followers(user_id);
CREATE INDEX idx_company_followers_followed_at ON company_followers(followed_at DESC);


CREATE TABLE company_approval_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    
    action TEXT NOT NULL CHECK (action IN ('SUBMITTED', 'APPROVED', 'REJECTED', 'RESUBMITTED')),
    performed_by_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    notes TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_approval_history_company_id ON company_approval_history(company_id);
CREATE INDEX idx_approval_history_created_at ON company_approval_history(created_at DESC);