package main

import (
    "log"
    "pizza-billing-backend/config"
    "pizza-billing-backend/controllers"
    "pizza-billing-backend/services"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Initialize database
    database, err := config.NewDatabase()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer database.Close()

    // Initialize services
    itemService := services.NewItemService(database.DB)
    invoiceService := services.NewInvoiceService(database.DB)

    // Initialize controllers
    itemController := controllers.NewItemController(itemService)
    invoiceController := controllers.NewInvoiceController(invoiceService)

    // Initialize router
    r := gin.Default()

    // CORS middleware
    r.Use(cors.New(cors.Config{
        AllowAllOrigins:  true,
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // API routes
    api := r.Group("/api")
    {
        // Item routes
        items := api.Group("/items")
        {
            items.GET("", itemController.GetItems)
            items.GET("/:id", itemController.GetItem)
            items.POST("", itemController.CreateItem)
            items.PUT("/:id", itemController.UpdateItem)
            items.DELETE("/:id", itemController.DeleteItem)
        }

        // Category routes
        api.GET("/categories", itemController.GetCategories)

        // Invoice routes
        invoices := api.Group("/invoices")
        {
            invoices.GET("", invoiceController.GetInvoices)
            invoices.GET("/:id", invoiceController.GetInvoice)
            invoices.POST("", invoiceController.CreateInvoice)
        }

        // Customer routes
        customers := api.Group("/customers")
        {
            customers.GET("", invoiceController.GetCustomers)
            customers.POST("", invoiceController.CreateCustomer)
        }
    }

    log.Println("Server starting on :8080")
    r.Run(":8080")
}