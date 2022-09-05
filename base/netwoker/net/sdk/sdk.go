// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func genIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func generateIps() net.IP {
	ip := make(net.IP, net.IPv6len)
	copy(ip, net.IPv4zero)
	i := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
	binary.BigEndian.PutUint32(ip.To4(), uint32(i))
	return ip.To16()
}

type IPv6Int [2]uint64

func RandomIpv6Int() (result [2]uint64) {
	result[0] = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	result[1] = rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
	return result
}

func main() {
	fmt.Println(genIpaddr())
	ip := generateIps()
	fmt.Println(ip.String())
}
