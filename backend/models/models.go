package models

import (
    "time"
)

type Category struct {
    ID          int       `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Description string    `json:"description" db:"description"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Item struct {
    ID          int     `json:"id" db:"id"`
    Name        string  `json:"name" db:"name"`
    Description string  `json:"description" db:"description"`
    Price       float64 `json:"price" db:"price"`
    CategoryID  int     `json:"category_id" db:"category_id"`
    IsActive    bool    `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    CategoryName string  `json:"category_name,omitempty"`
}

type Customer struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Phone     string    `json:"phone" db:"phone"`
    Email     string    `json:"email" db:"email"`
    Address   string    `json:"address" db:"address"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Invoice struct {
    ID            int           `json:"id" db:"id"`
    InvoiceNumber string        `json:"invoice_number" db:"invoice_number"`
    CustomerID    int           `json:"customer_id" db:"customer_id"`
    Subtotal      float64       `json:"subtotal" db:"subtotal"`
    TaxRate       float64       `json:"tax_rate" db:"tax_rate"`
    TaxAmount     float64       `json:"tax_amount" db:"tax_amount"`
    TotalAmount   float64       `json:"total_amount" db:"total_amount"`
    Status        string        `json:"status" db:"status"`
    CreatedAt     time.Time     `json:"created_at" db:"created_at"`
    Customer      *Customer     `json:"customer,omitempty"`
    Items         []InvoiceItem `json:"items,omitempty"`
}

type InvoiceItem struct {
    ID         int     `json:"id" db:"id"`
    InvoiceID  int     `json:"invoice_id" db:"invoice_id"`
    ItemID     int     `json:"item_id" db:"item_id"`
    Quantity   int     `json:"quantity" db:"quantity"`
    UnitPrice  float64 `json:"unit_price" db:"unit_price"`
    TotalPrice float64 `json:"total_price" db:"total_price"`
    ItemName   string  `json:"item_name,omitempty"`
}

type CreateInvoiceRequest struct {
    CustomerID int                     `json:"customer_id"`
    Items      []CreateInvoiceItemRequest `json:"items"`
}

type CreateInvoiceItemRequest struct {
    ItemID   int `json:"item_id"`
    Quantity int `json:"quantity"`
}