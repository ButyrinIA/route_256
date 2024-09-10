package storage

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var mutex sync.RWMutex

func WritePVZToFile(pvz PVZ, filename string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Запись данных о ПВЗ в файл
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Name: %s, Address: %s, Contact: %s\n", pvz.Name, pvz.Address, pvz.Contact)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("write ok ")
}

func ReadPVZFromFile(filename string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Чтение данных о ПВЗ из файла
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Читаем файл построчно и выводим его содержимое в консоль
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

}
