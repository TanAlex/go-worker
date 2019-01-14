package main

import (
	"fmt"
	Workers "workers"
)

func main() {
	ws := Workers.NewWorkers(3)
	ws.Jobs = []string{"test01", "test02", "test03", "test04"}
	ws.DoJob(func(w *Workers.Worker, job interface{}) {
		fmt.Printf("%s\n", job)
	})
	fmt.Println(ws.ToString())
	ws.Start()
}
