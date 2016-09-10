package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println(remoteAddr)

	return localAddr[0:idx]
}

func main() {
	//	GetLocalIP()
	//fmt.Println(GetOutboundIP())
	host, _ := os.Hostname()
	fmt.Println("host:", host)
	addrs, _ := net.LookupIP(host)
	fmt.Println("Addr:", addrs)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
		}
	}
}
