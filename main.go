package main

import (
	"github.com/Tell-Me-Your-Love-Way/castor/domains/amazon"
	"github.com/Tell-Me-Your-Love-Way/castor/domains/magalu"
	"github.com/Tell-Me-Your-Love-Way/castor/domains/scrapping"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil { // Load .env file}
		panic("Error loading .env file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "a-very-complex-password", // sem senha por padrão
		DB:       0,                         // banco padrão
	})
	amazon.ServiceInstance = amazon.NewService(rdb)
	magalu.ServiceInstance = magalu.NewService(rdb)
	scrapping.ServiceInstance = scrapping.NewService()
	router := gin.Default()
	router.POST("/amazon", amazon.HandlerQuery)
	router.POST("/magalu", magalu.Handler)
	router.POST("/parse", scrapping.Handler)
	err = router.Run("127.0.0.1:8080")
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
