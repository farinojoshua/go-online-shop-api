package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"online_shop/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: get products from database
		products, err := model.SelectProducts(db)
		if err != nil {
			log.Printf("failed to get products: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: give response
		c.JSON(http.StatusOK, products)
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: take id from request
		id := c.Param("id")

		// TODO: take product from database with id
		product, err := model.SelectProductByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}

			log.Printf("failed to get product: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// TODO: give response
		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		err := c.Bind(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		product.ID = uuid.New().String()
		err = model.InsertProduct(db, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var productReq model.Product
		err := c.Bind(&productReq)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "data produk tidak valid"})
			return
		}

		product, err := model.SelectProductByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil produk")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		if productReq.Name != "" {
			product.Name = productReq.Name
		}

		if productReq.Price != 0 {
			product.Price = productReq.Price
		}

		err = model.UpdateProduct(db, product)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengupdate produk: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		c.JSON(201, product)
	}
}

func DeleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
