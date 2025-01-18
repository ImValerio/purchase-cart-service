-- Active: 1736963717585@@127.0.0.1@5432@purchase_cart_db
CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    vat_rate DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "order" (
    id SERIAL PRIMARY KEY,
    total_price DECIMAL(10, 2),
    total_vat DECIMAL(10, 2) ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_item (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES "order"(id) ON DELETE CASCADE,
    product_id INT REFERENCES product(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    vat DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

