package models

import "time"

type Room struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	FirstMemberID  int64     `json:"first_member_id"`
	SecondMemberID int64     `json:"second_member_id"`
	Created        time.Time `json:"created"`
}
