package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

type HelloRpcService struct{}

func (p *HelloRpcService) Hello(request string, reply *string) error {
	*reply = "hello: " + request
	fmt.Println("grpc callback done.")
	return nil
}

func RpcClient() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal()
	}

	var reply string
	err = client.Call("HelloRpcService.Hello", "request by demo client", &reply)
	if err != nil {
		log.Fatal()
	}
	fmt.Println(reply)
}

func main() {
	err := rpc.RegisterName("HelloRpcService", new(HelloRpcService))
	if err != nil {
		log.Fatal("Register grpc function failed.")
	}
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		rpc.HandleHTTP()
		err := http.ListenAndServe(":9001", nil)
		wg.Wait()
		if err != nil {
			defer wg.Done()
			log.Fatal(err)
		}
	}()
	log.Println("serve on 9001 for httprouter")

	go func() {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:9002")
		tcpListener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			defer wg.Done()
			log.Fatal(err)
		}
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
	log.Println("serve on 9002 for rpc")

	go func() {
		tcpAddr, _ := net.ResolveTCPAddr("tcp", "localhost:9003")
		tcpListener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			wg.Done()
			log.Fatal(err)
		}

		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				continue
			}
			go jsonrpc.ServeConn(conn)
		}
	}()
	log.Println("serve on 9003 for json rpc")
	wg.Wait()
	//
	//listener, err := net.Listen("tcp", ":1234")
	//if err != nil {
	//	log.Fatal("ListenTcp occur errors: ", err)
	//}
	//fmt.Printf("%#v\n", []byte("hello, 世界"))
	//conn, err := listener.Accept()
	//if err != nil {
	//	log.Fatal("Accept occur errors: ", err)
	//}
	//rpc.ServeConn(conn)
}
