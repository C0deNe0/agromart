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
    
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    approval_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
    --setting null when admin is deleted
    approval_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(owner_id, name)
);

CREATE INDEX idx_company_owner_id ON companies(owner_id);
CREATE INDEX idx_company_is_approved ON companies(is_approved);
CREATE INDEX idx_company_is_active ON companies(is_active);



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
