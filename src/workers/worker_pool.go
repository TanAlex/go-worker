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

//NewWorker constructor
func NewWorker(id int) *Worker {
	return &Worker{ID: id, Chan: make(chan interface{})}
}

//JobFuncType type for the function used in Workders
//    You can use DoJob(jobFunc) to pass the job function to run
type JobFuncType func(w *Worker, payload interface{})

//ResultFuncType is for the function you pride to accept the result coming back from jobs
type ResultFuncType func(result interface{})

//Workers ...
type Workers struct {
	Workers         []Worker
	ChanResult      chan interface{}
	Jobs            []string
	CurrentJobIndex int
	JobFunc         JobFuncType
	ResultFunc      ResultFuncType
}

//DoJob add the func to Workers so all of the workers will run it
func (ws *Workers) DoJob(jobFunc JobFuncType) {
	ws.JobFunc = jobFunc
}

//ResultHandle add the func to Workers so all of the workers will run it
func (ws *Workers) ResultHandle(resFunc ResultFuncType) {
	ws.ResultFunc = resFunc
}

//AddWorker create a new Worker and add to Workers
func (ws *Workers) AddWorker() {
	length := len(ws.Workers)
	length++
	w := NewWorker(length)
	ws.Workers = append(ws.Workers, *w)
}

//NewWorkers constructor
func NewWorkers(length int) *Workers {
	ws := make([]Worker, length)
	for i := range ws {
		ws[i] = *NewWorker(i)
	}
	ch := make(chan interface{})
	return &Workers{Workers: ws, ChanResult: ch, Jobs: []string{}, CurrentJobIndex: 0}
}

//ToString is a test function to construct the string
func (ws *Workers) ToString() string {
	ids := []string{}
	for _, v := range ws.Workers {
		ids = append(ids, strconv.Itoa(v.ID))
	}
	return strings.Join(ids, ",")
}

//Start is the function to start running jobs in the pool using all workers
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
					if ws.JobFunc != nil {
						ws.JobFunc(&w, job)
					}

					chanResult <- fmt.Sprintf("job %d done", index)
				}
				wg.Done()
			}(v)
		}
	}
	go func() {
		for result := range chanResult {
			//fmt.Printf("%v\n", result)
			if ws.ResultFunc != nil {
				ws.ResultFunc(result)
			}
		}
	}()
	wg.Wait()
	close(chanResult)
}
