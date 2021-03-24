package handler

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/streadway/amqp"
)

//Result exported
type Result struct {
	resp string
}

func (r Result) String() string {
	return fmt.Sprint(r.resp)
}

func NewHandler() http.Handler {
	router := chi.NewRouter()
	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/executors", executors)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}

func Process() {
	fileToProcess := []string{
		"/home/hungaro/dev/go/gopoc-executer/data/t1.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t2.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t3.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t4.js",
	}

	ini := time.Now()
	r := make(chan Result)
	go readListFile(fileToProcess, r)
	for d := range r {
		fmt.Print(d.resp)
		/*
			v8-lib
			i.e cannot use console.log
			i.e need to return the var in the end of the file
				to catch the return
			ctx, _ := v8go.NewContext()
			val, _ := ctx.RunScript(d.resp, "value.js")
			fmt.Printf("v8 lib: %s\n", val)
		*/
	}

	fmt.Println("(Took ", time.Since(ini).Seconds(), "secs)")

}

func readListFile(fileToProcess []string, rchan chan Result) {
	defer close(rchan)
	var results = []chan Result{}

	for i, url := range fileToProcess {
		results = append(results, make(chan Result))
		//go v8Parallel(url, results[i])
		go execFileParallel(url, results[i])
	}

	for i := range results {
		for r1 := range results[i] {
			rchan <- r1
		}
	}
}

func execFileParallel(file string, rchan chan Result) {
	defer close(rchan)
	data, err := exec.Command("node", file).Output()
	if err != nil {
		panic(err)
	}
	var r Result
	r.resp = string(data)
	rchan <- r
}

// func v8Parallel(file string, rchan chan Result) {
// 	defer close(rchan)
// 	f, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return
// 	}
// 	var r Result
// 	r.resp = string(f)
// 	rchan <- r
// }

func StartConn(conn *amqp.Connection) (*amqp.Channel, *amqp.Queue) {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		FailOnError(err, "Failed to declare a queue")
	}
	return ch, &q
}
