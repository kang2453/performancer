package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	wkclnt "wksock/client"
	wkserv "wksock/server"
)

var argType = ""
var argIp = ""
var argPort int

// var recvMsg = make(chan string, 100)
var sendMsg = make(chan string, 100)

func readCmd(filename string) []string {

	cnt := 0

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s open Fail(%s)\n", filename, err.Error())
		return nil
	}
	defer file.Close()

	data := make([]byte, 2048)
	cnt, err = file.Read(data)
	if err != nil || cnt == 0 {
		fmt.Printf("%s read Fail(%s)\n", filename, err.Error())
		return nil
	}
	msg := string(data[:cnt])
	return strings.Split(msg, "\n")
}

func initArgument() {
	moduleType := flag.String("type", "", "Server(s) or Client(c)")
	ip := flag.String("ip", "", "Connection IP")
	port := flag.Int("port", 10099, "Connection Port")

	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	argType = strings.ToLower(*moduleType)
	argIp = *ip
	argPort = *port
}

func main() {

	num := 0

	initArgument()

	fmt.Println("moduleType : " + argType)
	fmt.Println("IP: " + argIp)
	fmt.Println("PORT: " + strconv.Itoa(argPort))

	// 모듈 기동
	if argType == "s" {
		serv := wkserv.Serv
		go serv.StartServer(argPort, sendMsg)
		for {
			cmdList := readCmd("conf/cmd.txt")
			if cmdList == nil {
				return
			}
			for idx, msg := range cmdList {
				fmt.Printf("%d %s\n", idx+1, msg)
			}
			fmt.Printf("Input The Cmd Number: (Quit:99)")
			cnt, err := fmt.Scanf("%d", &num)
			if err != nil {
				fmt.Println("wrong cmd number... ")
				continue
			}
			if num == 99 {
				break
			}
			// fmt.Printf("Cnt:%d, Number:%d\n", cnt, num)
			if len(cmdList) < num-1 {
				fmt.Println("Wrong cmd number...")
				continue
			}
			fmt.Printf("Number CMD: %s\n", cmdList[num-1])
			sendMsg <- cmdList[num-1]
		}
	} else if argType == "c" {
		clnt := wkclnt.Clnt
		go wksock.StartClient(argIp, argPort, sendMsg)
		for {
			select {
			case data, ok := <-sendMsg:
				if ok == true {
					retVal := ParsingMsg(data)
					if retVal == "END" {
						break
					}
				}
			}
			// 10초에 한번씩 Server에게 hb를 날린다
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("END")
}
