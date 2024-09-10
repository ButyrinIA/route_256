package storage

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestWritePVZToFile(t *testing.T) {
	pvz1 := PVZ{Name: "PVZ1", Address: "street1", Contact: "123"}
	pvz2 := PVZ{Name: "", Address: "", Contact: ""}
	tempFile := "temp_test.txt"

	t.Run("writing PVZ", func(t *testing.T) {
		file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		WritePVZToFile(pvz1, tempFile)

		fileContent, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		actual := string(fileContent)
		if actual != "Name: PVZ1, Address: street1, Contact: 123\n" {
			t.Errorf("Expected %q, got %q", "Name: PVZ1, Address: street1, Contact: 123\n", actual)
		}
	})
	t.Run("writing PVZ with empty fields", func(t *testing.T) {
		file, err := os.OpenFile(tempFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		WritePVZToFile(pvz2, tempFile)

		fileContent, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		actual := string(fileContent)
		if actual != "Name: , Address: , Contact: \n" {
			t.Errorf("Expected %q, got %q", "Name: , Address: , Contact: \n", actual)
		}
	})
}

func TestReadPVZFromFile(t *testing.T) {
	tempFilename := "test_read_pvz.txt"
	defer os.Remove(tempFilename)
	testData := "Name: PVZ1, Address: street1, Contact: 123\nName: PVZ2, Address: Street2, Contact: 456\n"
	file, err := os.Create(tempFilename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(testData)
	if err != nil {
		t.Fatalf("Failed to write to test file: %v", err)
	}

	r, w, _ := os.Pipe()
	os.Stdout = w
	ReadPVZFromFile(tempFilename)
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	expectedOutput := "Name: PVZ1, Address: street1, Contact: 123\nName: PVZ2, Address: Street2, Contact: 456\n"
	if buf.String() != expectedOutput {
		t.Errorf("Incorrect output. Expected: %s, Got: %s", expectedOutput, buf.String())
	}
}
