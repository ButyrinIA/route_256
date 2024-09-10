package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/Homework-1/internal/pkg/repository"
	"testing"
)

func Test_GetById(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz WHERE id=$1", gomock.Any()).Return(nil)

		// act
		user, err := s.repo.GetById(ctx, id)

		// assert
		require.NoError(t, err)
		assert.Equal(t, int64(0), user.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz WHERE id=$1", gomock.Any()).
				Return(sql.ErrNoRows)

			// act
			user, err := s.repo.GetById(ctx, id)
			// assert
			require.EqualError(t, err, "not found")
			require.True(t, errors.Is(err, repository.ErrObjectNotFound))
			assert.Nil(t, user)
		})
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			s.mockDB.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,address,contact FROM pvz WHERE id=$1", gomock.Any()).
				Return(assert.AnError)

			// act
			user, err := s.repo.GetById(ctx, id)
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")
			assert.Nil(t, user)
		})
	})

}
func Test_Add(t *testing.T) {
	t.Parallel()

}
