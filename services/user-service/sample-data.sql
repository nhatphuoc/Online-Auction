INSERT INTO users (
    full_name,
    email,
    birth_day,
    role,
    password,
    total_number_reviews,
    total_number_good_reviews,
    email_verified,
    is_seller_request_sent
) VALUES
-- ================= BIDDERS =================
('Alice Nguyen', 'alice.bidder@example.com', '1998-03-12', 'ROLE_BIDDER', NULL, 5, 4, true, false),
('Bob Tran', 'bob.bidder@example.com', '1997-07-25', 'ROLE_BIDDER', NULL, 2, 2, true, false),
('Charlie Le', 'charlie.bidder@example.com', '2000-01-18', 'ROLE_BIDDER', NULL, 0, 0, true, false),

-- ================= SELLERS =================
('David Pham', 'david.seller@example.com', '1992-11-05', 'ROLE_SELLER', NULL, 20, 18, true, true),
('Emma Vo', 'emma.seller@example.com', '1990-06-30', 'ROLE_SELLER', NULL, 15, 14, true, true),

-- ================= ADMIN =================
('Admin User', 'admin@example.com', '1985-01-01', 'ROLE_ADMIN', NULL, 0, 0, true, false);
