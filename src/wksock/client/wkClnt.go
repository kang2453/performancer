package wkClnt

import (
	"fmt"
	"net"
	"strconv"
)

type Clnt struct {
	localip   string
	localport int

	serverip   string
	serverport int

	sock net.Conn
}

func (c *Clnt) startClient(ip string, port int, sendMsg chan string) {

	c.serverip = ip
	c.serverport = port

	RemoteConn := ip + ":" + strconv.Itoa(port)

	sock, err := net.Dial("tcp", RemoteConn)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("END|%s", err.Error())
		sendMsg <- msg
	}
	fmt.Println(sock)

}
