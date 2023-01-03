package cache

import (
	"context"
	"sync"

	"github.com/l-orlov/orders-service/internal/db"
	"github.com/l-orlov/orders-service/internal/model"
	"github.com/pkg/errors"
)

type Cache struct {
	database *db.Database
	orders   map[string]*model.Order
	mutex    sync.RWMutex
}

func New(ctx context.Context, database *db.Database) (*Cache, error) {
	orders, err := database.GetOrders(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "database.GetOrders")
	}

	ordersMap := make(map[string]*model.Order, len(orders))
	for _, order := range orders {
		ordersMap[order.ID] = order
	}

	return &Cache{
		database: database,
		orders:   ordersMap,
		mutex:    sync.RWMutex{},
	}, nil
}

func (cache *Cache) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	order, ok := cache.orders[id]
	if ok {
		return order, nil
	}

	order, err := cache.database.GetOrder(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "database.GetOrder")
	}

	return order, nil
}

func (cache *Cache) CreateOrder(ctx context.Context, order *model.Order) error {
	cache.mutex.RLock()

	_, ok := cache.orders[order.ID]
	if ok {
		cache.mutex.RUnlock()
		return nil
	}

	cache.mutex.RUnlock()
	err := cache.database.CreateOrder(ctx, order)
	if err != nil {
		return errors.Wrap(err, "cache.database.CreateOrder")
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.orders[order.ID] = order

	return nil
}
