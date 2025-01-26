-- Drop tables if they exist
DROP TABLE IF EXISTS Orders;
DROP TABLE IF EXISTS Products;
DROP TABLE IF EXISTS Customers;

-- Create Customers table
CREATE TABLE Customers (
    customer_id INTEGER PRIMARY KEY,
    customer_name TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    join_date DATE NOT NULL
);

-- Create Products table
CREATE TABLE Products (
    product_id INTEGER PRIMARY KEY,
    product_name TEXT NOT NULL,
    price REAL NOT NULL  -- This works in PostgreSQL, MySQL, and SQLite
);

-- Create Orders table
CREATE TABLE Orders (
    order_id INTEGER PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    order_date DATE NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Customers(customer_id),
    FOREIGN KEY (product_id) REFERENCES Products(product_id)
);

-- Insert Customer data
INSERT INTO Customers (customer_id, customer_name, email, join_date) VALUES
    (1, 'Alice Johnson', 'alice@example.com', '2024-01-01'),
    (2, 'Bob Smith', 'bob@example.com', '2024-02-15'),
    (3, 'Charlie Davis', 'charlie@example.com', '2024-03-10');

-- Insert Product data
INSERT INTO Products (product_id, product_name, price) VALUES
    (101, 'Laptop', 1200.00),
    (102, 'Smartphone', 800.90),
    (103, 'Headphones', 150.00);

-- Insert Order data
INSERT INTO Orders (order_id, customer_id, product_id, order_date, quantity) VALUES
    (1001, 1, 101, '2024-06-01', 1),
    (1002, 2, 103, '2024-06-05', 2),
    (1003, 3, 102, '2024-06-07', 1),
    (1004, 1, 103, '2024-06-10', 3);