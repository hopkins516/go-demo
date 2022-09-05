// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func HttpResp() {
	httpRpc, err := rpc.DialHTTP("tcp", "localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	var ret string
	err = httpRpc.Call("HelloRpcService.Hello", "HTTP rpc", &ret)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ret)
}

func TcpResp() {
	tcpRpc, err := rpc.Dial("tcp", "localhost:9002")
	if err != nil {
		log.Fatal(err)
	}
	var ret string
	var _ = tcpRpc.Call("HelloRpcService.Hello", "TCP rpc", &ret)
	log.Println(ret)
}

func JsonRpcResp() {
	jsonRpc, _ := jsonrpc.Dial("tcp", "localhost:9003")
	var ret string
	var _ = jsonRpc.Call("HelloRpcService.Hello", "Json TCP rpc", &ret)
	log.Println(ret)
}

func main() {
	HttpResp()
	TcpResp()
	JsonRpcResp()
}
