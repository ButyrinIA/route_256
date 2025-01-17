package repository

import "errors"

var ErrObjectNotFound = errors.New("not found")

type PVZ struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Contact string `db:"contact"`
}
