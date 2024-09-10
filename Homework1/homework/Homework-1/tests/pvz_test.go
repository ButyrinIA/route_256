//go:build integration
// +build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/Homework-1/internal/pkg/repository"
	"homework/Homework-1/internal/pkg/repository/postgresql"
	"homework/Homework-1/tests/fixtures"
	"testing"
)

func TestGetPVZ(t *testing.T) {
	var (
		ctx = context.Background()
	)

	db.SetUp(t, "pvz")
	defer db.TearDown()
	// arrange
	repo := postgresql.NewPVZ(db.DB)
	respAdd, err := repo.Add(ctx, fixtures.PVZ().Valid().P())
	require.NoError(t, err)

	//act
	resp, err := repo.GetById(ctx, respAdd)
	//assert
	require.NoError(t, err)
	assert.Equal(t, resp, &repository.PVZ{
		ID:      respAdd,
		Name:    "pvz1",
		Address: "street1",
		Contact: "123",
	})
}
