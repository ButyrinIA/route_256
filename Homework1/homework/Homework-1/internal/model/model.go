package model

import "time"

type Order struct {
	ID          int
	RecipientID int
	StorageDate time.Time
}
