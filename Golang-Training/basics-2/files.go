package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func PlayWithFiles() {
	CreateFile()
	ReadFile()
}

func CreateFile() {
	wd, _ := os.Getwd()
	filePath := path.Join(wd, "basics-2/dist", "test.txt")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("failed to create a file, err: ", err)
	}

	defer file.Close()

	size, err := file.WriteString("This is a new file")
	if err != nil {
		log.Fatal("failed to write to a file")
	}

	fmt.Println("Name of file is: ", file.Name())
	fmt.Printf("Size of data written to file is: %d bytes\n", size)
}

func ReadFile() {
	wd, _ := os.Getwd()
	filePath := path.Join(wd, "basics-2/dist", "test.txt")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading data from file: %s", err)
	}

	fmt.Printf("\nSize: %d bytes", len(data))
	fmt.Printf("\nData: %s", data)
}
