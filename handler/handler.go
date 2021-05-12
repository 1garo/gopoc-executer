package handler

import (
	"fmt"
	// "log"
	// "log"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/streadway/amqp"
	"rogchap.com/v8go"
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
	// non lib mock
	fileToProcess := []string{
		//"/home/hungaro/dev/ts/ee-card/dist/src/index.js/home/hungaro/dev/go/gopoc-executer/nodecsv-lfiles/index.js",

    //"/home/hungaro/dev/ts/ee-card/dist/src/index.js",
    "dist/src/index.js",
		//"/home/hungaro/dev/go/gopoc-executer/data/t1.js",
		// "/home/hungaro/dev/go/gopoc-executer/data/t2.js",
		// "/home/hungaro/dev/go/gopoc-executer/data/t3.js",
		// "/home/hungaro/dev/go/gopoc-executer/data/t4.js",
	}
	// v8 mock
	// fileToProcess := []string{
	// 	"/home/hungaro/dev/go/gopoc-executer/data-v8/t1.js",
	// 	"/home/hungaro/dev/go/gopoc-executer/data-v8/t2.js",
	// 	"/home/hungaro/dev/go/gopoc-executer/data-v8/t3.js",
	// 	"/home/hungaro/dev/go/gopoc-executer/data-v8/t4.js",
	// }

	ini := time.Now()
	r := make(chan Result)
	go readListFile(fileToProcess, r)
	for d := range r {
		fmt.Print("data -> ", d.resp)
		ctx, _ := v8go.NewContext() // creates a new V8 context with a new Isolate aka VM
		val, err := ctx.RunScript("", "data/t2.js")
		if err != nil {
			switch err := err.(type) {
			case *v8go.JSError:
				// fmt.Println(err.Message)    // the message of the exception thrown
				// fmt.Println(err.Location)   // the filename, line number and the column where the error occured
				// fmt.Println(err.StackTrace) // the full stack trace of the error, if available

				fmt.Printf("javascript error: %v\n", err)        // will format the standard error message
				fmt.Printf("javascript stack trace: %+v\n", err) // will format the full error stack trace
			default:
				fmt.Println("v8 error -> ", err)
			}
		}
		fmt.Println("v8 process -> ", val)
	}

	fmt.Println("(Took ", time.Since(ini).Seconds(), "secs)")

}

func readListFile(fileToProcess []string, rchan chan Result) {
	defer close(rchan)
	var results = []chan Result{}

	for i, url := range fileToProcess {
		results = append(results, make(chan Result))
		// v8 lib parallel execution
		// go v8Parallel(url, results[i])
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
	cmd := exec.Command("node", file)
  // TODO: rewrite this function to be able to read stream from api output
	data, err := cmd.CombinedOutput()
	fmt.Println("Result: " + string(data))
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(data))
		return
	}
	fmt.Println("Result: " + string(data))
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
