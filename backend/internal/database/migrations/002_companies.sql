

-- APPROVAL ENUM 
CREATE TYPE approval_status AS ENUM ('PENDING','APPROVED','REJECTED');

-- COMPANIES TABLE AGAIN WITH EVERYTHING

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

    --FOLLOWER SYSTEM 
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

-- Comments

COMMENT ON TABLE companies IS 'Companies require admin approval before they can create products';
COMMENT ON COLUMN companies.approval_status IS 'PENDING: Awaiting review | APPROVED: Can operate | REJECTED: Needs resubmission';
COMMENT ON COLUMN companies.rejection_reason IS 'Admin provides reason when rejecting';

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



-- =============================================
-- COMPANY SUBMISSION HISTORY (AUDIT LOG)
-- =============================================

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


--TRIGGERS

--1. auto log submission history
CREATE OR REPLACE FUNCTION log_company_submission()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO company_approval_history (company_id, action, performed_by_id, notes)
        VALUES (NEW.id, 'SUBMITTED', NEW.owner_id, 'Initial submission');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_company_submission
AFTER INSERT ON companies
FOR EACH ROW 
EXECUTE FUNCTION log_company_submission();


-- 2. update follower count 
CREATE OR REPLACE FUNCTION update_company_follower_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE companies 
        SET follower_count = follower_count + 1 
        WHERE id = NEW.company_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE companies 
        SET follower_count = GREATEST(follower_count - 1, 0)
        WHERE id = OLD.company_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_company_follower_count
AFTER INSERT OR DELETE ON company_followers
FOR EACH ROW
EXECUTE FUNCTION update_company_follower_count();


-- 3. Prevent following non-approved companies
CREATE OR REPLACE FUNCTION check_company_followable()
RETURNS TRIGGER AS $$
DECLARE
    comp_status approval_status;
    comp_active BOOLEAN;
BEGIN
    SELECT approval_status, is_active 
    INTO comp_status, comp_active
    FROM companies 
    WHERE id = NEW.company_id;
    
    IF comp_status != 'APPROVED' THEN
        RAISE EXCEPTION 'Cannot follow company that is not approved';
    END IF;
    
    IF NOT comp_active THEN
        RAISE EXCEPTION 'Cannot follow inactive company';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_company_followable
BEFORE INSERT ON company_followers
FOR EACH ROW
EXECUTE FUNCTION check_company_followable();



-- =============================================
-- CONSTRAINT: Products require approved company
-- =============================================

-- This will be enforced at application level in service layer
-- COMMENT ON TABLE products IS 'Products can only be created by users who own an APPROVED company';




CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC NOT NULL,
    unit TEXT NOT NULL,
    origin TEXT,
    


    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE INDEX idx_product_company_id ON products(company_id);
CREATE INDEX idx_product_is_active ON products(is_active);


CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);


CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    label TEXT NOT NULL, --"250kg", "500kg", "1kg"
    quantity_value NUMERIC NOT NULL, --250, 500, 1
    quantity_unit TEXT NOT NULL, --"kg"
    price NUMERIC(10,2) NOT NULL, -- exact money, max: 99999999.99
    isActive BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
