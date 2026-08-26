-- PostgreSQL database dump for Product Inventory with Auth
-- Dumped on 2026-08-25

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- Clean existing tables
DROP TABLE IF EXISTS order_items CASCADE;
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Recreate schema
CREATE TABLE public.products (
    id BIGSERIAL PRIMARY KEY,
    name character varying(255) NOT NULL,
    description text NOT NULL,
    price bigint NOT NULL,
    stock_quantity integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT chk_stock_quantity_non_negative CHECK ((stock_quantity >= 0))
);

CREATE TABLE public.users (
    id BIGSERIAL PRIMARY KEY,
    username character varying(255) UNIQUE NOT NULL,
    password_hash character varying(255) NOT NULL,
    role character varying(50) DEFAULT 'customer'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE public.orders (
    id BIGSERIAL PRIMARY KEY,
    status character varying(50) DEFAULT 'PENDING'::character varying NOT NULL,
    total_amount bigint DEFAULT 0 NOT NULL,
    user_id bigint REFERENCES public.users(id) ON DELETE CASCADE,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    cancelled_at timestamp with time zone
);

CREATE TABLE public.order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id bigint NOT NULL REFERENCES public.orders(id) ON DELETE CASCADE,
    product_id bigint NOT NULL REFERENCES public.products(id) ON DELETE CASCADE,
    quantity integer NOT NULL,
    unit_price bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT order_items_quantity_check CHECK ((quantity > 0))
);

CREATE INDEX idx_products_name ON public.products USING btree (name);
CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);
CREATE INDEX idx_order_items_product_id ON public.order_items USING btree (product_id);

-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
INSERT INTO public.products (id, name, description, price, stock_quantity, created_at, updated_at) VALUES
(1, 'iPhone 15 Pro Max', 'Apple iPhone 15 Pro Max 256GB Natural Titanium', 119900, 15, now(), now()),
(2, 'Sony PlayStation 5', 'PlayStation 5 Console (Disc Version) with 825GB SSD', 49999, 10, now(), now()),
(3, 'Nintendo Switch OLED', 'Nintendo Switch Console OLED Model with Neon Blue & Neon Red Joy-Con', 34999, 20, now(), now()),
(4, 'Apple MacBook Air M3', 'MacBook Air 13-inch M3 Chip 8-Core CPU 8-Core GPU 8GB 256GB SSD', 109900, 8, now(), now()),
(5, 'Sony WH-1000XM5', 'Sony WH-1000XM5 Wireless Noise Canceling Headphones Black', 39800, 30, now(), now());

-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
-- Seed admin/adminpassword and customer/customerpassword
INSERT INTO public.users (id, username, password_hash, role, created_at, updated_at) VALUES
(1, 'admin', '$2a$10$gza/9WGwkqCDN2GIe6vGL.2FispTpJsMW/H0dPnW79So4m5m1ms/q', 'admin', now(), now()),
(2, 'customer', '$2a$10$tq1nzO418LKRHsBM5WSnqeIqa/B5nJwn1.XP52qGPguT1jTmwBEmC', 'customer', now(), now());

SELECT pg_catalog.setval('public.products_id_seq', 5, true);
SELECT pg_catalog.setval('public.users_id_seq', 2, true);
SELECT pg_catalog.setval('public.orders_id_seq', 1, false);
SELECT pg_catalog.setval('public.order_items_id_seq', 1, false);
