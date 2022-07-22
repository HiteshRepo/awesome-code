package main

import (
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/config"
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/coordinator"
	"github.com/hiteshpattanayak-tw/awesome-code/multithreading/splitter"
)

func main() {
	cnf := config.ProvideInputConfigs()
	parts := splitter.SplitFileIntoParts(cnf)
	c := coordinator.ProvideCoordinator(cnf, parts)
	c.SpawnWorkers()
	c.Start()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
