package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/Homework-1/internal/pkg/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_GetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().GetById(gomock.Any(), id).Return(&repository.PVZ{
			ID:      1,
			Name:    "PVZ1",
			Address: "street1",
			Contact: "123",
		}, nil)
		result, status := s.srv.get(ctx, id)

		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"Name\":\"PVZ1\",\"Address\":\"street1\",\"Contact\":\"123\"}", string(result))
	})
	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := errors.New("repository error")
		s.mockPvz.EXPECT().GetById(gomock.Any(), id).Return(nil, expectedErr)
		result, status := s.srv.get(ctx, id)

		require.Equal(t, http.StatusInternalServerError, status)
		assert.Nil(t, result)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().GetById(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
		result, status := s.srv.get(ctx, id)

		require.Equal(t, http.StatusNotFound, status)
		assert.Nil(t, result)
	})
}

func Test_Create(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		pvz = repository.PVZ{
			ID:      1,
			Name:    "PVZ1",
			Address: "street1",
			Contact: "123",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().Add(gomock.Any(), &pvz).Return(int64(1), nil)
		result, status := s.srv.add(ctx, &pvz)

		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"PVZ1\",\"address\":\"street1\",\"contact\":\"123\"}", string(result))
	})
	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := errors.New("repository error")
		s.mockPvz.EXPECT().Add(gomock.Any(), &pvz).Return(int64(0), expectedErr)
		result, status := s.srv.add(ctx, &pvz)

		require.Equal(t, http.StatusInternalServerError, status)
		assert.Nil(t, result)
	})
	t.Run("no rows", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().Add(gomock.Any(), &repository.PVZ{}).Return(int64(0), sql.ErrNoRows)
		result, status := s.srv.add(ctx, &repository.PVZ{})

		require.Equal(t, http.StatusBadRequest, status)
		assert.Nil(t, result)
	})

}
func Test_Delete(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err, status := s.srv.delete(ctx, id)

		require.Equal(t, http.StatusOK, status)
		assert.Nil(t, err)
	})
	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := errors.New("repository error")
		s.mockPvz.EXPECT().Delete(gomock.Any(), id).Return(expectedErr)
		err, status := s.srv.delete(ctx, id)

		require.Equal(t, http.StatusInternalServerError, status)
		assert.NotNil(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := repository.ErrObjectNotFound
		s.mockPvz.EXPECT().Delete(gomock.Any(), id).Return(expectedErr)
		err, status := s.srv.delete(ctx, id)

		require.Equal(t, http.StatusNotFound, status)
		assert.NotNil(t, err)
	})

}
func Test_List(t *testing.T) {
	t.Parallel()
	var (
		responsePVZs []*repository.PVZ
		pvz          = repository.PVZ{
			ID:      1,
			Name:    "PVZ1",
			Address: "street1",
			Contact: "123",
		}
	)

	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()

		s.mockPvz.EXPECT().ListOfPVZ(gomock.Any()).Return([]*repository.PVZ{&pvz}, nil)

		req := httptest.NewRequest("GET", "/pvz", nil)
		w := httptest.NewRecorder()
		s.srv.List(w, req)

		err := json.Unmarshal(w.Body.Bytes(), &responsePVZs)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		require.Len(t, responsePVZs, 1)
		assert.Equal(t, &pvz, responsePVZs[0])

	})
	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := errors.New("repository error")
		s.mockPvz.EXPECT().ListOfPVZ(gomock.Any()).Return(nil, expectedErr)
		req := httptest.NewRequest("GET", "/pvz", nil)
		w := httptest.NewRecorder()
		s.srv.List(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Empty(t, w.Body.String())
	})
	t.Run("emptry list", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()

		s.mockPvz.EXPECT().ListOfPVZ(gomock.Any()).Return([]*repository.PVZ{}, nil) // Возвращаем пустой список
		req := httptest.NewRequest("GET", "/pvz", nil)
		w := httptest.NewRecorder()
		s.srv.List(w, req)

		err := json.Unmarshal(w.Body.Bytes(), &responsePVZs)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, responsePVZs)
	})
}

func Test_Update(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		pvz = repository.PVZ{
			ID:      1,
			Name:    "PVZ1",
			Address: "STREET1",
			Contact: "12345",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		s.mockPvz.EXPECT().Update(gomock.Any(), &pvz).Return(nil)
		result, status := s.srv.update(ctx, &pvz)

		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"ID\":1,\"name\":\"PVZ1\",\"address\":\"STREET1\",\"contact\":\"12345\"}", string(result))
	})
	t.Run("repository error", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := errors.New("repository error")
		s.mockPvz.EXPECT().Update(gomock.Any(), &pvz).Return(expectedErr)
		result, status := s.srv.update(ctx, &pvz)

		require.Equal(t, http.StatusInternalServerError, status)
		assert.Nil(t, result)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		expectedErr := repository.ErrObjectNotFound
		s.mockPvz.EXPECT().Update(gomock.Any(), &pvz).Return(expectedErr)
		result, status := s.srv.update(ctx, &pvz)

		require.Equal(t, http.StatusNotFound, status)
		assert.Nil(t, result)
	})
	t.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		s := setUp(t)
		defer s.tearDown()
		req := httptest.NewRequest("POST", "/pvz/-1", strings.NewReader("{\"Name\":\"PVZ1\",\"Address\":\"STREET1\",\"Contact\":\"12345\"}"))
		w := httptest.NewRecorder()
		s.srv.Update(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		assert.Empty(t, w.Body.String())
	})

}

func Test_validateGetByID(t *testing.T) {
	t.Parallel()
	t.Run("ok", func(t *testing.T) {
		result := validateGetByID(1)
		assert.True(t, result)
	})
	t.Run("fail", func(t *testing.T) {
		result := validateGetByID(-1)
		assert.False(t, result)
	})
}
