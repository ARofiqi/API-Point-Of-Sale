CREATE DATABASE db_aro_shop;

USE db_aro_shop;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(50) NOT NULL
);

CREATE TABLE transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATE NOT NULL DEFAULT (CURRENT_DATE),
    total DECIMAL(10,2) NOT NULL
);

CREATE TABLE transaction_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    sub_total DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);


# CREATE TABLE customers (
#     id INT AUTO_INCREMENT PRIMARY KEY,
#     name VARCHAR(100) NOT NULL,
#     no_telephone VARCHAR(10),
#     balance INT DEFAULT 0.00
# );

# CREATE TABLE transactions (
#     id INT AUTO_INCREMENT PRIMARY KEY,
#     product_id INT NOT NULL,
#     quantity INT NOT NULL,
#     total DECIMAL(10,2) NOT NULL,
#     FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
# );