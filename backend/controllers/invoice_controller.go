package controllers

import (
    "net/http"
    "pizza-billing-backend/models"
    "pizza-billing-backend/services"
    "strconv"

    "github.com/gin-gonic/gin"
)

type InvoiceController struct {
    invoiceService *services.InvoiceService
}

func NewInvoiceController(invoiceService *services.InvoiceService) *InvoiceController {
    return &InvoiceController{invoiceService: invoiceService}
}

func (ic *InvoiceController) GetInvoices(c *gin.Context) {
    invoices, err := ic.invoiceService.GetAllInvoices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, invoices)
}

func (ic *InvoiceController) GetInvoice(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
        return
    }

    invoice, err := ic.invoiceService.GetInvoiceByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }
    c.JSON(http.StatusOK, invoice)
}

func (ic *InvoiceController) CreateInvoice(c *gin.Context) {
    var req models.CreateInvoiceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    invoice, err := ic.invoiceService.CreateInvoice(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, invoice)
}

func (ic *InvoiceController) GetCustomers(c *gin.Context) {
    customers, err := ic.invoiceService.GetAllCustomers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, customers)
}

func (ic *InvoiceController) CreateCustomer(c *gin.Context) {
    var customer models.Customer
    if err := c.ShouldBindJSON(&customer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdCustomer, err := ic.invoiceService.CreateCustomer(&customer)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, createdCustomer)
}