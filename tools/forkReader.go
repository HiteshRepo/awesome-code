package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ForkReader(r io.Reader, w1, w2 io.Writer) error {
	toggle := false
	for {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if !toggle {
			_, err := w1.Write(buf)
			if err != nil {
				return err
			}
		} else {
			_, err := w2.Write(buf)
			if err != nil {
				return err
			}
		}

		toggle = !toggle
	}

	return nil
}
func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16 * 1024 * 1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16 * 1024 * 1024)

	str := readLine(reader)

	r := strings.NewReader(str)
	rw1, rw2 := bytes.NewBuffer([]byte{}), bytes.NewBuffer([]byte{})
	err = ForkReader(r, rw1, rw2)
	if err != nil {
		panic(err)
	}
	res1, err := ioutil.ReadAll(rw1)
	if err != nil {
		panic(err)
	}
	res2, err := ioutil.ReadAll(rw2)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(writer, "%s\n", string(res1))
	fmt.Fprintf(writer, "%s\n", string(res2))

	writer.Flush()
}
