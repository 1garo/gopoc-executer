package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1garo/gopoc-executer/handler"
	"github.com/robertkrimen/otto"
	"github.com/streadway/amqp"
	"rogchap.com/v8go"
)

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
		// start using otto lib
		vm := otto.New()
		value, _ := vm.Run(`
			abc = 2 + 2;
			//console.log("The value of abc is " + abc); // 4`)
		log.Printf("otto lib -> %s\n", value)
		// finish otto lib

		// start using v8 lib
		ctx, _ := v8go.NewContext()                             // creates a new V8 context with a new Isolate aka VM
		ctx.RunScript("const add = (a, b) => a + b", "math.js") // executes a script on the global context
		ctx.RunScript("const result = add(3, 4)", "main.js")    // any functions previously added to the context can be called
		val, _ := ctx.RunScript("result", "value.js")           // return a value in JavaScript back to Go
		log.Printf("v8 lib -> %s\n", val)
		// finish v8 lib

		// non-lib method and rabbit conn
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		handler.FailOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, q := handler.StartConn(conn)

		handler.Sender(q, ch)
		handler.Consumer(q, ch)
		defer ch.Close()

		server.Serve(listener)
	}()
	defer Stop(server)

	log.Printf("Started server on %s", addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
