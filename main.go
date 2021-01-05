package main

import (
	"io"
	"log"
	"messanger/controllers"
	"messanger/middleware"
	"messanger/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()

	db, err := models.SetupModels()
	if err != nil {
		log.Fatal(err)

	}

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Use(middleware.AuthorizeJWT())

	users := r.Group("/users")
	users.POST("", controllers.CreateUser)
	users.PUT("", controllers.UpdateUser)
	users.GET("/:id", controllers.GetUserByID)
	users.DELETE("/:id", controllers.DeleteUserByID)
	users.POST("/login", controllers.GetToken)

	messages := r.Group("/messages")
	messages.POST("", controllers.CreateMessage)
	messages.PUT("", controllers.UpdateMessage)
	messages.DELETE("/:id", controllers.DeleteMessageByID)
	messages.GET("/:id", controllers.GetMessages)

	rooms := r.Group("/rooms")

	rooms.POST("", controllers.CreateRoom)

	r.Run(":" + os.Getenv("PORT"))
}
