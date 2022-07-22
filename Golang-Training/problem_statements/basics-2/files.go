package basics_2

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var data = []string{
	"Line 1\n",
	"Line 2\n",
	"Line 3\n",
	"Line 4\n",
	"Line 5\n",
	"Line 6\n",
}

func ForkReader() {
	wd, _ := os.Getwd()
	file1Path := path.Join(wd, "problem_statements/basics-2/dist", "file1.txt")
	file2Path := path.Join(wd, "problem_statements/basics-2/dist", "file2.txt")

	CreateFile(file1Path)
	CreateFile(file2Path)

	toggle := true
	for _, line := range data {
		if toggle {
			WriteFile(file1Path, line)
		} else {
			WriteFile(file2Path, line)
		}

		toggle = !toggle
	}

	ReadFile(file1Path)
	ReadFile(file2Path)
}

func CreateFile(filePath string) {
	_, err := os.Create(filePath)
	if err != nil {
		log.Fatal("failed to create a file, err: ", err)
	}
}

func WriteFile(filePath string, data string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("failed to open a file, err: ", err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal("failed to write to a file, err: ", err)
	}
}

func ReadFile(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading data from file: %s", err)
	}
	fmt.Printf("\nContents of file %s: \n%s", filePath, data)
}
