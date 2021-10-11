package main

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
	"time"
)

const Protocol = "tcp"

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}
	port := ":" + arguments[1]
	server, err := net.Listen(Protocol, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Close()

	con, err := server.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()

	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		data = strings.TrimRight(strings.TrimLeft(data, " "), " \n")

		switch data[0:4] {
		case "ECHO", "echo":
			data = data[5:]+"\n"
			logrus.Info(fmt.Sprintf("->: "+data))
			con.Write([]byte(data))
		case "TIME", "time":
			if len(data) > 4 {
				con.Write([]byte("unknown command, type 'help' to view the available commands\n"))
			} else {
				t := time.Now()
				sTime := t.Format(time.RFC3339) + "\n"
				con.Write([]byte(sTime))
			}
		case "EXIT", "exit":
			logrus.Info("Stop TCP server")
			con.Write([]byte("exit\n"))
			return
		default:
			con.Write([]byte("unknown command, type 'help' to view the available commands\n"))
		}
	}
}

