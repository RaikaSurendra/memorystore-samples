package controllers

import (
	"log"
	"strconv"

	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/cache"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/models"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/repositories"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "The total number of cache hits",
	})
	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "The total number of cache misses",
	})
)

func init() {
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
	// Initialize to 0 so they are exported immediately
	cacheHits.Add(0)
	cacheMisses.Add(0)
}

// DataController mirrors the Java DataController which handles business logic (Cache + DB)
type DataController struct {
	Repo *repositories.ItemsRepository
}

func NewDataController() *DataController {
	return &DataController{
		Repo: repositories.NewItemsRepository(),
	}
}

func (c *DataController) Get(id int64) (*models.Item, error) {
	idStr := strconv.FormatInt(id, 10)

	// 1. Check Cache
	val, err := cache.Get(idStr)
	if err == nil && val != "" {
		item, err := models.ItemFromJSON(val)
		if err == nil {
			cacheHits.Inc()
			item.FromCache = true
			return item, nil
		}
		log.Printf("Error parsing cache JSON: %v", err)
	}

	// Cache Miss
	cacheMisses.Inc()

	// 2. Check Database
	item, err := c.Repo.Get(id)
	if err != nil {
		return nil, err
	}

	// 3. Update Cache
	itemJson, err := item.ToJSON()
	if err == nil {
		_ = cache.Set(idStr, itemJson)
	} else {
		log.Printf("Error serializing item to JSON: %v", err)
	}

	return item, nil
}

func (c *DataController) GetMultiple(amount int) ([]models.Item, error) {
	return c.Repo.GetMultiple(amount)
}

func (c *DataController) Create(item *models.Item) (int64, error) {
	id, err := c.Repo.Create(item)
	if err != nil {
		return 0, err
	}

	// Write-Through / Cache-Update
	item.ID = id
	itemJson, err := item.ToJSON()
	if err == nil {
		_ = cache.Set(strconv.FormatInt(id, 10), itemJson)
	}

	return id, nil
}

func (c *DataController) Delete(id int64) error {
	err := c.Repo.Delete(id)
	if err != nil {
		log.Printf("Error deleting from DB: %v", err)
	}

	// Always delete from cache
	_ = cache.Delete(strconv.FormatInt(id, 10))
	return nil
}
