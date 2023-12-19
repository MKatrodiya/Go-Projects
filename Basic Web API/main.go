package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

var products = []product{
	{ID: "1", Name: "Laptop", Price: 999.99, Quantity: 10},
	{ID: "2", Name: "Mouse", Price: 14.49, Quantity: 50},
	{ID: "3", Name: "Keyboard", Price: 100, Quantity: 25},
}

func getProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	p, err := getProductByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Product not found!",
		})
	}
	c.IndentedJSON(http.StatusOK, p)
}

func getProductByID(id string) (*product, error) {
	for i, p := range products {
		if p.ID == id {
			return &products[i], nil
		}
	}
	return nil, errors.New("Product not found!")
}

func buyProduct(c *gin.Context) {
	id := c.Query("id")
	quantity := c.Query("quantity")
	if id == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Product id is missing from query parameter!",
		})
		return
	}
	p, err := getProductByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Product not found!",
		})
		return
	}
	if p.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Product not available!",
		})
		return
	}
	if quantity == "" {
		p.Quantity -= 1
	} else {
		q, err := strconv.Atoi(quantity)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid quantity",
			})
			return
		}
		if q > p.Quantity {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"message": "Product not available in provided quantity",
			})
			return
		}
		p.Quantity -= q
	}
	c.IndentedJSON(http.StatusOK, p)
}

func addProduct(c *gin.Context) {
	var p product
	err := c.BindJSON(&p)
	if err != nil {
		return
	}

	products = append(products, p)
	c.IndentedJSON(http.StatusCreated, products)
}

func main() {
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/products/:id", getProduct)
	router.POST("/products", addProduct)
	router.PATCH("/products/buy", buyProduct)
	router.Run("localhost:8080")
}

// {"id": "4", "name":"Monitor", "price": "599.99", "quantity": "15"}
