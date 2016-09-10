package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func requestHandler(c net.Conn) {

	fmt.Println("requestHandler start")

	msg := make(chan string, 10)
	wg := new(sync.WaitGroup)

	c.Write([]byte("hello "))

	go func(c net.Conn) {
		wg.Add(1)
		data := make([]byte, 4096)

		for {
			n, err := c.Read(data)
			if err != nil {
				msg <- "EOF"
				wg.Done()
				// return
			} else {
				fmt.Println(string(data[:n]))
				msg <- string(data[:n])
			}
			time.Sleep(1 * time.Second)
		}
	}(c)

	go func(c net.Conn) {
		wg.Add(1)
		for {
			select {
			case data, ok := <-msg:
				if ok == true {
					if data == "EOF" {
						wg.Done()
						// return
					} else {
						_, err := c.Write([]byte(data))
						if err != nil {
							wg.Done()
							// return
						}
					}
				}
			}
		}
	}(c)
	wg.Wait()
	fmt.Println("requestHandler end")
}

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer conn.Close()
		go requestHandler(conn)
	}
}
