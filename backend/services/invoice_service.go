package services

import (
    "database/sql"
    "fmt"
    "pizza-billing-backend/models"
    "time"
)

type InvoiceService struct {
    db *sql.DB
}

func NewInvoiceService(db *sql.DB) *InvoiceService {
    return &InvoiceService{db: db}
}

func (s *InvoiceService) GetAllInvoices() ([]models.Invoice, error) {
    query := `
        SELECT i.id, i.invoice_number, i.customer_id, i.subtotal, i.tax_rate, 
               i.tax_amount, i.total_amount, i.status, i.created_at,
               c.name, c.phone, c.email, c.address
        FROM invoices i
        LEFT JOIN customers c ON i.customer_id = c.id
        ORDER BY i.created_at DESC
    `
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var invoices []models.Invoice
    for rows.Next() {
        var invoice models.Invoice
        var customer models.Customer
        
        err := rows.Scan(&invoice.ID, &invoice.InvoiceNumber, &invoice.CustomerID,
                        &invoice.Subtotal, &invoice.TaxRate, &invoice.TaxAmount,
                        &invoice.TotalAmount, &invoice.Status, &invoice.CreatedAt,
                        &customer.Name, &customer.Phone, &customer.Email, &customer.Address)
        if err != nil {
            return nil, err
        }
        
        customer.ID = invoice.CustomerID
        invoice.Customer = &customer
        
        // Get invoice items
        items, _ := s.getInvoiceItems(invoice.ID)
        invoice.Items = items
        
        invoices = append(invoices, invoice)
    }

    return invoices, nil
}

func (s *InvoiceService) GetInvoiceByID(id int) (*models.Invoice, error) {
    query := `
        SELECT i.id, i.invoice_number, i.customer_id, i.subtotal, i.tax_rate,
               i.tax_amount, i.total_amount, i.status, i.created_at,
               c.name, c.phone, c.email, c.address
        FROM invoices i
        LEFT JOIN customers c ON i.customer_id = c.id
        WHERE i.id = $1
    `
    
    var invoice models.Invoice
    var customer models.Customer
    
    err := s.db.QueryRow(query, id).Scan(&invoice.ID, &invoice.InvoiceNumber,
                                        &invoice.CustomerID, &invoice.Subtotal,
                                        &invoice.TaxRate, &invoice.TaxAmount,
                                        &invoice.TotalAmount, &invoice.Status,
                                        &invoice.CreatedAt, &customer.Name,
                                        &customer.Phone, &customer.Email, &customer.Address)
    if err != nil {
        return nil, err
    }
    
    customer.ID = invoice.CustomerID
    invoice.Customer = &customer
    
    // Get invoice items
    items, err := s.getInvoiceItems(invoice.ID)
    if err != nil {
        return nil, err
    }
    invoice.Items = items

    return &invoice, nil
}

func (s *InvoiceService) CreateInvoice(req *models.CreateInvoiceRequest) (*models.Invoice, error) {
    tx, err := s.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Generate invoice number
    invoiceNumber := s.generateInvoiceNumber()
    
    // Calculate totals
    var subtotal float64
    for _, item := range req.Items {
        // Get item price
        var price float64
        err := tx.QueryRow("SELECT price FROM items WHERE id = $1", item.ItemID).Scan(&price)
        if err != nil {
            return nil, err
        }
        subtotal += price * float64(item.Quantity)
    }
    
    taxRate := 10.0 // 10% tax
    taxAmount := subtotal * taxRate / 100
    totalAmount := subtotal + taxAmount

    // Create invoice
    var invoiceID int
    query := `
        INSERT INTO invoices (invoice_number, customer_id, subtotal, tax_rate, tax_amount, total_amount, status)
        VALUES ($1, $2, $3, $4, $5, $6, 'pending')
        RETURNING id
    `
    
    err = tx.QueryRow(query, invoiceNumber, req.CustomerID, subtotal, taxRate, taxAmount, totalAmount).Scan(&invoiceID)
    if err != nil {
        return nil, err
    }

    // Create invoice items
    for _, item := range req.Items {
        var price float64
        err := tx.QueryRow("SELECT price FROM items WHERE id = $1", item.ItemID).Scan(&price)
        if err != nil {
            return nil, err
        }
        
        totalPrice := price * float64(item.Quantity)
        
        _, err = tx.Exec(`
            INSERT INTO invoice_items (invoice_id, item_id, quantity, unit_price, total_price)
            VALUES ($1, $2, $3, $4, $5)
        `, invoiceID, item.ItemID, item.Quantity, price, totalPrice)
        if err != nil {
            return nil, err
        }
    }

    err = tx.Commit()
    if err != nil {
        return nil, err
    }

    return s.GetInvoiceByID(invoiceID)
}

func (s *InvoiceService) getInvoiceItems(invoiceID int) ([]models.InvoiceItem, error) {
    query := `
        SELECT ii.id, ii.invoice_id, ii.item_id, ii.quantity, ii.unit_price, ii.total_price, i.name
        FROM invoice_items ii
        LEFT JOIN items i ON ii.item_id = i.id
        WHERE ii.invoice_id = $1
    `
    
    rows, err := s.db.Query(query, invoiceID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.InvoiceItem
    for rows.Next() {
        var item models.InvoiceItem
        err := rows.Scan(&item.ID, &item.InvoiceID, &item.ItemID, &item.Quantity,
                        &item.UnitPrice, &item.TotalPrice, &item.ItemName)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    return items, nil
}

func (s *InvoiceService) generateInvoiceNumber() string {
    return fmt.Sprintf("INV-%d", time.Now().Unix())
}

func (s *InvoiceService) GetAllCustomers() ([]models.Customer, error) {
    query := `SELECT id, name, phone, email, address, created_at FROM customers ORDER BY name`
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []models.Customer
    for rows.Next() {
        var customer models.Customer
        err := rows.Scan(&customer.ID, &customer.Name, &customer.Phone,
                        &customer.Email, &customer.Address, &customer.CreatedAt)
        if err != nil {
            return nil, err
        }
        customers = append(customers, customer)
    }

    return customers, nil
}

func (s *InvoiceService) CreateCustomer(customer *models.Customer) (*models.Customer, error) {
    query := `
        INSERT INTO customers (name, phone, email, address)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
    
    err := s.db.QueryRow(query, customer.Name, customer.Phone, customer.Email, customer.Address).
        Scan(&customer.ID, &customer.CreatedAt)
    if err != nil {
        return nil, err
    }

    return customer, nil
}