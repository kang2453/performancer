package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {

	client, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	go func(c net.Conn) {
		data := make([]byte, 4096)
		var s string

		for {
			n, err := c.Read(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			s = string(data[:n])
			fmt.Println(s)
			if s == "quit" {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}(client)

	go func(c net.Conn) {
		i := 0
		retVal := false
		for {

			var s string
			if i == 10 {
				s = "quit"
				retVal = true
			} else {
				s = "Hello" + strconv.Itoa(i)
			}
			_, err := c.Write([]byte(s))
			if err != nil {
				fmt.Println(err)
				return
			}

			if retVal == true {
				return
			}
			i++
			time.Sleep(1 * time.Second)
		}
	}(client)

	fmt.Scanln()
}
