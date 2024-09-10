package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/Homework-1/internal/pkg/repository"
)

//go:generate mockgen -source ./pvz.go -destination ./mocks/mock_database.go -package=mock_database

type DBops interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetPool(_ context.Context) *pgxpool.Pool
	Close(_ context.Context)
}
type PVZRepo struct {
	db DBops
}

func NewPVZ(database DBops) *PVZRepo {
	return &PVZRepo{db: database}
}

func (r *PVZRepo) Add(ctx context.Context, pvz *repository.PVZ) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO pvz(name,address,contact) VALUES ($1,$2,$3) RETURNING id;`, pvz.Name, pvz.Address, pvz.Contact).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err

}

func (r *PVZRepo) GetById(ctx context.Context, id int64) (*repository.PVZ, error) {
	var pvz repository.PVZ
	err := r.db.Get(ctx, &pvz, "SELECT id,name,address,contact FROM pvz WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err

	}
	return &pvz, nil
}
func (r *PVZRepo) ListOfPVZ(ctx context.Context) ([]*repository.PVZ, error) {
	var pvzs []*repository.PVZ
	err := r.db.Select(ctx, &pvzs, "SELECT id, name, address, contact FROM pvz")
	if err != nil {
		return nil, err
	}
	return pvzs, nil
}

func (r *PVZRepo) Update(ctx context.Context, pvz *repository.PVZ) error {
	_, err := r.db.Exec(ctx, "UPDATE pvz SET name = $1, address = $2, contact = $3 WHERE id = $4", pvz.Name, pvz.Address, pvz.Contact, pvz.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		return err
	}
	return nil
}

func (r *PVZRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM pvz WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		return err
	}
	return nil
}
