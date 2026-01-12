-- UP: 00007_insert_seed_data (INLINE CONSTANTS - UUIDS FIXED)

-- =============================================
-- 1. UUID CONSTANTS (FIXED TO BE VALID HEX)
-- =============================================
-- User UUIDs (Prefix 'a' changed to 'b'):
-- Admin: b0000000-0000-4000-8000-000000000001
-- Owner: b0000000-0000-4000-8000-000000000002
-- Buyer: b0000000-0000-4000-8000-000000000003
-- Company UUID (Prefix 'c' changed to '1'): 10000000-0000-4000-8000-000000000001
-- Category UUIDs (Prefix 'd' changed to '2'): 20000000-0000-4000-8000-000000000001, 20000000-0000-4000-8000-000000000002, 20000000-0000-4000-8000-000000000003
-- Product UUIDs (Prefix 'p' changed to '3'): 30000000-0000-4000-8000-000000000001, 30000000-0000-4000-8000-000000000002

-- Hash for "password123": $2a$10$v0/fFk8o5q4tX.QxI2/0yOB33e6F9rA3.09xPq9pY3S8H4j2y/nJ2

-- =============================================
-- 2. INSERT USERS
-- =============================================

INSERT INTO users (id, email, name, role, email_verified) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'admin@agromart.com', 'Super Admin', 'ADMIN', TRUE),
    ('b0000000-0000-4000-8000-000000000002', 'john.doe@farm.com', 'John Doe (Owner)', 'USER', TRUE),
    ('b0000000-0000-4000-8000-000000000003', 'jane.smith@buyer.com', 'Jane Smith (Buyer)', 'USER', TRUE);


-- =============================================
-- 2.1. INSERT USER AUTH METHODS (Login Credentials)
-- =============================================

-- ADMIN Credentials
INSERT INTO user_auth_methods (user_id, auth_provider, password_hash) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'LOCAL', '$2a$10$v0/fFk8o5q4tX.QxI2/0yOB33e6F9rA3.09xPq9pY3S8H4j2y/nJ2');

-- OWNER/NORMAL USER Credentials
INSERT INTO user_auth_methods (user_id, auth_provider, password_hash) VALUES
    ('b0000000-0000-4000-8000-000000000002', 'LOCAL', '$2a$10$v0/fFk8o5q4tX.QxI2/0yOB33e6F9rA3.09xPq9pY3S8H4j2y/nJ2');

-- Buyer/Normal User Credentials
INSERT INTO user_auth_methods (user_id, auth_provider, password_hash) VALUES
    ('b0000000-0000-4000-8000-000000000003', 'LOCAL', '$2a$10$v0/fFk8o5q4tX.QxI2/0yOB33e6F9rA3.09xPq9pY3S8H4j2y/nJ2');


-- =============================================
-- 3. INSERT COMPANIES
-- =============================================

INSERT INTO companies (id, owner_id, name, description, city, approval_status, follower_count) VALUES
    ('10000000-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000002', 'Green Valley Organics', 'Leading supplier of certified organic produce.', 'Bangalore', 'APPROVED', 1);

INSERT INTO company_approval_history (company_id, action, performed_by_id, notes) VALUES
    ('10000000-0000-4000-8000-000000000001', 'SUBMITTED', 'b0000000-0000-4000-8000-000000000002', 'Initial submission'),
    ('10000000-0000-4000-8000-000000000001', 'APPROVED', 'b0000000-0000-4000-8000-000000000001', 'Company approved during seed data insertion');

INSERT INTO companies (owner_id, name, description, city, approval_status) VALUES
    ('b0000000-0000-4000-8000-000000000003', 'Future Farmers Co.', 'A new cooperative awaiting admin approval.', 'Mysore', 'PENDING');


-- =============================================
-- 4. INSERT COMPANY RELATIONS
-- =============================================

-- Make Jane Smith follow Green Valley Organics
INSERT INTO company_followers (company_id, user_id) VALUES
    ('10000000-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000003');


-- =============================================
-- 5. INSERT CATEGORIES
-- =============================================

INSERT INTO categories (id, name, slug) VALUES
    ('20000000-0000-4000-8000-000000000001', 'Fresh Fruits', 'fresh-fruits'),
    ('20000000-0000-4000-8000-000000000002', 'Root Vegetables', 'root-vegetables'),
    ('20000000-0000-4000-8000-000000000003', 'Cereals and Grains', 'cereals-and-grains');


-- =============================================
-- 6. INSERT PRODUCTS
-- =============================================

-- INSERT INTO products (id, company_id, category_id, name, description, unit, origin, base_price, approval_status) VALUES
--     ('30000000-0000-4000-8000-000000000001', '10000000-0000-4000-8000-000000000001', '20000000-0000-4000-8000-000000000001', 'Organic Fuji Apple', 'Sweet and crisp apples from Himachal Pradesh.', 'kg', 'India', 150.00, 'APPROVED'),
--     ('30000000-0000-4000-8000-000000000002', '10000000-0000-4000-8000-000000000001', '20000000-0000-4000-8000-000000000003', 'High-Protein Wheat Grain', 'Premium quality wheat, ideal for baking.', 'quintal', 'India', 3500.00, 'APPROVED');

-- INSERT INTO products (company_id, category_id, name, description, unit, origin, base_price, approval_status) VALUES
--     ('10000000-0000-4000-8000-000000000001', '20000000-0000-4000-8000-000000000002', 'Farm Fresh Carrot', 'Sweet and crunchy carrots, currently under review.', 'kg', 'India', 60.00, 'PENDING');


-- =============================================
-- 7. INSERT PRODUCT RELATIONS
-- =============================================

-- Product Variants (for Apples)
-- INSERT INTO product_variants (product_id, label, quantity_value, quantity_unit, price, stock_quantity) VALUES
--     ('30000000-0000-4000-8000-000000000001', '1kg Pack', 1.0, 'kg', 180.00, 150),
--     ('30000000-0000-4000-8000-000000000001', '5kg Box', 5.0, 'kg', 850.00, 30);

-- Product Images (for Apples)
INSERT INTO product_images (product_id, image_url, s3_key, is_primary) VALUES
    ('30000000-0000-4000-8000-000000000001', 'https://s3.aws.com/product/apples-main.jpg', 'product/apples/main.jpg', TRUE),
    ('30000000-0000-4000-8000-000000000001', 'https://s3.aws.com/product/apples-side.jpg', 'product/apples/side.jpg', FALSE);

-- Favorites
INSERT INTO favorites (user_id, product_id) VALUES
    ('b0000000-0000-4000-8000-000000000003', '30000000-0000-4000-8000-000000000001');


-- DOWN: 00007_insert_seed_data (UUIDS FIXED)

-- =============================================
-- 1. DELETE DATA (Order must respect Foreign Keys)
-- =============================================

-- 1. Tables referencing Products
DELETE FROM product_variants WHERE product_id = '30000000-0000-4000-8000-000000000001';
DELETE FROM product_images WHERE product_id = '30000000-0000-4000-8000-000000000001';
DELETE FROM favorites WHERE product_id IN (
    '30000000-0000-4000-8000-000000000001', 
    '30000000-0000-4000-8000-000000000002'
);

-- 2. Products
DELETE FROM products WHERE id IN (
    '30000000-0000-4000-8000-000000000001', 
    '30000000-0000-4000-8000-000000000002'
);
-- Delete the pending product (using unique attributes)
DELETE FROM products WHERE company_id = '10000000-0000-4000-8000-000000000001' AND name = 'Farm Fresh Carrot';

-- 3. Tables referencing Companies
DELETE FROM company_followers WHERE company_id = '10000000-0000-4000-8000-000000000001';
DELETE FROM company_approval_history WHERE company_id = '10000000-0000-4000-8000-000000000001';

-- 4. Companies
DELETE FROM companies WHERE id = '10000000-0000-4000-8000-000000000001';
-- Delete the pending company (using unique attributes)
DELETE FROM companies WHERE owner_id = 'b0000000-0000-4000-8000-000000000003' AND name = 'Future Farmers Co.';

-- 5. Categories
DELETE FROM categories WHERE id IN (
    '20000000-0000-4000-8000-000000000001', 
    '20000000-0000-4000-8000-000000000002', 
    '20000000-0000-4000-8000-000000000003'
);

-- 6. Tables referencing Users
DELETE FROM user_auth_methods WHERE user_id IN (
    'b0000000-0000-4000-8000-000000000001', 
    'b0000000-0000-4000-8000-000000000002', 
    'b0000000-0000-4000-8000-000000000003'
);

-- 7. Users
DELETE FROM users WHERE id IN (
    'b0000000-0000-4000-8000-000000000001', 
    'b0000000-0000-4000-8000-000000000002', 
    'b0000000-0000-4000-8000-000000000003'
);