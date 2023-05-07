package dto

import (
	"time"

	"github.com/tmp-friends/anirec-batch/functions/models"
)

type Record struct {
	ID          int         `json:"id"`
	Comment     string      `json:"comment"`
	RatingState string      `json:"rating_state"`
	CreatedAt   time.Time   `json:"created_at"`
	User        User        `json:"user"`
	Work        models.Work `json:"work"`
}
