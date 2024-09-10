package executor

import (
	"bufio"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"homework/Homework-1/internal/infrastructure/kafka"
	"homework/Homework-1/internal/pkg/db"
	"homework/Homework-1/internal/pkg/repository/postgresql"
	"homework/Homework-1/internal/pkg/service"
	servicecli "homework/Homework-1/internal/service"
	"homework/Homework-1/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"
)

func RunServiceMode() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(ctx)

	pvzRepo := postgresql.NewPVZ(database)
	Producer, err := kafka.NewProducer()

	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer Producer.Close()
	implementation := service.NewServer(pvzRepo, Producer)

	log.Println("service started on port" + service.Port)
	router := service.CreateRouter(implementation)

	srv := &http.Server{
		Addr:    service.Port,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	Consumer, err := kafka.NewConsumer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer Consumer.Close()
	partitionConsumer, err := Consumer.SingleConsumer.ConsumePartition("methods", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Чтение и вывод событий в консоль
	for {
		select {
		case message := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(message.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error from partition consumer: %v\n", err)
		}
	}
}

func RunInteractiveMode() {
	stor, err := storage.NewStorage()
	if err != nil {
		fmt.Println("не удалось подключиться к хранилищу")
		return
	}
	for {
		var wg sync.WaitGroup
		gorutChan := make(chan string)

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)

		fmt.Print("Введите команду: ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		command := strings.TrimSpace(input)

		serv := servicecli.NewService(&stor)

		go func() {
			for {
				var numGoroutines = <-gorutChan
				var numCPU = runtime.NumCPU()

				fmt.Printf("The work of Goroutines: %s\n", numGoroutines)
				fmt.Printf("Number of CPUs: %d\n", numCPU)
			}
		}()

		select {
		case <-signalChan:
			err := stor.Storage.Close()
			if err != nil {
				fmt.Println("ошибка закрытия")
				return
			}
			fmt.Println("Received signal, exiting...")
			os.Exit(0)
		default:
			switch command {
			case "help":
				serv.Help()

			case "accept":
				var id, rid, weight int
				var packType string
				fmt.Println("введите ID заказа:")
				_, err := fmt.Scanln(&id)
				if err != nil {
					fmt.Print(err)
					return
				}
				fmt.Println("введите ID получателя:")
				_, err = fmt.Scanln(&rid)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Введите тип упаковки (Packet, Box или Film):")
				_, err = fmt.Scanln(&packType)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Введите вес посылки:")
				_, err = fmt.Scanln(&weight)
				if err != nil {
					fmt.Println(err)
					return
				}
				t := time.Now().AddDate(0, 0, 7)
				err = serv.AcceptOrder(id, rid, t, packType, weight)
				if err != nil {
					fmt.Println(err)

				} else {
					fmt.Println("заказ курьера принят")
				}

			case "delete":
				var id int
				fmt.Println("введите ID заказа:")
				_, err = fmt.Scanln(&id)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = serv.Delete(id)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("заказ вернули курьеру")
				}

			case "issue":
				var n int
				fmt.Println("введите количество заказов для выдачи:")
				_, err = fmt.Scanln(&n)
				if err != nil {
					fmt.Println(err)
					return
				}
				id := make([]int, n)
				for i := 0; i < n; i++ {
					_, err := fmt.Scanln(&id[i])
					if err != nil {
						return
					}
				}
				err := serv.IssueOrder(id)
				if err != nil {
					fmt.Println(err)

				} else {
					fmt.Println("заказы выданы")
				}

			case "list":
				var rId, N int
				fmt.Println("введите ID получателя:")
				_, err = fmt.Scanln(&rId)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Print("введите N (число заказов в списке), введите N=0, чтобы получить список заказов в пункте пвз:")
				_, err = fmt.Scan(&N)
				if err != nil {
					fmt.Println(err)
					return
				}
				if N == 0 {
					err1 := serv.ListOrders(rId, N, true)
					if err1 != nil {
						fmt.Println(err1)
					}
				}
				if N > 0 {
					err2 := serv.ListOrders(rId, N, false)
					if err2 != nil {
						fmt.Println(err2)
					}
				}
			case "return":
				var rId, id int
				fmt.Println("введите ID получателя:")
				_, err = fmt.Scanln(&rId)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("введите ID заказа:")
				_, err = fmt.Scanln(&id)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = serv.ReturnOrder(rId, id)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("заказ вернули в пвз")
				}

			case "returnlist":
				var pageSize, pageNumber int
				fmt.Println("введите количество элементовна странице :")
				_, err = fmt.Scanln(&pageSize)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("введите номер страницы:")
				_, err = fmt.Scan(&pageNumber)
				if err != nil {
					fmt.Println(err)
					return
				}
				err := serv.ReturnList(pageSize, pageNumber)
				if err != nil {
					fmt.Println(err)
				}

			case "write":
				gorutChan <- "write"
				wg.Add(1)
				go func() {
					var pvz = storage.PVZ{}
					fmt.Println("name: ")
					_, err = fmt.Scanln(&pvz.Name)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("address: ")
					_, err = fmt.Scanln(&pvz.Address)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("Contact: ")
					_, err = fmt.Scanln(&pvz.Contact)
					if err != nil {
						fmt.Println(err)
						return
					}
					storage.WritePVZToFile(pvz, "file/pvz.txt")
					wg.Done()

				}()
				wg.Wait()

			case "read":
				gorutChan <- "read"
				wg.Add(1)
				go func() {
					storage.ReadPVZFromFile("file/pvz.txt")
					wg.Done()
				}()
				wg.Wait()
			default:
				fmt.Println("Команда не распознана. Введите 'help' для получения списка команд.")
			}
		}
	}
}
