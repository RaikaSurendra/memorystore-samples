package repositories

import (
	"context"

	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/db"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/models"
)

type ItemsRepository struct{}

func NewItemsRepository() *ItemsRepository {
	return &ItemsRepository{}
}

func (r *ItemsRepository) Get(id int64) (*models.Item, error) {
	var item models.Item
	err := db.Pool.QueryRow(context.Background(), "SELECT id, name, description, price FROM items WHERE id=$1", id).Scan(&item.ID, &item.Name, &item.Description, &item.Price)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemsRepository) GetMultiple(amount int) ([]models.Item, error) {
	// For demo, just get latest 'amount' items
	rows, err := db.Pool.Query(context.Background(), "SELECT id, name, description, price FROM items ORDER BY id DESC LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemsRepository) Create(item *models.Item) (int64, error) {
	var id int64
	err := db.Pool.QueryRow(context.Background(), "INSERT INTO items (name, description, price) VALUES ($1, $2, $3) RETURNING id", item.Name, item.Description, item.Price).Scan(&id)
	return id, err
}

func (r *ItemsRepository) Delete(id int64) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM items WHERE id=$1", id)
	return err
}

func (r *ItemsRepository) Exists(id int64) (bool, error) {
	var exists bool
	err := db.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM items WHERE id=$1)", id).Scan(&exists)
	return exists, err
}
