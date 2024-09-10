package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"homework/Homework-1/internal/infrastructure/kafka"
	"homework/Homework-1/internal/pkg/middleware"
	"homework/Homework-1/internal/pkg/repository"
	"homework/Homework-1/internal/pkg/repository/in_memory_cache"
	"homework/Homework-1/internal/pkg/repository/redis"
	"io"
	"net/http"
	"strconv"
	"time"
)

const Port = ":9000"
const queryParamKey = "key"

//go:generate mockgen -source ./httpservice.go -destination ./mocks/mock_repository.go -package=mock_repository

type PVZRepo interface {
	Add(ctx context.Context, pvz *repository.PVZ) (int64, error)
	GetById(ctx context.Context, id int64) (*repository.PVZ, error)
	ListOfPVZ(ctx context.Context) ([]*repository.PVZ, error)
	Update(ctx context.Context, pvz *repository.PVZ) error
	Delete(ctx context.Context, id int64) error
}

type Server struct {
	repo     PVZRepo
	producer *kafka.Producer
	cache    *in_memory_cache.InMemoryCache
	redis    *redis.Redis
}

func NewServer(repo PVZRepo, kafkaProducer *kafka.Producer) Server {
	return Server{
		repo:     repo,
		producer: kafkaProducer,
		cache:    in_memory_cache.NewInMemoryCache(1 * time.Minute),
		redis:    redis.NewRedis(),
	}
}

type addPVZRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type addPVZResponse struct {
	ID      int64  `json:"ID"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func CreateRouter(implemetation Server) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/pvz", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implemetation.Create(w, req)
		case http.MethodGet:
			implemetation.List(w, req)
		default:
			fmt.Println("error")
		}
	})

	router.HandleFunc(fmt.Sprintf("/pvz/{%s:[0-9]+}", queryParamKey), func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			implemetation.GetByID(w, req)
		case http.MethodPut:
			implemetation.Update(w, req)
		case http.MethodDelete:
			implemetation.Delete(w, req)
		default:
			fmt.Println("error")
		}
	})
	return router
}

//метод create pvz

func (s *Server) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.producer.SendMessage(req, body)
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
	}

	var pvz addPVZRequest
	if err := json.Unmarshal(body, &pvz); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pvzRepo := &repository.PVZ{
		Name:    pvz.Name,
		Address: pvz.Address,
		Contact: pvz.Contact,
	}
	key, pvzJson, status := s.add(req.Context(), pvzRepo)
	if status == http.StatusOK {
		s.cache.Set(key, pvzJson)
	}
	w.WriteHeader(status)
	w.Write(pvzJson)
}

//метод get pvz by id

func (s *Server) GetByID(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.producer.SendMessage(req, body)
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
	}
	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateGetByID(keyInt) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cachedData, err := s.redis.Get(req.Context(), key)
	if err == nil {
		w.Write(cachedData)
		return
	}
	/*cachedData, ok := s.cache.Get(keyInt)
	if ok {
		w.Write(cachedData)
		return
	}*/
	data, status := s.get(req.Context(), keyInt)
	if status == http.StatusOK {
		s.redis.Set(req.Context(), key, data)
		//s.cache.Set(keyInt, data)
	}
	w.WriteHeader(status)
	w.Write(data)
}

//метод delete pvz

func (s *Server) Delete(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.producer.SendMessage(req, body)
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
	}

	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateGetByID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err, status := s.delete(req.Context(), id)
	if err != nil {
		fmt.Println(err)
	}
	s.cache.Delete(id)
	w.WriteHeader(status)
}

//метод get pvz list

func (s *Server) List(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.producer.SendMessage(req, body)
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
	}

	pvzs, err := s.repo.ListOfPVZ(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, pvz := range pvzs {
		pvzJson, _ := json.Marshal(pvz)
		s.cache.Set(pvz.ID, pvzJson)

	}
	err = json.NewEncoder(w).Encode(pvzs)
	if err != nil {
		return
	}

}

// метод update pvz

func (s *Server) Update(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.producer.SendMessage(req, body)
	if err != nil {
		fmt.Printf("Error producing message: %v\n", err)
	}

	key, ok := mux.Vars(req)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !validateGetByID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var pvz repository.PVZ
	if err = json.Unmarshal(body, &pvz); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pvz.ID = id
	s.cache.Delete(id)
	pvzJson, status := s.update(req.Context(), &pvz)
	if status == http.StatusOK {
		s.cache.Set(id, pvzJson)
	}
	w.WriteHeader(status)
	w.Write(pvzJson)
}

func (s *Server) delete(ctx context.Context, id int64) (error, int) {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return err, http.StatusNotFound
		}
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
func (s *Server) update(ctx context.Context, pvz *repository.PVZ) ([]byte, int) {
	err := s.repo.Update(ctx, pvz)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	pvzJson, _ := json.Marshal(&addPVZResponse{
		ID:      pvz.ID,
		Name:    pvz.Name,
		Address: pvz.Address,
		Contact: pvz.Contact,
	})
	return pvzJson, http.StatusOK
}

func (s *Server) add(ctx context.Context, pvzRepo *repository.PVZ) (int64, []byte, int) {

	id, err := s.repo.Add(ctx, pvzRepo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil, http.StatusBadRequest
		}
		return 0, nil, http.StatusInternalServerError
	}

	pvzJson, _ := json.Marshal(&addPVZResponse{
		ID:      id,
		Name:    pvzRepo.Name,
		Address: pvzRepo.Address,
		Contact: pvzRepo.Contact,
	})
	return id, pvzJson, http.StatusOK

}

func (s *Server) get(ctx context.Context, key int64) ([]byte, int) {
	pvz, err := s.repo.GetById(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	pvzJson, err := json.Marshal(pvz)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return pvzJson, http.StatusOK
}

func validateGetByID(key int64) bool {
	if key <= 0 {
		return false
	}
	return true
}
