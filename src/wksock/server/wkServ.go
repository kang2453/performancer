package wkServ

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type CliInfo struct {
	sock     net.Conn
	ipaddr   string
	lastTime time.Time
	use      bool
}

type Serv struct {
	ip      string
	port    int
	cnt     int
	cliList *[]CliInfo
	ln      net.Listener
}

func (s *Serv) SetClient(sock net.Conn) bool {

	retVal := false

	if len(*s.cliList) == 0 {
		info := new(CliInfo)
		info.sock = sock
		info.ipaddr = sock.RemoteAddr().String()
		info.lastTime = time.Now()
		info.use = true
		*s.cliList = append(*s.cliList, *info)
	} else {
		for idx, info := range *s.cliList {
			if info.use == false {
				info := new(CliInfo)
				info.sock = sock
				info.ipaddr = sock.RemoteAddr().String()
				info.lastTime = time.Now()
				info.use = true
				*s.cliList = append(*s.cliList, *info)
				retVal = true
			}
		}
	}
	return retVal
}

func (s *Serv) StartServer(port int, sendMsg chan string) {

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer ln.Close()

	go s.ClientHandle(sendMsg)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		defer conn.Close()
		retVal := s.SetClient(conn)
		if retVal == false {
			fmt.Println("Client Connection Fail(client number more than 100)")
			continue
		}
		fmt.Printf("Client Conn:%s\n", conn.RemoteAddr().String())
		//go s.ClientHandle(sendMsg)
	}
}

func (s *Serv) ClientHandle(msgChan chan string) {

	fmt.Println("requestHandler start")
	interMsg := make(chan string, 10)
	wg := new(sync.WaitGroup)

	// c.Write([]byte("hello "))
	// info := s.cliList[idx]
	// sock := info.sock

	for {
		for _, info := range *s.cliList {
			// info
		}
	}
	// Recv gorutine
	go func(sock net.Conn) {
		wg.Add(1)
		data := make([]byte, 4096)
		for {
			n, err := sock.Read(data)
			if err != nil {
				interMsg <- "EOF"
				wg.Done()
				// return
			} else {
				fmt.Println(string(data[:n]))
				msg <- string(data[:n])
			}
			time.Sleep(1 * time.Second)
			info.lastTime = time.Now()
		}
	}(sock)

	// write gorutine
	go func(sock net.Conn) {
		wg.Add(1)
		for {
			select {
			case data, ok := <-msgChan:
				if ok == true {
					if data == "EOF" {
						wg.Done()
						// return
					} else {
						_, err := sock.Write([]byte(data))
						if err != nil {
							wg.Done()
							// return
						}
					}
				}
			case data, ok := <-interMsg:
				if ok == true {
					wg.Done()
				}
			}
		}
	}(sock)
	wg.Wait()
	fmt.Println("requestHandler end")
}
