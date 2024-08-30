package main

import (
	"log"
	"net/http"
	"task-api/internal/item"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// POST 	/items
// GET 		/items?status=xxxxx
// PATCH	/items/:id

// GET 		/items/:id
// PUT		/items/:id
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

	config := cors.DefaultConfig()
	// frontend URL
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
	}
	r.Use(cors.New(config))

	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	// Register router
	r.POST("/items", controller.CreateItem)
	r.GET("/items", controller.FindItems)
	r.PATCH("/items/:id", controller.UpdateItemStatus)

	// Start server
	if err := r.Run(); err != nil {
		log.Panic(err)
	}
}

type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

// TableName overrides the table name used by User to `profiles`
func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}

// GetLatestDBVersion returns the latest applied version from the goose_db_version table.
func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion

	// Query to get the latest version applied
	err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error
	if err != nil {
		return 0, err
	}

	return version.VersionID, nil
}
