package servicecli

import (
	"fmt"
	"time"
)

type storage interface {
	AcceptOrder(id int, rid int, t time.Time, packType string, weight int) error
	Delete(id int) error
	IssueOrder(idSlice []int) error
	ListOrder(rID int, N int, inPVZ bool) error
	ReturnOrder(rID int, ID int) error
	ReturnList(pageSize int, pageNumber int) error
}
type Service struct {
	s storage
}

func NewService(s storage) Service {
	return Service{s: s}
}

func (s Service) Help() {
	fmt.Println("Список доступных команд:")
	fmt.Println("accept - принять заказ от курьера")
	fmt.Println("delete - вернуть заказ курьеру")
	fmt.Println("issue - выдать заказ клиенту")
	fmt.Println("list - получить список заказов")
	fmt.Println("return - вернуть")
	fmt.Println("returnlist - получить список возвратов")
	fmt.Println("write - записать информацию о пвз")
	fmt.Println("read - прочитать информацию о пвз")
}

// принятие заказа от курьера

func (s Service) AcceptOrder(id int, rid int, t time.Time, packType string, weight int) error {
	return s.s.AcceptOrder(id, rid, t, packType, weight)
}

// вернуть заказ курьеру

func (s Service) Delete(id int) error {
	return s.s.Delete(id)
}

// выдать заказ клиенту

func (s Service) IssueOrder(idSlice []int) error {
	return s.s.IssueOrder(idSlice)

}

// получить список заказов

func (s Service) ListOrders(rID int, N int, inPVZ bool) error {
	return s.s.ListOrder(rID, N, inPVZ)

}

// возврат заказа

func (s Service) ReturnOrder(rID int, ID int) error {
	return s.s.ReturnOrder(rID, ID)
}

// получить список возвратов

func (s Service) ReturnList(pageSize int, pageNumber int) error {
	return s.s.ReturnList(pageSize, pageNumber)
}
