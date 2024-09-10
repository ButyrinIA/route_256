package storage

import "time"

type OrderDTO struct {
	ID          int
	RecipientID int
	StorageDate time.Time
	IssueDate   time.Time
	WasIssued   bool
	WasReturn   bool
	OurPVZ      bool
	Packaging   string
	Cost        int
}

type OrderOut struct {
	ID          int
	RecipientID int
	StorageDate time.Time
	WasIssued   bool
}

type PVZ struct {
	Name    string
	Address string
	Contact string
}
