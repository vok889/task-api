package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-api/internal/auth"
	"task-api/internal/item"
	"task-api/internal/user"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect database
	db, err := gorm.Open(
		postgres.Open(
			os.Getenv("URL_DB"),
		),
	)
	if err != nil {
		log.Panic(err)
	}

	// Controller
	controller := item.NewController(db)
	userController := user.NewController(db, "secret")

	// verifyToken middleware
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	verifyToken := auth.Guard(secret)

	// Router setup
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}

	config.AllowCredentials = true
	r.Use(cors.New(config))

	// frontend URL
	config.AllowOrigins = []string{
		"http://localhost:8000",
		"http://127.0.0.1:8000",
	}
	r.Use(cors.New(config))
	// r.Use(mylog.Logger())
	// r.Use(mylog.Logger2())

	r.GET("/version", func(c *gin.Context) {
		version, err := GetLatestDBVersion(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": version})
	})

	r.POST("/items", verifyToken, controller.CreateItem)
	r.GET("/items", verifyToken, controller.FindItems)
	r.GET("/items/:id", verifyToken, controller.FindItemByID)
	r.PUT("/items/:id", verifyToken, controller.UpdateItem)
	r.PATCH("/items/:id", verifyToken, controller.UpdateItemStatus)
	r.DELETE("/items/:id", verifyToken, controller.DeleteItem)
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)

	r.GET("/test", func(ctx *gin.Context) {
		fmt.Println("------test------")
		value, exists := ctx.Get("example2")
		if exists {
			log.Printf("example: %v, %T\n", value, value)
		} else {
			log.Printf("example does not exists: %v, %T\n", value, value)
		}
		// for i := 0; i < 20; i++ {
		// 	fmt.Println(i)
		// 	time.Sleep(1 * time.Second)
		// }
		ctx.JSON(210, "test response")
	})

	// // Register router
	// items := r.Group("/items")
	// //item.Use(mylog.Logger2())
	// // item.Use(auth.BasicAuth([]auth.Credential{
	// // 	{"admin", "secret"},
	// // 	{"admin2", "1234"},
	// // }))
	// items.Use(auth.Guard(os.Getenv("JWT_SECRET")))
	// {
	// 	items.POST("", controller.CreateItem)
	// 	items.GET("", controller.FindItems)
	// 	items.PATCH("/:id", controller.UpdateItemStatus)
	// }

	// // Start server
	// fmt.Printf("RUN: %v,%T\n", port, port)
	// // if err := r.Run(":" + os.Getenv("PORT")); err != nil {
	// // 	log.Panic(err)
	// // }

	// // for non-windows server
	// //endless.DefaultHammerTime = 10 * time.Second
	// // if err := endless.ListenAndServe(":8081", r); err != nil {
	// // 	log.Panic(err)
	// // }

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":" + getPort(),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}

func Logger() {
	panic("unimplemented")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	return port
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
