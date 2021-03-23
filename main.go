package main

import (
	"fmt"
  "os/exec"
	"io/ioutil"
	"time"
	"rogchap.com/v8go"
)

//Result exported
type Result struct {
	resp string
}

func (r Result) String() string {
	return fmt.Sprint(r.resp)
}

func main() {
	fileToProcess := []string{
		"/home/hungaro/dev/go/gopoc-executer/data/t1.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t2.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t3.js",
		"/home/hungaro/dev/go/gopoc-executer/data/t4.js",
	}

	ini := time.Now()
	r := make(chan Result)
	go readListFile(fileToProcess, r)
	fmt.Println("With goroutines:")
  for  d := range r {
    fmt.Println("entrou\n")
    fmt.Println(d.resp)
    // exec with v8 lib
    // some particular things:
    /*
      i.e cannot use console.log 
      i.e need to write the var in the end of the file 
          to catch the return
    */
    ctx, _ := v8go.NewContext()
    val, _ := ctx.RunScript(d.resp, "value.js")
    // exec with golang std lib
    data, err := exec.Command("node", "/home/hungaro/dev/go/gopoc-executer/index.js").Output()
    if err != nil {
      panic(err)
    }
    fmt.Printf("std lib: %s\n", data)
		fmt.Printf("v8 lib: %s\n", val)
	}

	fmt.Println("(Took ", time.Since(ini).Seconds(), "secs)")
}

func readListFile(fileToProcess []string, rchan chan Result) {
	defer close(rchan)
	var results = []chan Result{}

	for i, url := range fileToProcess {
		results = append(results, make(chan Result))
		go scrapParallel(url, results[i])
	}


	for i := range results {
		for r1 := range results[i] {
			rchan <- r1
		}
	}
}

func scrapParallel(file string, rchan chan Result) {
	defer close(rchan)
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	var r Result
  r.resp = string(f)
	rchan <- r
}
