-- Create categories table if not exists
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id BIGINT,
    level INTEGER NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Create products table if not exists
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id BIGINT NOT NULL,
    seller_id BIGINT NOT NULL,
    starting_price DOUBLE PRECISION NOT NULL,
    current_price DOUBLE PRECISION,
    buy_now_price DOUBLE PRECISION,
    step_price DOUBLE PRECISION NOT NULL,
    status VARCHAR(255) NOT NULL,
    thumbnail_url TEXT,
    auto_extend BOOLEAN NOT NULL DEFAULT false,
    current_bidder BIGINT,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
    CONSTRAINT products_status_check CHECK (status IN ('ACTIVE', 'FINISHED', 'PENDING', 'REJECTED'))
);

-- Create indices for better performance
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_level ON categories(level);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);
CREATE INDEX IF NOT EXISTS idx_products_end_at ON products(end_at);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);

-- Insert sample categories
INSERT INTO categories (name, slug, description, level) VALUES
    ('Electronics', 'electronics', 'Electronic devices and gadgets', 1),
    ('Fashion', 'fashion', 'Clothing and accessories', 1),
    ('Home & Garden', 'home-garden', 'Home improvement and gardening', 1)
ON CONFLICT (slug) DO NOTHING;

-- Insert sample products
INSERT INTO products (name, description, category_id, seller_id, starting_price, current_price, step_price, status, end_at) VALUES
    ('iPhone 15 Pro Max', 'Latest Apple smartphone with advanced camera system', 1, 1, 20000000, 20000000, 500000, 'ACTIVE', NOW() + INTERVAL '7 days'),
    ('Samsung Galaxy S24 Ultra', 'Premium Android smartphone with S Pen', 1, 1, 18000000, 18000000, 500000, 'ACTIVE', NOW() + INTERVAL '5 days'),
    ('MacBook Pro M3', 'Professional laptop for creators', 1, 2, 35000000, 35000000, 1000000, 'ACTIVE', NOW() + INTERVAL '10 days')
ON CONFLICT DO NOTHING;
