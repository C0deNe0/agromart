-- UP: 00006_create_db_triggers

-- =============================================
-- 1. COMPANY TRIGGERS
-- =============================================

-- FUNCTION: auto log submission history
CREATE OR REPLACE FUNCTION log_company_submission()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO company_approval_history (company_id, action, performed_by_id, notes)
        VALUES (NEW.id, 'SUBMITTED', NEW.owner_id, 'Initial submission');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation
DROP TRIGGER IF EXISTS trigger_log_company_submission ON companies;
CREATE TRIGGER trigger_log_company_submission
AFTER INSERT ON companies
FOR EACH ROW
EXECUTE FUNCTION log_company_submission();


-- FUNCTION: update follower count
CREATE OR REPLACE FUNCTION update_company_follower_count()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE companies
        SET follower_count = follower_count + 1
        WHERE id = NEW.company_id;
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE companies
        SET follower_count = GREATEST(follower_count - 1, 0)
        WHERE id = OLD.company_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation
DROP TRIGGER IF EXISTS trigger_update_company_follower_count ON company_followers;
CREATE TRIGGER trigger_update_company_follower_count
AFTER INSERT OR DELETE ON company_followers
FOR EACH ROW
EXECUTE FUNCTION update_company_follower_count();


-- FUNCTION: Prevent following non-approved companies
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
        RAISE EXCEPTION 'Cannot follow unapproved company';
    END IF;

    IF NOT comp_active THEN
        RAISE EXCEPTION 'Cannot follow inactive company';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation
DROP TRIGGER IF EXISTS trigger_check_company_followable ON company_followers;
CREATE TRIGGER trigger_check_company_followable
BEFORE INSERT ON company_followers
FOR EACH ROW
EXECUTE FUNCTION check_company_followable();


-- =============================================
-- 2. PRODUCT TRIGGERS
-- =============================================

-- FUNCTION: Auto-log product submission
-- NOTE: Corrected the performed_by_id lookup here to use the company owner, 
--       as the original SELECT logic was more accurate.
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
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation
DROP TRIGGER IF EXISTS trigger_log_product_submission ON products;
CREATE TRIGGER trigger_log_product_submission
AFTER INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION log_product_submission();


-- FUNCTION: Prevent creating products for unapproved companies
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
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation
DROP TRIGGER IF EXISTS trigger_check_company_approved_for_product ON products;
CREATE TRIGGER trigger_check_company_approved_for_product
BEFORE INSERT ON products
FOR EACH ROW
EXECUTE FUNCTION check_company_approved_for_product();

-- FUNCTION: Update product updated_at on variant/image changes
CREATE OR REPLACE FUNCTION update_product_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE products 
    SET updated_at = NOW() 
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; -- <-- SEMICOLON ADDED HERE

-- CREATE TRIGGER: FIX: Drop trigger before creation (for variants)
DROP TRIGGER IF EXISTS trigger_update_product_timestamp_on_variant ON product_variants;
CREATE TRIGGER trigger_update_product_timestamp_on_variant
AFTER INSERT OR UPDATE ON product_variants
FOR EACH ROW
EXECUTE FUNCTION update_product_timestamp();

-- CREATE TRIGGER: FIX: Drop trigger before creation (for images)
DROP TRIGGER IF EXISTS trigger_update_product_timestamp_on_image ON product_images;
CREATE TRIGGER trigger_update_product_timestamp_on_image
AFTER INSERT OR UPDATE OR DELETE ON product_images
FOR EACH ROW
EXECUTE FUNCTION update_product_timestamp();