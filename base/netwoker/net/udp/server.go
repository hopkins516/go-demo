package udp

import (
	"fmt"
	"net"
)

func BuildUdpServer() {

	var listen, err = net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("create udp listen error, ", err)
		return
	}
	fmt.Println("addr", listen.Addr())
	fmt.Println("network", listen.Addr().Network())
	fmt.Println("string", listen.Addr().String())

	defer listen.Close()
	//var udpListen, _ = net.ListenUDP("udp", &net.UDPAddr{IP: ""})
}