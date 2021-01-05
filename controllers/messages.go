package controllers

import (
	"log"
	"messanger/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateMessage(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Permission denided"})
		return
	}
	db := c.MustGet("db").(*pgxpool.Pool)
	message := &models.Message{}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	message.CreatedBy = id.(int64)
	err := db.QueryRow(c, `
	INSERT INTO messages(room_id, created_by, text ) VALUES($1, $2, $3)  RETURNING id,created;
	`, message.RoomID, message.CreatedBy, message.Text).Scan(&message.ID, &message.Created)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": message})
}

func UpdateMessage(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Permission denided"})
		return
	}

	db := c.MustGet("db").(*pgxpool.Pool)
	message := &models.Message{}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err := db.QueryRow(c, `
	UPDATE messages SET text =$1 WHERE id = $2 AND created_by =$3 RETURNING room_id, created_by, created;
	`, message.Text, message.ID, id.(int64)).Scan(&message.RoomID, &message.CreatedBy, &message.Created)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Update error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func DeleteMessageByID(c *gin.Context) {
	userId, ok := c.Get("user_id")
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
	DELETE FROM messages WHERE id= $1 AND created_by = $2
	`, id.ID, userId)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}

func GetMessages(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	id := &models.ID{}

	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	items := make([]*models.Message, 0)
	rows, err := db.Query(c, `
	SELECT id,room_id, created_by, text, created FROM messages WHERE room_id= $1 ORDER BY created DESC LIMIT 500;
	`, id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Messages not found!"})
	}

	for rows.Next() {
		item := &models.Message{}
		rows.Scan(&item.ID, &item.RoomID, &item.CreatedBy, &item.Text, &item.Created)
		items = append(items, item)
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
