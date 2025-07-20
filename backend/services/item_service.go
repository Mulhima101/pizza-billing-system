package services

import (
    "database/sql"
    "pizza-billing-backend/models"
)

type ItemService struct {
    db *sql.DB
}

func NewItemService(db *sql.DB) *ItemService {
    return &ItemService{db: db}
}

func (s *ItemService) GetAllItems() ([]models.Item, error) {
    query := `
        SELECT i.id, i.name, i.description, i.price, i.category_id, i.is_active, i.created_at, c.name
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.id
        WHERE i.is_active = true
        ORDER BY i.name
    `
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.Item
    for rows.Next() {
        var item models.Item
        err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, 
                        &item.CategoryID, &item.IsActive, &item.CreatedAt, &item.CategoryName)
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    return items, nil
}

func (s *ItemService) GetItemByID(id int) (*models.Item, error) {
    query := `
        SELECT i.id, i.name, i.description, i.price, i.category_id, i.is_active, i.created_at, c.name
        FROM items i
        LEFT JOIN categories c ON i.category_id = c.id
        WHERE i.id = $1
    `
    
    var item models.Item
    err := s.db.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, 
                                         &item.Price, &item.CategoryID, &item.IsActive, 
                                         &item.CreatedAt, &item.CategoryName)
    if err != nil {
        return nil, err
    }

    return &item, nil
}

func (s *ItemService) CreateItem(item *models.Item) (*models.Item, error) {
    query := `
        INSERT INTO items (name, description, price, category_id, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `
    
    err := s.db.QueryRow(query, item.Name, item.Description, item.Price, 
                        item.CategoryID, item.IsActive).Scan(&item.ID, &item.CreatedAt)
    if err != nil {
        return nil, err
    }

    return item, nil
}

func (s *ItemService) UpdateItem(item *models.Item) (*models.Item, error) {
    query := `
        UPDATE items 
        SET name = $1, description = $2, price = $3, category_id = $4, is_active = $5
        WHERE id = $6
        RETURNING created_at
    `
    
    err := s.db.QueryRow(query, item.Name, item.Description, item.Price, 
                        item.CategoryID, item.IsActive, item.ID).Scan(&item.CreatedAt)
    if err != nil {
        return nil, err
    }

    return item, nil
}

func (s *ItemService) DeleteItem(id int) error {
    query := `UPDATE items SET is_active = false WHERE id = $1`
    _, err := s.db.Exec(query, id)
    return err
}

func (s *ItemService) GetAllCategories() ([]models.Category, error) {
    query := `SELECT id, name, description, created_at FROM categories ORDER BY name`
    
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var category models.Category
        err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }

    return categories, nil
}