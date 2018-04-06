package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"time"
)

var online = 0

type User struct {
	address  net.Addr
	connection net.Conn
	username string
	active   bool
}

var users []User

func handleConnection(conn net.Conn) {
	user := User{
		address:  conn.RemoteAddr(),
		connection: conn,
		username: "",
		active:   false,
	}
	sc := bufio.NewScanner(conn)

	conn.Write([]byte("Hello! Please, enter your nickname: "))
	for sc.Scan() {
		if sc.Text() != "" {
			user.username = sc.Text()
			break
		}
		conn.Write([]byte("Hello! Please, enter your nickname: "))
	}

	conn.Write([]byte(timing() + " Welcome, " + user.username + "!\n\r"))
	for _, us := range users {
		if us.active {
			us.connection.Write([]byte(timing() + " " + user.username + " connected.\n\r"))
		}
	}
	user.active = true
	online++
	users = append(users, user)
	defer func() {
		conn.Close()
		online--
		fmt.Println(timing(), "user \""+user.username+"\" addr", user.address, "disconnected from the server.")
	}()

ReadLoop:
	for sc.Scan() {
		message := sc.Text()
		switch message {
		case ".help":
			conn.Write([]byte("" +
				"--------------------\n\r" +
				".help to see available commands\n\r" +
				".exit to leave this chat\n\r" +
				".online to see current online\n\r" +
				"--------------------\n\r"))
		case ".exit":
			conn.Write([]byte("You have been disconnected. Goodbye!\n\r"))
			user.active = false
			for _, us := range users {
				if us.active {
					us.connection.Write([]byte(timing() + " " + user.username + " disconnected\n\r"))
				}
			}
			break ReadLoop
		case ".online":
			conn.Write([]byte("Current online is " + strconv.Itoa(online) + "\n\r"))
		case "":
		default:
			for _, us := range users {
				if us.active {
					us.connection.Write([]byte(timing() + " " + user.username + ": " + message + "\n\r"))
				}
			}
		}
	}
}

func timing() string {
	return time.Now().Format("[15:04:05]")
}

func startServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println(timing(), "Server started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println(timing(), conn.RemoteAddr(), "connected to server.")

		go handleConnection(conn)
	}
}

func showOnline() {
	for {
		time.Sleep(10*time.Second)
		fmt.Println(timing(), "ONLINE:", online)
	}
}

func main() {
	go startServer()
	showOnline()
}
