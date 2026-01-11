package main

import (
	"log"

	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/cache"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/controllers"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/db"
	"github.com/GoogleCloudPlatform/memorystore-samples/caching/go/sample-demo-app/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize Database and Cache
	db.InitDB()
	cache.InitRedis()
	defer db.Pool.Close()
	defer cache.Rdb.Close()

	r := gin.Default()
	r.Use(middleware.PrometheusMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Serve Static Files
	r.Static("/static", "./static")

	// Load Templates
	r.LoadHTMLGlob("templates/*")

	// Controllers
	homeCtrl := controllers.NewHomeController()
	itemCtrl := controllers.NewItemController()

	// Home Route
	r.GET("/", homeCtrl.Home)

	// API Routes
	api := r.Group("/api/item")
	{
		api.GET("/:id", itemCtrl.Get)
		api.GET("/random", itemCtrl.GetRandom)
		api.POST("/create", itemCtrl.Create)
		api.DELETE("/delete/:id", itemCtrl.Delete)
	}

	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
