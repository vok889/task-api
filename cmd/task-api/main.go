package main

import (
	"log"
	"task-api/internal/item"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// POST 	/items
// GET 		/items?status=xxxxx
// GET 		/items/:id
// PUT		/items/:id
// PATCH	/items/:id
// DELETE 	/items/:id

func main() {
	// Connect database
	db, err := gorm.Open(
		postgres.Open(
			"postgres://postgres:password@localhost:5432/task",
		),
	)
	if err != nil {
		log.Panic(err)
	}

	// Controller
	controller := item.NewController(db)

	// Router
	r := gin.Default()

	// Register router
	r.POST("/items", controller.CreateItem)
	r.GET("/items", controller.FindItems)
	r.PATCH("/items/:id", controller.UpdateItemStatus)

	// Start server
	if err := r.Run(); err != nil {
		log.Panic(err)
	}
}
