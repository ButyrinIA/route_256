package service

import (
	"github.com/golang/mock/gomock"
	mock_repository "homework/Homework-1/internal/pkg/service/mocks"
	"testing"
)

type pvzRepoFixtures struct {
	ctrl    *gomock.Controller
	srv     Server
	mockPvz *mock_repository.MockPVZRepo
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPvz := mock_repository.NewMockPVZRepo(ctrl)
	srv := Server{mockPvz, nil}
	return pvzRepoFixtures{
		ctrl:    ctrl,
		mockPvz: mockPvz,
		srv:     srv,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}
