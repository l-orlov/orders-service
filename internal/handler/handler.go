package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/l-orlov/orders-service/internal/cache"
	db2 "github.com/l-orlov/orders-service/internal/db"
	"github.com/l-orlov/orders-service/internal/model"
	"github.com/pkg/errors"
)

type Handler struct {
	database  *db2.Database
	cacheImpl *cache.Cache
}

func New(database *db2.Database, cacheImpl *cache.Cache) http.Handler {
	h := &Handler{
		database:  database,
		cacheImpl: cacheImpl,
	}

	router := gin.Default()
	router.GET("/orders", h.getOrders)
	router.GET("/orders/:id", h.getOrderByID)
	router.POST("/orders", h.postOrders)

	return router
}

// getOrders возвращает имеющиеся заказы
func (h *Handler) getOrders(c *gin.Context) {
	orders, err := h.database.GetOrders(c)
	if err != nil {
		log.Printf("h.database.GetOrders: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []*model.Order{}
	}

	c.IndentedJSON(http.StatusOK, orders)
}

// postOrders создает заказ
func (h *Handler) postOrders(c *gin.Context) {
	newOrder := &model.Order{}
	err := c.BindJSON(newOrder)
	if err != nil {
		log.Printf("c.BindJSON: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.cacheImpl.CreateOrder(c, newOrder)
	if err != nil {
		log.Printf("h.database.CreateOrder: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusCreated, newOrder)
}

// getOrders возвращает заказ по id
func (h *Handler) getOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.cacheImpl.GetOrder(c, id)
	if err != nil {
		if errors.Is(err, db2.ErrNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
			return
		}

		log.Printf("h.database.GetOrder: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.IndentedJSON(http.StatusOK, order)
}
