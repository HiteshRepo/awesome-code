package main

import (
	"fmt"
	"log"
	"os"
	"path"
)

type FileManager interface {
	Read() (string, error)
	Create() error
	Close() error
	Write() (int, error)
}

type File1Manager struct {
	filePath string
	file     *os.File
}

func (f1m *File1Manager) Read() (string, error) {
	return read(f1m.file)
}

func (f1m *File1Manager) Create() {
	file, err := create(f1m.filePath)
	if err != nil {
		log.Fatal("cannot create file: ", f1m.filePath)
	}
	f1m.file = file
}

func (f1m *File1Manager) Write(data string) (int, error) {
	return write(f1m.file, data)
}

func (f1m *File1Manager) Close() error {
	err := f1m.file.Close()
	if err != nil {
		return err
	}
	return nil
}

type File2Manager struct {
	filePath string
	file     *os.File
}

func (f2m *File2Manager) Read() (string, error) {
	return read(f2m.file)
}

func (f2m *File2Manager) Create() {
	file, err := create(f2m.filePath)
	if err != nil {
		log.Fatal("cannot create file: ", f2m.filePath)
	}
	f2m.file = file
}

func (f2m *File2Manager) Write(data string) (int, error) {
	return write(f2m.file, data)
}

func (f2m *File2Manager) Close() error {
	err := f2m.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func read(file *os.File) (string, error) {
	var data []byte
	_, err := file.Read(data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func write(file *os.File, data string) (int, error) {
	size, err := file.WriteString(data)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func create(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func main() {
	wd, _ := os.Getwd()
	file1Path := path.Join(wd, "file_manager/dist", "file1.txt")
	f1m := File1Manager{filePath: file1Path}
	f1m.Create()
	defer f1m.Close()

	file2Path := path.Join(wd, "file_manager/dist", "file2.txt")
	f2m := File2Manager{filePath: file2Path}
	f2m.Create()
	defer f2m.Close()

	data := []string{"Line1\n","Line2\n","Line3\n","Line4\n"}
	toggle := false

	for _,d := range data {
		if !toggle {
			f1m.Write(d)
		} else {
			f2m.Write(d)
		}

		toggle = !toggle
	}

	data1, _ := f1m.Read()
	data2, _ := f2m.Read()

	fmt.Println(data1)
	fmt.Println(data2)
}


