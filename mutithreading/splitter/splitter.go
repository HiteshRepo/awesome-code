package splitter

import (
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/config"
	"os"
)

func SplitFileIntoParts(config *config.Config) [][]int64 {
	file, err := os.Open(config.FileName)
	check(err)
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileChunk := fileSize / config.NoOfWorkers

	parts := make([][]int64, 0)

	for i := int64(0); i < config.NoOfWorkers; i++ {
		start := i*fileChunk
		end := (i+1)*fileChunk

		if start != 0 {
			start += 1
		}

		if i == config.NoOfWorkers - 1 {
			end = fileSize
		}

		parts = append(parts, []int64{start, end})
	}

	return parts
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}