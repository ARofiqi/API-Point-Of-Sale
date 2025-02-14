CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    category_id INT(11) NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);