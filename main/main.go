package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/felixwqp/geerpc"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}


func startServer(addr chan string) {
	// pick a free port
	// You may use port 0 to indicate you're not specifying an exact port but you want a free, available port selected by the system:
	// start a listeners,
	var foo Foo
	if err := geerpc.Register(&foo); err != nil{
		log.Fatal("Register Error: ", err)
	}

	l, err := net.Listen("tcp", ":0")
	// bind the port to a certain process,
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil{
				log.Fatal("Call Foo.Sum error, ", err)
			}
			log.Println("%d + %d = %d ", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
