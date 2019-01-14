package Workers

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Worker struct for Worker that does jobs
type Worker struct {
	ID   int
	Chan chan interface{}
}

//NewWorker ...
func NewWorker(id int) *Worker {
	return &Worker{ID: id, Chan: make(chan interface{})}
}

type JobFuncType func(w *Worker, payload interface{})

//Workers ...
type Workers struct {
	Workers         []Worker
	ChanResult      chan interface{}
	Jobs            []string
	CurrentJobIndex int
	JobFunc         JobFuncType
}

func (ws *Workers) DoJob(jobFunc JobFuncType) {
	ws.JobFunc = jobFunc
}

func (ws *Workers) AddWorker() {
	length := len(ws.Workers)
	length++
	w := NewWorker(length)
	ws.Workers = append(ws.Workers, *w)
}

//NewWorkers ...
func NewWorkers(length int) *Workers {
	ws := make([]Worker, length)
	for i, _ := range ws {
		ws[i] = *NewWorker(i)
	}
	ch := make(chan interface{})
	return &Workers{Workers: ws, ChanResult: ch, Jobs: []string{}, CurrentJobIndex: 0}
}

func (ws *Workers) ToString() string {
	ids := []string{}
	for _, v := range ws.Workers {
		ids = append(ids, strconv.Itoa(v.ID))
	}
	return strings.Join(ids, ",")
}

func (ws *Workers) Start() {
	chanResult := ws.ChanResult
	workers := ws.Workers
	jobLength := len(ws.Jobs)
	index := ws.CurrentJobIndex
	var mutex = &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(len(workers))
	if index < jobLength {
		for i, v := range workers {
			//w := workers[i]
			_ = i
			go func(w Worker) {
				for {
					mutex.Lock()
					index := ws.CurrentJobIndex
					length := len(ws.Jobs)
					if index >= length {
						mutex.Unlock()
						break
					}
					job := ws.Jobs[index]
					ws.CurrentJobIndex++
					mutex.Unlock()
					//fmt.Printf("%v\n", job)
					ws.JobFunc(&w, job)

					chanResult <- fmt.Sprintf("job %d done", index)
				}
				wg.Done()
			}(v)
		}
	}
	go func() {
		for result := range chanResult {
			fmt.Printf("%v\n", result)
		}
	}()
	wg.Wait()
	close(chanResult)
}
