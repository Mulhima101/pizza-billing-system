-- Create database
CREATE DATABASE pizza_billing;

-- Use the database
\c pizza_billing;

-- Categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Items table
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Customers table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(15),
    email VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Invoices table
CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(20) UNIQUE NOT NULL,
    customer_id INTEGER REFERENCES customers(id),
    subtotal DECIMAL(10,2) NOT NULL,
    tax_rate DECIMAL(5,2) DEFAULT 10.00,
    tax_amount DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Invoice items table
CREATE TABLE invoice_items (
    id SERIAL PRIMARY KEY,
    invoice_id INTEGER REFERENCES invoices(id) ON DELETE CASCADE,
    item_id INTEGER REFERENCES items(id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL
);

-- Insert sample data
INSERT INTO categories (name, description) VALUES
('Pizza', 'Various pizza types'),
('Toppings', 'Pizza toppings'),
('Beverages', 'Drinks and beverages'),
('Sides', 'Side dishes');

INSERT INTO items (name, description, price, category_id) VALUES
('Margherita Pizza', 'Classic pizza with tomato and mozzarella', 12.99, 1),
('Pepperoni Pizza', 'Pizza with pepperoni', 15.99, 1),
('Extra Cheese', 'Additional cheese topping', 2.50, 2),
('Mushrooms', 'Fresh mushroom topping', 2.00, 2),
('Coca Cola', 'Cold beverage', 2.99, 3),
('Garlic Bread', 'Side dish', 4.99, 4);

INSERT INTO customers (name, phone, email, address) VALUES
('John Doe', '+1234567890', 'john@email.com', '123 Main St'),
('Jane Smith', '+0987654321', 'jane@email.com', '456 Oak Ave');