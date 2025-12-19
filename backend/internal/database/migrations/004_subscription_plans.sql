CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    description TEXT,
    billing_cycle TEXT NOT NULL CHECK (billing_cycle in ('MONTHLY', 'YEARLY')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    --LIMTS to put 
    max_products INT,
    max_products_images INT,
    max_variants_per_product INT,


    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE company_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL
        REFERENCES companies(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL
        REFERENCES subscription_plans(id) ON DELETE CASCADE,
    
    status TEXT NOT NULL CHECK (status in ('ACTIVE','PAUSED', 'EXPIRED', 'CANCELLED')),

    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT company_active_subscription_unique UNIQUE (company_id, status) DEFERRABLE INITIALLY DEFERRED
);