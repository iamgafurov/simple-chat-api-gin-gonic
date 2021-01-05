package middleware

import (
	"fmt"
	"log"
	"messanger/service"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) > 0 {
			tokenString := authHeader
			//log.Print(tokenString)
			id, err := service.DecodeToken(tokenString)
			if err != nil || id == 0 {
				log.Print(err)
				fmt.Println("Not Authenication1")
				return
			}
			log.Print(id)
			c.Set("user_id", id)
			c.Next()

		} else {
			fmt.Println("Not Authenication")
			return
		}

	}
}
