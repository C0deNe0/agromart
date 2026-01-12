-- UP: 00004_create_product_tables

-- =============================================
-- PRODUCTS TABLE
-- =============================================

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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CHECK (length(name) >= 3)
);

CREATE INDEX idx_products_company_id ON products(company_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_approval_status ON products(approval_status);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_products_submitted_at ON products(submitted_at DESC);

COMMENT ON TABLE products IS 'Products require admin approval before being visible';
COMMENT ON COLUMN products.approval_status IS 'PENDING: Awaiting review | APPROVED: Visible to users | REJECTED: Needs modification';


-- =============================================
-- PRODUCT IMAGES
-- =============================================

CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,

    image_url TEXT NOT NULL,
    s3_key TEXT NOT NULL, 
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


-- =============================================
-- PRODUCT VARIANTS
-- =============================================

CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    
    -- variant details
    label TEXT NOT NULL, 
    
    quantity_value NUMERIC NOT NULL CHECK (quantity_value >= 0), 
    quantity_unit TEXT NOT NULL, 
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0), 

    stock_quantity INTEGER,
    low_stock_threshold INTEGER,
    
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

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