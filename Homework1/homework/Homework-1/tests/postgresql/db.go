package postgresql

import (
	"context"
	"fmt"
	"homework/Homework-1/internal/pkg/db"
	"strings"
	"testing"
)

type TDB struct {
	DB db.DBops
}

func NewFromEnv() *TDB {
	db, err := db.NewDb(context.Background())
	if err != nil {
		panic(err)
	}
	return &TDB{DB: db}
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)

}

func (d *TDB) TearDown() {

}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {

	q := fmt.Sprintf("TRUNCATE table %s", strings.Join(tableName, ""))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}