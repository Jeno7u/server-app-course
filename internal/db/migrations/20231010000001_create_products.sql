-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    price REAL NOT NULL,
    count INTEGER NOT NULL
);

INSERT INTO products (title, price, count) VALUES ('Laptop', 999.99, 10);
INSERT INTO products (title, price, count) VALUES ('Mouse', 25.50, 50);

-- +goose Down
DROP TABLE products;