-- Extension để dùng UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-------------------------------
-- Bảng roles
-------------------------------
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-------------------------------
-- Bảng categories
-------------------------------
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-------------------------------
-- Bảng sub_categories
-------------------------------
CREATE TABLE sub_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng users
-------------------------------
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    image VARCHAR(255),
    subscribe INT DEFAULT 0,
    is_active BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'none' CHECK (status IN ('none', 'approved', 'rejected')),
    description TEXT,
    facebook VARCHAR(255),
    twitter VARCHAR(255),
    linkedin VARCHAR(255),
    youtube VARCHAR(255),
    rejection_reason TEXT,
    category_id UUID,
    sub_category_id UUID,
    role_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_user_sub_category FOREIGN KEY (sub_category_id) REFERENCES sub_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL
);

-------------------------------
-- Bảng courses
-------------------------------
CREATE TABLE courses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    short_description TEXT,
    requirements TEXT,
    student_learn TEXT,
    image VARCHAR(255),
    intro_video VARCHAR(255),
    price NUMERIC(12,2) DEFAULT 0,
    slug VARCHAR(255),
    sub_category_id UUID,
    user_id UUID,
    total_sold INT DEFAULT 0,
    total_view INT DEFAULT 0,
    total_rating INT DEFAULT 0,
    total_enrolled INT DEFAULT 0,
    require_log_in BOOLEAN DEFAULT FALSE,
    require_enroll BOOLEAN DEFAULT FALSE,
    discount NUMERIC(12,2) DEFAULT 0,
    is_deleted BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('pending','approved','rejected','draft')),
    rejection_reason TEXT,
    approved_at TIMESTAMP,
    is_bestseller BOOLEAN DEFAULT FALSE,
    likes INT DEFAULT 0,
    dislikes INT DEFAULT 0,
    levels JSONB DEFAULT '[]',
    captions JSONB DEFAULT '[]',
    language JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_course_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_course_sub_category FOREIGN KEY (sub_category_id) REFERENCES sub_categories(id) ON DELETE SET NULL
);

-------------------------------
-- Bảng sections
-------------------------------
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    course_id UUID NOT NULL,
    slug VARCHAR(255),
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_section_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng lectures
-------------------------------
CREATE TABLE lectures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL,
    title VARCHAR(255),
    description TEXT,
    free_preview BOOLEAN DEFAULT FALSE,
    video_url VARCHAR(255),
    video_poster_url VARCHAR(255),
    duration VARCHAR(50),
    uploaded_files JSONB DEFAULT '[]',
    slug VARCHAR(255),
    is_deleted BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'not learned',
    type VARCHAR(50) DEFAULT 'lecture',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_lecture_section FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng assignments
-------------------------------
CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL,
    title VARCHAR(255),
    description TEXT,
    time_duration INT,
    total_number INT,
    min_pass_number INT,
    upload_limit INT,
    max_attachment_size INT,
    slug VARCHAR(255),
    uploaded_files JSONB DEFAULT '[]',
    is_deleted BOOLEAN DEFAULT FALSE,
    type VARCHAR(50) DEFAULT 'assignment',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_assignment_section FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng quizzes
-------------------------------
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255),
    description TEXT,
    section_id UUID NOT NULL,
    time_limit INT,
    quiz_gradable BOOLEAN DEFAULT FALSE,
    passing_score NUMERIC(5,2),
    number_of_questions INT,
    show_time BOOLEAN DEFAULT FALSE,
    question_limit INT,
    slug VARCHAR(255),
    is_deleted BOOLEAN DEFAULT FALSE,
    questions JSONB NOT NULL,
    type VARCHAR(50) DEFAULT 'quiz',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_quiz_section FOREIGN KEY (section_id) REFERENCES sections(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng carts
-------------------------------
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_cart_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng cart_items
-------------------------------
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID NOT NULL,
    course_id UUID NOT NULL,
    title VARCHAR(255),
    sub_category VARCHAR(255),
    author VARCHAR(255),
    price NUMERIC(12,2) DEFAULT 0,
    image VARCHAR(255),
    slug VARCHAR(255),
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_cartitem_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE
);

-------------------------------
-- Trung gian: liked_courses
-------------------------------
CREATE TABLE liked_courses (
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    PRIMARY KEY(user_id, course_id),
    CONSTRAINT fk_liked_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_liked_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Trung gian: disliked_courses
-------------------------------
CREATE TABLE disliked_courses (
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    PRIMARY KEY(user_id, course_id),
    CONSTRAINT fk_disliked_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_disliked_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Trung gian: subscribe_channels
-------------------------------
CREATE TABLE subscribe_channels (
    user_id UUID NOT NULL,
    channel_id UUID NOT NULL,
    PRIMARY KEY(user_id, channel_id),
    CONSTRAINT fk_subscribe_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscribe_channel FOREIGN KEY (channel_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng quiz_results
-------------------------------
CREATE TABLE quiz_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID NOT NULL,
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    score NUMERIC(5,2),
    result VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_quiz_result_quiz FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_result_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_quiz_result_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng certificates
-------------------------------
CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255),
    email_address VARCHAR(255),
    phone_number VARCHAR(20),
    sub_category_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_certificate_sub_category FOREIGN KEY (sub_category_id) REFERENCES sub_categories(id) ON DELETE SET NULL
);

-------------------------------
-- Bảng progress
-------------------------------
CREATE TABLE progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_progress JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_progress_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng orders
-------------------------------
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    quantity INT DEFAULT 0,
    total_amount NUMERIC(12,2) DEFAULT 0,
    payment_status VARCHAR(50) DEFAULT 'pending', -- pending, paid, failed
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_order_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng order_details
-------------------------------
CREATE TABLE order_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    course_id UUID NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_order_detail_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_order_detail_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng purchased
-------------------------------
CREATE TABLE purchased (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_ids JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_purchased_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng reviews
-------------------------------
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_review_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng payment_providers
-------------------------------
CREATE TABLE payment_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    provider_code VARCHAR(50) UNIQUE,
    config JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-------------------------------
-- Bảng payments
-------------------------------
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    provider_id UUID NOT NULL,
    transaction_id VARCHAR(255),
    amount NUMERIC(12,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'pending',   -- pending, success, failed, refunded
    payment_method VARCHAR(50),
    payment_response JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_provider FOREIGN KEY (provider_id) REFERENCES payment_providers(id) ON DELETE SET NULL
);

-------------------------------
-- Bảng payment_logs
-------------------------------
CREATE TABLE payment_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    event_type VARCHAR(100),
    payload JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_payment_log FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng wallets
-------------------------------
CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE,
    balance NUMERIC(12,2) DEFAULT 0,
    total_deposit NUMERIC(12,2) DEFAULT 0,
    total_withdraw NUMERIC(12,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_wallet_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-------------------------------
-- Bảng wallet_transactions
-------------------------------
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    order_id UUID,
    type VARCHAR(10) NOT NULL,          -- credit/debit
    amount NUMERIC(12,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_wallet_transaction FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    CONSTRAINT fk_wallet_transaction_order FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
);

-------------------------------
-- Bảng withdrawal_requests
-------------------------------
CREATE TABLE withdrawal_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    requested_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_withdrawal_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
);
