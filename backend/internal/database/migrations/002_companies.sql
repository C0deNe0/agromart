

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
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,

    name TEXT NOT NULL,
    description TEXT,
    
    unit TEXT NOT NULL,
    origin TEXT,
    
    base_price NUMERIC(10,2) NOT NULL CHECK( base_price >=0),
    
    
    approval_status approval_status NOT NULL DEFAULT 'PENDING',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    reviewed_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMPTZ,
    rejection_reason TEXT,
    

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    
    CHECK (length(name) >= 3)
    );

CREATE INDEX idx_products_company_id ON products(company_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_approval_status ON products(approval_status);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_submitted_at ON products(submitted_at DESC);

COMMENT ON TABLE products IS 'Products require admin approval before being visible';
COMMENT ON COLUMN products.approval_status IS 'PENDING: Awaiting review | APPROVED: Visible to users | REJECTED: Needs modification';


CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    image_url TEXT NOT NULL,
    s3_key TEXT NOT NULL, -- to delete from s3 to
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_is_primary ON product_images(is_primary);

-- Ensure only one primary image per product
CREATE UNIQUE INDEX idx_product_images_primary 
ON product_images(product_id) 
WHERE is_primary = true;

COMMENT ON TABLE product_images IS 'Product images stored in S3';
COMMENT ON COLUMN product_images.s3_key IS 'S3 object key for management and deletion';

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    -- variant details
    label TEXT NOT NULL, --"250kg", "500kg", "1kg"
    
    quantity_value NUMERIC NOT NULL CHECK (quantity_value >= 0), --250, 500, 1
    quantity_unit TEXT NOT NULL, --"kg"
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0), -- exact money, max: 99999999.99

    stock_quantity INTEGER,
    low_stock_threshold INTEGER,
    
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

    UNIQUE(product_id, label)

);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_is_available ON product_variants(is_available);

COMMENT ON TABLE product_variants IS 'Different packaging/quantity options for products';
COMMENT ON COLUMN product_variants.label IS 'Display name for variant (e.g., "1kg Pack", "Bulk 25kg")';



-- =============================================
-- PRODUCT APPROVAL HISTORY
-- =============================================

CREATE TABLE product_approval_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    action TEXT NOT NULL CHECK (action IN ('SUBMITTED', 'APPROVED', 'REJECTED', 'RESUBMITTED')),
    performed_by_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT,
    notes TEXT,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_product_approval_history_product_id ON product_approval_history(product_id);
CREATE INDEX idx_product_approval_history_created_at ON product_approval_history(created_at DESC);

-- =============================================
-- TRIGGERS
-- =============================================

-- 1. Auto-log product submission
CREATE OR REPLACE FUNCTION log_product_submission()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO product_approval_history (product_id, action, performed_by_id, notes)
        SELECT NEW.id, 'SUBMITTED', c.owner_id, 'Initial product submission'
        FROM companies c WHERE c.id = NEW.company_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_log_product_submission
AFTER INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION log_product_submission();

-- 2. Prevent creating products for unapproved companies
CREATE OR REPLACE FUNCTION check_company_approved_for_product()
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
        RAISE EXCEPTION 'Cannot create product for unapproved company';
    END IF;
    
    IF NOT comp_active THEN
        RAISE EXCEPTION 'Cannot create product for inactive company';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_company_approved_for_product
BEFORE INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION check_company_approved_for_product();

-- 3. Update product updated_at on variant changes
CREATE OR REPLACE FUNCTION update_product_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products 
    SET updated_at = NOW() 
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_product_timestamp_on_variant
AFTER INSERT OR UPDATE ON product_variants
FOR EACH ROW
EXECUTE FUNCTION update_product_timestamp();

CREATE TRIGGER trigger_update_product_timestamp_on_image
AFTER INSERT OR UPDATE OR DELETE ON product_images
FOR EACH ROW
EXECUTE FUNCTION update_product_timestamp();