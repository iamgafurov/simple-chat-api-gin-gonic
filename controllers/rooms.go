package controllers

import (
	"log"
	"messanger/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateRoom(c *gin.Context) {
	db := c.MustGet("db").(*pgxpool.Pool)
	room := &models.Room{}

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	log.Print(room)

	roomID := getRoomIDByMembersID(c, room.FirstMemberID, room.SecondMemberID)

	if roomID != 0 {
		c.JSON(http.StatusOK, gin.H{"id": roomID})
		return
	}

	roomID = getRoomIDByMembersID(c, room.SecondMemberID, room.FirstMemberID)

	if roomID != 0 {
		c.JSON(http.StatusOK, gin.H{"id": roomID})
		return
	}

	err := db.QueryRow(c, `
	INSERT INTO rooms(name,first_member_id,second_member_id) VALUES ($1,$2,$3) RETURNING id;
	`, room.Name, room.FirstMemberID, room.SecondMemberID).Scan(&room.ID)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insert error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": room.ID})
}

func getRoomIDByMembersID(c *gin.Context, first_id int64, second_id int64) int64 {
	db := c.MustGet("db").(*pgxpool.Pool)
	id := int64(0)

	err := db.QueryRow(c, `
	SELECT id FROM rooms WHERE first_member_id = $1 and second_member_id = $2
	`, first_id, second_id).Scan(&id)
	if err != nil {
		log.Print(err)
		return 0
	}

	return id
}
