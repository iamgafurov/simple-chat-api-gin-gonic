package models

import "time"

type Message struct {
	ID        int64     `json:"id`
	RoomID    int64     `json:"room_id"`
	CreatedBy int64     `json:"created_by"`
	Text      string    `json:"text"`
	Created   time.Time `json:"time"`
}
