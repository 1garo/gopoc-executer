package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/1garo/gopoc-executer/handler"
	//"rogchap.com/v8go"
)

//Result exported
type Result struct {
	resp string
}

func (r Result) String() string {
	return fmt.Sprint(r.resp)
}

func main() {
	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	httpHandler := handler.NewHandler()
	server := &http.Server{
		Handler: httpHandler,
	}

	go func() {
		process()
		server.Serve(listener)
	}()
	defer Stop(server)

	log.Printf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func process() {
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

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
