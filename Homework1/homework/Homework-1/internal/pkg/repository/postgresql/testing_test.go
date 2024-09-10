package postgresql

import (
	"github.com/golang/mock/gomock"
	mock_database "homework/Homework-1/internal/pkg/repository/postgresql/mocks"
	"homework/Homework-1/internal/pkg/service"
	"testing"
)

type pvzRepoFixtures struct {
	ctrl   *gomock.Controller
	repo   service.PVZRepo
	mockDB *mock_database.MockDBops
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockDB := mock_database.NewMockDBops(ctrl)
	repo := NewPVZ(mockDB)
	return pvzRepoFixtures{
		ctrl:   ctrl,
		repo:   repo,
		mockDB: mockDB,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}
