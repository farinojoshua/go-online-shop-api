package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"online_shop/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO : mengambil data pesanan dari request
		var checkoutOrder model.Checkout
		err := c.BindJSON(&checkoutOrder)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "data produk tidak valid"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)

		for _, o := range checkoutOrder.Products {
			ids = append(ids, o.ID)
			orderQty[o.ID] = o.Quantity
		}

		// TODO : mengambil data produk dari database
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil produk")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// TODO : membuat kata sandi
		passcode := generatePasscode(5)

		// TODO : hash kata sandi
		hashcode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)
		if err != nil {
			log.Printf("Terjadi kesalahan saat menghasilkan hash: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		hashcodeString := string(hashcode)

		// TODO : buat order dan detail
		order := model.Order{
			ID:         uuid.New().String(),
			Email:      checkoutOrder.Email,
			Address:    checkoutOrder.Address,
			PassCode:   &hashcodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		for _, p := range products {
			total := p.Price * int64(orderQty[p.ID])

			detail := model.OrderDetail{
				ID:        uuid.New().String(),
				OrderID:   order.ID,
				ProductID: p.ID,
				Quantity:  orderQty[p.ID],
				Price:     p.Price,
				Total:     total,
			}

			details = append(details, detail)

			order.GrandTotal += total
		}

		model.CreateOrder(db, order, details)

		orderWithDetail := model.OrderWithDetail{
			Order:   order,
			Details: details,
		}

		orderWithDetail.Order.PassCode = &passcode

		c.JSON(http.StatusCreated, orderWithDetail)
	}
}

func generatePasscode(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func ConfirmOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO : mengambil id order dari request
		id := c.Param("id")

		// TODO : membaca request body
		var confirmReq model.ConfirmOrder
		err := c.BindJSON(&confirmReq)
		if err != nil {
			log.Printf("Terjadi kesalahan saat membaca request body: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "data pesanan tidak valid"})
			return
		}

		// TODO : mengambil order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil pesanan: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		if order.PassCode == nil {
			log.Println("passcode tidak valid")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// TODO : memverifikasi kata sandi
		if err := bcrypt.CompareHashAndPassword([]byte(*order.PassCode), []byte(confirmReq.Passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat memverifikasi kata sandi: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// TODO : pastikan pesanan belum dibayarkan
		if order.PaidAt != nil {
			log.Println("pesanan sudah dibayarkan")
			c.JSON(http.StatusBadRequest, gin.H{"error": "pesanan sudah dibayarkan"})
			return
		}

		// TODO : Cocokkan jumlah pembayaran
		if order.GrandTotal != confirmReq.Amount {
			log.Printf("Terjadi kesalahan saat memverifikasi pembayaran: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// TODO : update status order
		current := time.Now()
		err = model.UpdateOrderByID(db, id, confirmReq, current)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengupdate pesanan: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		order.PassCode = nil
		order.PaidAt = &current
		order.PaidBank = &confirmReq.Bank
		order.PaidAccountNumber = &confirmReq.AccountNumber

		// TODO : kirim response
		c.JSON(http.StatusOK, order)
	}
}

func GetOrder(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO : mengambil id order dari request
		id := c.Param("id")

		// TODO : ambil passcode dari query parameter
		passcode := c.Query("passcode")

		// TODO : membaca order dari database
		order, err := model.SelectOrderByID(db, id)
		if err != nil {
			log.Printf("Terjadi kesalahan saat mengambil pesanan: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		// TODO : pastikan order valid
		if order.PassCode == nil {
			log.Printf("Terjadi kesalahan saat mengambil pesanan: passcode tidak valid\n")
			c.JSON(http.StatusBadRequest, gin.H{"error": "passcode tidak valid"})
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(*order.PassCode), []byte(passcode)); err != nil {
			log.Printf("Terjadi kesalahan saat mengambil pesanan: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "terjadi kesalahan pada server"})
			return
		}

		order.PassCode = nil

		// TODO : kirim response
		c.JSON(http.StatusOK, order)
	}
}
