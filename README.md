# A Workers pool library in GO

This utility library is to create a generic way to run a bunch of jobs in parallel

Usage:
```go
//Create a Workers instance with 3 Workers in the pool, so it can run 3 jobs in parallel
ws := Workers.NewWorkers(3)
//The jobs array, they will be passed to the Workers 
ws.Jobs = []string{"test01", "test02", "test03", "test04"}
//Define the job function to run
ws.DoJob(func(w *Workers.Worker, job interface{}) {
    fmt.Printf("%s\n", job)
})
//Add the function to handle each result from workers
//the result doesn't have to be a string
//It can be a structure that has each worker's id etc
ws.ResultHandle(func(result interface{}) {
    fmt.Printf("%s\n", result)
})
//Call Start() to run the jobs
ws.Start()

```

run
```sh
export GOPATH=`pwd`
go build -o bin/worker.exe src/*.go
bin/worker.exe 
```