package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"homework/Homework-1/internal/model"
	"io"
	"os"
	"time"
)

const storageName = "file/storage"

var ErrProductStorage = errors.New("срок хранения в прошлом")
var ErrProductIssue = errors.New("не все заказы принадлежат вам")
var ErrTestId = errors.New("заказ не найден")
var ErrPackagingType = errors.New("invalid type of packaging")
var ErrMaxWeight = errors.New("the weight of the order exceeds the maximum allowed")
var ErrDelete = errors.New("срок хранения не истек или заказ был выдан")
var ErrNoOreders = errors.New("нет заказов для выдачи")
var ErrReturnOrder = errors.New("выдача не из нашего пвз")

type Storage struct {
	Storage *os.File
}

// создание файла

func NewStorage() (Storage, error) {
	file, err := os.OpenFile(storageName, os.O_CREATE, 0777)
	if err != nil {
		return Storage{}, err
	}
	return Storage{Storage: file}, nil

}

// добавление в файл заказа

func (s *Storage) AcceptOrder(id int, rid int, t time.Time, packType string, weight int) error {
	if t.Before(time.Now()) {
		return ErrProductStorage
	}
	all, err := s.listAll()
	if err != nil {
		return err
	}
	pack, ok := packagingType[packType]
	if !ok {
		return ErrPackagingType
	}
	if pack.MaxWeight > 0 && weight >= pack.MaxWeight {
		return ErrMaxWeight
	}

	newOrder := OrderDTO{
		ID:          id,
		RecipientID: rid,
		StorageDate: t,
		IssueDate:   time.Time{},
		WasIssued:   false,
		Packaging:   pack.Name,
	}
	newOrder.Cost += pack.CostIncrease
	all = append(all, newOrder)
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil

}

// удаление из файла заказа

func (s *Storage) Delete(id int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}

	err = testId(id, all)
	if err != nil {
		return err
	}
	for index, order := range all {
		if order.ID == id && order.StorageDate.Before(time.Now()) && !order.WasIssued {
			all = remove(all, index)
			break
		}
		if order.ID == id && (!order.StorageDate.Before(time.Now()) || order.WasIssued) {
			return ErrDelete
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// выдать заказы клиенту

func (s *Storage) IssueOrder(idSlice []int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	for _, id := range idSlice {
		err := testId(id, all)
		if err != nil {
			return err
		}
	}
	IssueList := make([]model.Order, 0, len(all))
	quantity := 0
	for _, id := range idSlice {
		for index, order := range all {
			if id == order.ID && !order.WasIssued && order.StorageDate.After(time.Now()) {
				all[index].IssueDate = time.Now()
				all[index].WasIssued = true
				all[index].OurPVZ = true
				IssueList = append(IssueList, model.Order{
					ID:          order.ID,
					RecipientID: order.RecipientID,
				})
				quantity += 1
				break
			}
			if id == order.ID && !order.WasIssued && !order.StorageDate.After(time.Now()) {
				break
			}
		}
	}
	if quantity == 0 {
		return ErrNoOreders
	}
	for index, _ := range IssueList {
		if IssueList[0].RecipientID != IssueList[index].RecipientID {
			return ErrProductIssue
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// список заказов

func (s *Storage) ListOrder(rID int, N int, inPVZ bool) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	listOrders := make([]OrderOut, 0, N)
	if N != 0 && !inPVZ {
		for index := len(all) - 1; index >= 0; index-- {
			if all[index].RecipientID == rID && len(listOrders) < N {
				listOrders = append(listOrders, OrderOut{
					ID:          all[index].ID,
					RecipientID: all[index].RecipientID,
					StorageDate: all[index].StorageDate,
					WasIssued:   all[index].WasIssued,
				})
			}
		}
	}
	if N == 0 && inPVZ {
		for _, order := range all {
			if order.RecipientID == rID && !order.WasIssued && !order.WasReturn {
				listOrders = append(listOrders, OrderOut{
					ID:          order.ID,
					RecipientID: order.RecipientID,
					StorageDate: order.StorageDate,
					WasIssued:   order.WasIssued,
				})
			}
		}

	}
	for _, order := range listOrders {
		fmt.Println(order)
	}
	return nil
}

// возврат от клиента

func (s *Storage) ReturnOrder(rID int, ID int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	for index, order := range all {
		if time.Since(order.IssueDate).Hours() <= 48 && order.RecipientID == rID && order.ID == ID && order.WasIssued && order.OurPVZ {
			all[index].WasIssued = false
			all[index].WasReturn = true
		}
		if time.Since(order.IssueDate).Hours() <= 48 && order.RecipientID == rID && order.ID == ID && order.WasIssued && !order.OurPVZ {
			return ErrReturnOrder
		}
	}
	err = writeBytes(all)
	if err != nil {
		return err
	}
	return nil
}

// список возвратов пагинировано

func (s *Storage) ReturnList(pageSize int, pageNumber int) error {
	all, err := s.listAll()
	if err != nil {
		return err
	}
	returnlist := make([]OrderOut, 0, len(all))
	for _, orders := range all {
		if orders.WasReturn {
			returnlist = append(returnlist, OrderOut{
				ID:          orders.ID,
				RecipientID: orders.RecipientID,
				StorageDate: orders.StorageDate,
				WasIssued:   orders.WasIssued,
			})
		}
	}
	startI := (pageNumber - 1) * pageSize
	endI := startI + pageSize

	if startI >= len(returnlist) {
		startI = len(returnlist)

	}

	if endI > len(returnlist) {
		endI = len(returnlist)
	}
	for _, order := range returnlist[startI:endI] {
		fmt.Println(order)
	}
	return nil
}

// запись в файл
func writeBytes(order []OrderDTO) error {
	rawBytes, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = os.WriteFile(storageName, rawBytes, 0777)
	if err != nil {
		return err
	}
	return nil
}

// удаление из файла
func remove(slice []OrderDTO, s int) []OrderDTO {
	return append(slice[:s], slice[s+1:]...)
}

// чтение из файла в буфер
func (s *Storage) listAll() ([]OrderDTO, error) {
	reader := bufio.NewReader(s.Storage)
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var orders []OrderDTO
	if len(rawBytes) == 0 {
		return orders, nil
	}
	err = json.Unmarshal(rawBytes, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func testId(id int, all []OrderDTO) error {

	test := false
	// слайз id заказов
	testId := make([]int, 0, len(all))
	for _, order := range all {
		testId = append(testId, order.ID)
	}
	for index, _ := range testId {
		if testId[index] == id {
			test = true
		}

	}
	if test == true {
		return nil
	} else {
		return ErrTestId
	}
}
