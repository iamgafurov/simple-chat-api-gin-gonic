package controllers

import (
	"log"
	"messanger/models"
	"net/http"

	"messanger/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	user := &models.User{}
	item := models.Registration{}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
	}
	log.Print(item)
	err = db.QueryRow(c, `
	INSERT INTO users(name,login,password) VALUES ($1,$2,$3) ON CONFLICT (login) DO NOTHING RETURNING id, name, login,created;
	`, item.Name, item.Login, hash).Scan(&user.ID, &user.Name, &user.Login, &user.Created)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	user := &models.User{}
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Permission denided"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	user.ID = id.(int64)
	err := db.QueryRow(c, `
	UPDATE users SET name =$1,login=$2 WHERE id =$3 RETURNING created
	`, user.Name, user.Login, user.ID).Scan(&user.Created)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUserByID(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	user := &models.User{}
	id := &models.ID{}

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := db.QueryRow(c, `
	SELECT id, name, login, created FROM users WHERE id = $1
	`, id.ID).Scan(&user.ID, &user.Name, &user.Login, &user.Created)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func DeleteUserByID(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Permission denided"})
		return
	}
	db := c.MustGet("db").(*pgxpool.Pool)
	id := &models.ID{}

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	_, err := db.Exec(c, `
	DELETE FROM users WHERE id= $1 
	`, id.ID)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func GetToken(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	auth := &models.Auth{}
	user := &models.Registration{}
	var id int64
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := db.QueryRow(c, `
	SELECT id, password  FROM users WHERE login = $1
	`, auth.Login).Scan(&id, &user.Password)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login or password!"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login or password!"})
		return
	}

	token, err := service.GenerateToken(id, user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal"})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
