package services

import (
	"database/sql"
	"errors"
)

type orderRepo interface {
	// GetOrderById(id string) (*models.Order, error)
}

type orderService struct {
	ordersCache *map[string][]byte
	orderRepo   orderRepo
}

func (o *orderService) AddOrder() error {

	return nil
}

func (o *orderService) GetOrderById(id string) ([]byte, error) {
	if data, ok := (*o.ordersCache)[id]; ok {
		return data, nil
	}
	return nil, errors.New("Not exist")

}

func NewOrderService(db *sql.DB, cache *map[string][]byte, orderRepo orderRepo) *orderService {
	return &orderService{ordersCache: cache, orderRepo: orderRepo}
}
