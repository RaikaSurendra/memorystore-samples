package models

import "encoding/json"

type Item struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	FromCache   bool    `json:"fromCache"`
}

func (i *Item) ToJSON() (string, error) {
	bytes, err := json.Marshal(i)
	return string(bytes), err
}

func ItemFromJSON(data string) (*Item, error) {
	var item Item
	err := json.Unmarshal([]byte(data), &item)
	return &item, err
}
