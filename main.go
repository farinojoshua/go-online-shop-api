package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"online_shop/handler"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/database?sslmode=disable")
	if err != nil {
		fmt.Printf("failed to open database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("failed to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("connected to database")

	// _, err = migrate(db)
	// if err != nil {
	// 	fmt.Printf("failed to migrate database: %v\n", err)
	// 	os.Exit(1)
	// }

	r := gin.Default()

	r.GET("/api/v1/products", handler.ListProducts(db))
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout")

	r.POST("/api/v1/orders/:id/confirm")
	r.GET("/api/v1/orders/:id")

	r.POST("/admin/products", handler.CreateProduct(db))
	r.PUT("/admin/products/:id", handler.UpdateProduct(db))
	r.DELETE("/admin/products/:id", handler.DeleteProduct(db))

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}
