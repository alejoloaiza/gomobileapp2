package main

import (
	"bufio"
	"fmt"

	"io"
	"net"
	"strings"
	"time"
)

func StartIRCprocess(in chan string) {

MainCycle:
	for {
		Connection, err := net.Dial("tcp", Localconfig.IRCServerPort)

		if err != nil {
			fmt.Println(err)
			time.Sleep(2000 * time.Millisecond)
			continue MainCycle
		}

		fmt.Fprintln(Connection, "NICK "+Localconfig.IRCNick)
		fmt.Fprintln(Connection, "USER "+Localconfig.IRCUser)
		fmt.Fprintln(Connection, "JOIN "+Localconfig.IRCChannels)
		MyReader := bufio.NewReader(Connection)
	ReaderCycle:
		for {

			message, err := MyReader.ReadString('\n')
			// atomixxx: To handle if connection is closed, and jump to next execution.
			if err != nil {
				fmt.Println(time.Now().Format(time.Stamp) + ">>>" + err.Error())
				if io.EOF == err {
					Connection.Close()
					fmt.Println("server closed connection")
				}
				time.Sleep(2000 * time.Millisecond)
				break ReaderCycle
			}

			fmt.Print(time.Now().Format(time.Stamp) + ">>" + message)

			// atomixxx: Split the message into words to better compare between different commands
			text := strings.Split(message, " ")
			//fmt.Println("Number of objects in text: "+ strconv.Itoa(len(text)))
			var respond bool = false
			var response string
			// atomixxx: Logic to detect messages, BOT logic should go inside this
			if len(text) >= 4 && text[1] == "PRIVMSG" {
				respond = true
				var repeat bool = true
				var respondTo string
				//atomixxx logic to differ if message is channel or private from user
				if text[2][0:1] == "#" {
					// logic to respond the same thing to a channel / repeater BOT
					respondTo = text[2]
				} else {
					userto := strings.Split(text[0], "!")
					respondTo = userto[0][1:]
					// logic to respond the same thing to a user / repeater BOT
				}
				// If its a command BOT will execute the command given
				if text[3] == ":!cmd" {
					repeat = false
					commandresponse := ProcessCommand(text[4:])
					in <- "ok"
					response = "PRIVMSG " + respondTo + " :" + commandresponse

				}
				// If is not a command BOT will repeat the same thing
				if repeat == true {
					response = "PRIVMSG " + respondTo + " " + strings.Join(text[3:], " ")

				}
			}
			// atomixxx: Ping/Pong handler to avoid timeout disconnect from the irc server
			if len(text) == 2 && text[0] == "PING" {
				response = "PONG " + text[1]
				respond = true
			}
			// This checks if the received text requires response or not, and respond according to the above logic

			if respond == true {
				fmt.Fprintln(Connection, response)
				fmt.Println(time.Now().Format(time.Stamp) + "<<" + response)
			}

		}
		// atomixxx: If connection is closed, will try to reconnect after 2 seconds
		time.Sleep(2000 * time.Millisecond)
	}

}

func ProcessCommand(command []string) string {
	var bodyString string

	bodyString = "Command received... processing"

	return bodyString
}
