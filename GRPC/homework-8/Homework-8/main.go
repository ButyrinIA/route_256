package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	pb "workshop-8/internal/pkg/pb"
)

type server struct {
	pb.UnimplementedDeliveryServiceServer
	orders        map[int64]pb.Order
	ordersFile    string
	issuedCounter prometheus.Counter
}

func saveOrdersToFile(orders map[int64]pb.Order, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, order := range orders {
		line := fmt.Sprintf("%s,%s,%s\n", order.OrderId, order.RecipientId, order.ExpirationDate.AsTime().Format(time.RFC3339))
		_, err = file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadOrdersFromFile(filename string) (map[int64]pb.Order, error) {
	orders := make(map[int64]pb.Order)
	file, err := os.Open(filename)
	if err != nil {
		return orders, nil
	}
	defer file.Close()
	var orderId, recipientId int64
	var expirationDate string
	for {
		_, err := fmt.Fscanf(file, "%s,%s,%s\n", &orderId, &recipientId, &expirationDate)
		if err != nil {
			break
		}

		expirationTime, err := time.Parse(time.RFC3339, expirationDate)
		if err != nil {
			return nil, err
		}

		orders[orderId] = pb.Order{
			OrderId:        orderId,
			RecipientId:    recipientId,
			ExpirationDate: timestamppb.New(expirationTime),
			IssueDate:      timestamppb.New(expirationTime),
		}
	}

	return orders, nil
}

func (s *server) ReceiveOrder(ctx context.Context, req *pb.ReceiveOrderRequest) (*pb.ReceiveOrderResponse, error) {
	if _, exists := s.orders[req.OrderId]; exists {
		return nil, fmt.Errorf("заказ %s уже существует", req.OrderId)
	}

	/*if req.ExpirationDate.AsTime().Before(time.Now()) {
		return nil, fmt.Errorf("срок годности остался в прошлом")
	}*/

	s.orders[req.OrderId] = pb.Order{
		OrderId:        req.OrderId,
		RecipientId:    req.RecipientId,
		ExpirationDate: timestamppb.New(time.Now().AddDate(0, 0, 7)),
	}

	err := saveOrdersToFile(s.orders, s.ordersFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось сохранить заказы в файл: %w", err)
	}

	return &pb.ReceiveOrderResponse{Success: true}, nil
}

func (s *server) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*pb.ReturnOrderResponse, error) {
	order, exists := s.orders[req.OrderId]
	if !exists {
		return nil, fmt.Errorf("заказа %s не существует", req.OrderId)
	}

	if order.ExpirationDate.AsTime().After(time.Now()) {
		return nil, fmt.Errorf("невозможно вернуть заказ %s до истечения срока действия", req.OrderId)
	}

	delete(s.orders, req.OrderId)

	err := saveOrdersToFile(s.orders, s.ordersFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось сохранить заказы в файл: %w", err)
	}

	return &pb.ReturnOrderResponse{Success: true}, nil
}

func (s *server) IssueOrder(ctx context.Context, req *pb.IssueOrderRequest) (*pb.IssueOrderResponse, error) {
	for _, orderId := range req.OrderIds {
		order, exists := s.orders[orderId]
		if !exists {
			return nil, fmt.Errorf("заказа %s не существует", orderId)
		}

		if order.ExpirationDate.AsTime().After(time.Now()) {
			s.orders[orderId] = pb.Order{
				OrderId:        order.OrderId,
				RecipientId:    order.RecipientId,
				ExpirationDate: order.ExpirationDate,
				IssueDate:      timestamppb.New(time.Now()),
			}
			s.issuedCounter.Inc()
		}

	}

	return &pb.IssueOrderResponse{Success: true}, nil
}

func (s *server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	response := &pb.ListOrdersResponse{}

	for _, order := range s.orders {
		if req.RecipientId == order.RecipientId {
			response.Orders = append(response.Orders, &order)
		}
	}

	return response, nil
}

func (s *server) AcceptReturn(ctx context.Context, req *pb.AcceptReturnRequest) (*pb.AcceptReturnResponse, error) {
	order, exists := s.orders[req.OrderId]
	if !exists {
		return nil, fmt.Errorf("заказа %s не существует", req.OrderId)
	}
	if time.Since(order.IssueDate.AsTime()).Hours() > 48 {
		return nil, fmt.Errorf("заказ %s нельзя вернуть", req.OrderId)
	}

	return &pb.AcceptReturnResponse{Success: true}, nil
}

func (s *server) ListReturns(ctx context.Context, req *pb.ListReturnsRequest) (*pb.ListReturnsResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	response := &pb.ListReturnsResponse{}

	for _, order := range s.orders {
		if order.IssueDate != nil {
			issueTime := order.IssueDate.AsTime()

			if time.Since(issueTime).Hours() <= 48 {
				response.Returns = append(response.Returns, &order)
			}
		}
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start > len(response.Returns) {
		start = len(response.Returns)
	}

	if end > len(response.Returns) {
		end = len(response.Returns)
	}

	response.Returns = response.Returns[start:end]

	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9093")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	ordersFile := "orders.txt"
	orders, err := loadOrdersFromFile(ordersFile)
	if err != nil {
		log.Fatalf("не удалось загрузить заказы: %v", err)
	}

	issuedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "issued_orders_total",
		Help: "Total number of issued orders",
	})
	prometheus.MustRegister(issuedCounter)

	pb.RegisterDeliveryServiceServer(s, &server{
		orders:        orders,
		ordersFile:    ordersFile,
		issuedCounter: issuedCounter,
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9091", nil)) // Prometheus endpoint for metrics
	}()

	log.Printf("Starting gRPC server at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
